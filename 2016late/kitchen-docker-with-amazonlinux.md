# ローカル環境でAmzon Linuxをテストする

## はじめに

可搬性(ポータビリティ)に立ちふさがる密林、Amazon Linux。担当プロジェクトでもご多分に漏れず、Amazon LinuxをOSとしたシステムを組んでいます。
chefを用いて構成管理(商用環境は、[OpsWorks](https://aws.amazon.com/jp/opsworks/))、[Test Kitchen](http://kitchen.ci/) + [Serverspec](http://serverspec.org/)でインフラテストしているがOSはCentOS6.7(packerで生成したVagrant box)...OSが違うのに、何が構成管理か、何が可搬性か。  
[AWS公式のDockerイメージ](https://hub.docker.com/_/amazonlinux/)もリリースされたことですし、Dockerを武器にローカストで可搬性を手に入れることが出来ると思い立ったが吉日。
1営業日くらいで可搬性のあるインフラテスト環境を手に入れることが出来ました :dragon_face:

## 変更のサマリ

| -        | OpsWorks用cookbook | ローカル環境用cookbook | Test KitchenでプロビジョニングするホストのOS | Test Kitchenのドライバ | 初回のテスト完了までの所要時間 |
|----------|--------------------|------------------------|----------------------------------------------|------------------------|--------------------------------|
| ビフォー | example            | example_local          | CentOS6.7                                    | Vagrant                | 20分程度                       |
| アフター | example            | -                      | Amazon Linux                                 | Docker                 | 5分程度                        |

主たる変更は、[kitchen-docker](https://github.com/test-kitchen/kitchen-docker)を導入し、Test KitchenのドライバとしてDockerを採用したことです。
変更により、以下のメリットを享受出来るようになりました。
- __cookbookの統一__
- __プロビジョニング対象ホストのOSをOpsWorksと合わせることが出来る__
- __初回テスト時の所要時間の短縮__

以下で、詳細を記載していきます。

## ビフォー

![before](./before.png)

1. exampleアプリケーションを稼働させるためのcookbook(example)を開発する
	- OpsWorksを利用しているので、リポジトリ直下にcookbookを配置する必要がある
1. example cookbookをTest Kitchenで稼働させるために、example_localを開発する
	- example::defaultをincludeし、attributeを上書きするcookbook
1. Serverspecでテストコードを開発する
	- Test Kitchenにも対応するので、example_localに配置
1. Test Kitchenでexample cookbookの動作確認とテストを実施
	- プロビジョニング:```bundle exec kitchen converge default-example```
	- テスト:```bundle exec kitchen verify default-example```
1. テストが通ったら、packerに対応させる
	- ```packer/example.json```を開発し、cookbookと一緒にAtlasにアップロード
	- AtlasのpackerでVagrant boxを生成する(リモートビルド)
1. 開発担当者がAtlasからVagrant boxをダウンロードし、exampleアプリケーションを開発

__...だるい__

設計し、実装を終えた当初は、堅牢で素晴らしいワークフローだと思ったものです。
このケースの課題は
- __cookbook開発のオーバーヘッドが無視できない__
	- cookbookの見通しが悪い。```_local```って何??
	- Atlasは、日本の日勤帯の時間は、ネットワーク帯域を絞っているのでboxのダウンロードが激重
- __attributeでパッケージ名を上書きしている___
	- OSが異なることが諸悪の根源
- __学習コストが結構かかる__
	- chefとServerspecは必要経費だとしても、Test Kitchen、Vagrant、Packer、Atlas.... 

ゴールデンイメージなVagrant boxを中心に添えたワークフローは、堅牢ですがミドルウェアのアップデートのスピードには着いてことが難しいというのが個人的な印象です。

## アフター

![after](./after.png)

1. exampleアプリケーションを稼働させるためのcookbook(example)を開発する
	- OpsWorksを利用しているので、リポジトリ直下にcookbookを配置する必要がある
1. Serverspecでテストコードを開発する
	- Test Kitchenにも対応するので、example_localに配置
1. Test Kitchenでexample cookbookの動作確認とテストを実施
	- Test Kitchenのドライバとして、Dockerを指定
	- kitchen-dockerでAmzon LinuxのイメージをPULLし、example cookbookをプロビジョニング
	- プロビジョニング:```bundle exec kitchen converge default-example```
	- テスト:```bundle exec kitchen verify default-example```
1. テストが通ったら、リモートリポジトリにPUSH
1. OpsWorksのスタックをsetupし、AWS上のインフラをプロビジョニング

cookbook開発だけでなく、アプリケーションのローカル開発環境の見直しも併せて行った(Dockerの導入)ので、かなりシンプルはワークフローとなりました。
メリットをまとめると
- __同一のOSに対してプロビジョニングし、テスト出来る__
	- AWS公式Dockerイメージに対してプロビジョニングした結果をテスト出来る
- __コードベースがシンプルになる__
	- example_local cookbookを破棄することが出来る
- __cookbook開発に注力出来る__
	- 以前に比べて、学習費用対効果が良い
	- 必要経費(chef、Serverspec)のみで開発出来る

## Tips

導入に際してのポイントをコードを交えて説明してみます。

### kitchen-dockerの導入

#### gemのインストール

kitchen-dockerというgemが必要になりますので、Gemfileに追記しましょう。
Kitchen-vagrantは不要となりますので、削除しました。
```Gemfile
source "https://rubygems.org"

gem "chef"
gem "berkshelf"
gem "test-kitchen"
gem "kitchen-docker"
```

Gemfileの編集が終わったらbundlerでインストールします。

```zsh
bundle install --path ./vendor/bundle
```

#### .Kitchen.ymlの編集

[kitchen-dockerのREADME]()を参考に、Test Kitchenの設定ファイルを編集していきます。


```yaml
provisioner:
  name: chef_solo
  environments_path: ../environments
  solo_rb:
    environment: local
```
Test KitchenのドライバとしてDockerを指定します。sudoするとdockerコマンドへのパスが見つからない場合が多いかと思うので、```use_sudo: false```と指定しておきます。
OpsWorks環境とは異なるattributeがあれば、environmentsとしてJSON形式のファイルにまとめておくことをおすすめします。

```yaml
platforms:
  - name: amazonlinux-201609
    driver_config:
      platform: rhel
      image: amazonlinux:2016.09
```
name要素でプロビジョニング対象のホストの名前を指定します。
Amazon Linuxは、Test Kitchenの分類ではRHEL系のディストリビューションとなるので、```platform: rhel```と指定します。
image要素で、Dockerレジストリから取得したいイメージを指定します。OpsWorks環境のOSバージョンと合わせるため、```image: amazonlinux:2016.09```と指定します。

```yaml
suites:
  - name: default
    run_list:
      - recipe[example::amazonlinux-local]
      - recipe[example::default]
    attributes:
```
テストしたいrecipeをリストします。
amazonlinux-localというレシピについては、後述します。

上記をまとめると、以下の様な.Kitchen.ymlとなります。

```yaml
---
driver:
  name: docker
  use_sudo: false

provisioner:
  name: chef_solo
  environments_path: ../environments
  solo_rb:
    environment: local

platforms:
  - name: amazonlinux-201609
    driver_config:
      platform: rhel
      image: amazonlinux:2016.09

suites:
  - name: default
    run_list:
      - recipe[example::amazonlinux-local]
      - recipe[example::default]
    attributes:
```

### amazonlinuxの導入

[amazonlinux](https://hub.docker.com/_/amazonlinux/)は、AWS公式のDockerイメージですが、AWS環境と全く同一のものではありません。あくまでも、可搬性の観点で現状取りうる手段の中で最もAWS環境で稼働するAmazon Linuxに近い仮想環境です。
そのため、AWS環境との差異(OSの初期設定含む)については、自分たちで対応する必要があります。
私達のチームでは、差異を埋めるrecipeを開発して対応しました。

```
# sysconfigファイル
default[:sysconfigs] = ['i18n', 'network']
```

システム設定ファイル(```/etc/sysconfig```)の差異をattributeにまとめ、templateとして配置しておきます。

```zsh
$ tree /path/to/example/example/templates/default/
/path/to/example/example/templates/default/
├── i18n.erb
├── config.erb
└── network.erb
```

```ruby
# ローカル開発環境(Amazon LinuxのDockerイメージ)でのレシピ開発のための各種ファイルの配置
node[:sysconfigs].each do |sysconfig|
  template "/etc/sysconfig/#{sysconfig}" do
    owner "root"
    group "root"
    mode 0644
    not_if { File.exist?("/etc/sysconfig/#{sysconfig}") }
  end
end
```
配置したtemplateを上記の様なrecipeでプロビジョニングしていきます。

```ruby
# SELinuxは明示的にDisable
include_recipe 'selinux'
template "/etc/selinux/config" do
  owner "root"
  group "root"
  mode 0644
  not_if { File.exist?("/etc/selinux/config") }
end
```
これは、テストのためのプロビジョニングです。
ServerspecでSELinuxの設定をテストしているのですが、```/etc/selinux/config```ファイルが必須の内容となっておりましたので、ローカル環境でのrecipe開発のためだけにプロビジョニングしています。

```ruby
# peclを利用するのでコンパイラをインストール
package "gcc"
```
gccがインストールされておりませんので、テスト対象のrecipeをプロビジョニングする前にインストールしておきます。

上記をまとめると、以下の様なrecipe(amazonlinux-local.rb)となります。
このrecipeをテスト対象のrecipe(example::default)の前にプロビジョニングします(前述の.kitchen.yml参照)。

```ruby
# Cookbook Name:: example
# Recipe:: amazonlinux-local

# ローカル開発環境(Amazon LinuxのDockerイメージ)でのレシピ開発のための各種ファイルの配置
node[:sysconfigs].each do |sysconfig|
  template "/etc/sysconfig/#{sysconfig}" do
    owner "root"
    group "root"
    mode 0644
    not_if { File.exist?("/etc/sysconfig/#{sysconfig}") }
  end
end

# SELinuxは明示的にDisable
include_recipe 'selinux'
template "/etc/selinux/config" do
  owner "root"
  group "root"
  mode 0644
  not_if { File.exist?("/etc/selinux/config") }
end

# peclを利用するのでコンパイラをインストール
package "gcc"
```

## Test Kitchen

kitchen-dockerの導入を紹介してきましたが、Test Kitchenのコマンド集をまとめておきます。

```
# cookbookのディレクトリに移動
$ cd /path/to/example/example

# インスタンスの確認
$ bundle exec kitchen list
Instance                    Driver  Provisioner  Verifier  Transport  Last Action  Last Error
default-amazonlinux-201609  Docker  ChefSolo     Busser    Ssh        Verified     <None>

# インスタンスを作成
$ bundle exec kitchen create default-amazonlinux-201609

# インスタンスのセットアップ
$ bundle exec kitchen converge default-amazonlinux-201609

# インスタンスへログイン(tty)
$ bundle exec kitchen login default-amazonlinux-201609
Last login: Thu Jan  5 20:23:42 2017 from 172.17.0.1
[kitchen@685fdd7c1349 ~]$
[kitchen@685fdd7c1349 ~]$ pwd
/home/kitchen

# テスト(serverspec)の実行
$ bundle exec kitchen verify default-amazonlinux-201609

# インスタンスの破棄
% bundle exec kitchen destroy default-amazonlinux-201609
```

## おわりに

OpsWorksを導入して以来、Amazon Linuxをどうテストするかという課題を抱えていましたが、AWS公式のDockerイメージが提供されたことにより、ローコストで可搬性のあるテストを実施出来るようになりました。
今後の課題は、プロビジョニングをCIすることを目標に必要なライブラリの導入や開発、変更を加えて行きたいと思います。
増え続けるcookbookや優秀なchefの育成に困っている開発チームの課題解決の一助となれば幸いです。

## 参考

- [開発環境のDocker化 その後 - ランサーズ（Lancers）エンジニアブログ](https://engineer.blog.lancers.jp/2016/12/%E9%96%8B%E7%99%BA%E7%92%B0%E5%A2%83%E3%81%AEdocker%E5%8C%96-%E3%81%9D%E3%81%AE%E5%BE%8C/)
- [kitchen-dockerでEC2の料金をかけず、高速にAmazon LinuxでのCookbookのテストを行う方法 | A Convenient Engineer's Note](http://blog.memolib.com/memo/942/)
- [kitchen-docker_cliというTest-Kitchen Driver Pluginを作った - Miscellaneous notes](http://marcy.hatenablog.com/entry/2014/12/25/020211)
- [Can't run kitchen-docker with Docker for Mac · Issue #207 · test-kitchen/kitchen-docker · GitHub](https://github.com/test-kitchen/kitchen-docker/issues/207)
