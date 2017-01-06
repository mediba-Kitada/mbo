# ローカル環境でAmzon Linuxをテストする

## はじめに

可搬性(ポータビリティ)に立ちふさがる密林、Amazon Linux。担当プロジェクトでもご多分に漏れず、Amazon LinuxをOSとしたシステムを組んでいる。
chefを用いて構成管理(商用環境は、[OpsWorks](https://aws.amazon.com/jp/opsworks/))、[Test Kitchen](http://kitchen.ci/) + [Serverspec](http://serverspec.org/)でインフラテストしているがOSはCentOS6.7(packerで生成したVagrant box)...OSが違うのに、何が構成管理か、何が可搬性か。  
[AWSお墨付きのDockerイメージ](https://hub.docker.com/_/amazonlinux/)もリリースされたことですし、Dockerを武器に可搬性を手に入れることが出来ると思い立ったが吉日。
1営業日くらいで可搬性を手に入れることが出来ました :dragon_face:

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
	- プロビジョニング:```bundle exec kitchen converge default-example```
	- テスト:```bundle exec kitchen verify default-example```
1. テストが通ったら、リモートリポジトリにPUSHし、OpsWorksのスタックをsetup

cookbook開発だけでなく、アプリケーションのローカル開発環境の見直しも併せて行った(Dockerの導入)ので、かなりシンプルはワークフローとなりました。
メリットをまとめると
- __同一のOSに対してプロビジョニングし、テスト出来る__
	- 公式Dockerイメージに対してプロビジョニングした結果をテスト出来る
- __コードベースがシンプルになる__
	- example_local cookbookを破棄することが出来る
- __cookbook開発に注力出来る__
	- 学習費用対効果が良い
	- 必要経費(chef、Serverspec)のみで開発出来る

## Tips

### kitchen-dockerの導入

### amazonlinuxの導入

## おわりに

## 参考
