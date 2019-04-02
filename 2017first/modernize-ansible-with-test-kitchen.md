# Test KitchenでAnsible Playbook開発を近代化する

TechLeadな会でもアナウンス致しましたが、インフラ部ではAnsible Playbook開発の近代化に取り組んでいます。
各部のプロジェクトチームでもAnsible Playbookを開発、運用されているかと思いますので、導入を検討したいチーム、エンジニアは北田までお声がけ頂ければと思います。

## 動機と経緯

プロビジョニングとその周辺ツール、コードは、その性質上、秘伝のタレ化しやすい分野です。
[Infrastructure as Code](https://ja.wikipedia.org/wiki/Infrastructure_as_Code)が常識的となった現在でも、プロビジョニングは北田さんで...というような状況は発生しているかと思います。
その原因のひとつに、プロビジョニングというプロセスがテスタブルではないからなのでは??という課題意識が前々あり、Ansible Playbookを開発する機会があれば、取り組んでみようと機会を伺っておりました。
今回、[shiromane](https://github.com/mediba-system/shiromane)プロジェクトでAnsible Playbookを開発する機会を得たので、[Test Kitchen](http://kitchen.ci/)を用いてテスタブルに、近代化してみました。

## ワークフロー

(0. 下ごしらえ)
1. Ansible Playbookを書く
1. テストコードを書く
1. Test Kitchenを用いてテスト実行
1. テストをパスするまで、Ansible Playbookを開発、修正する
1. プロジェクトで利用する環境をプロビジョニング

実際のプロジェクトでは、1~4がイテレーティブなプロセスになるかと思います。

## 下ごしらえ

### AWSアカウント
(北田が)管理しやすいAWSアカウント `techtools`を選定しました。
コスト管理は、インフラ部DevOpsチームで実施しております。 

### IAMユーザ
セキュリティの観点で悩みどころでしたが、必要最低限の権限に絞り、`test-kitchen`というIAMユーザを`techtools`に新設しました。
アクセスキーとシークレットアクセスキーは、導入の際にお伝えします。

### EC2のキーペア
セキュリティの観点で悩みどころでしたが、`test-kitchen`というキーペアを`techtools`アカウントにて発行しました。
秘密鍵は、導入の際にお伝えします。

### ライブラリの導入
rbenv等を用いて、Rubyのバージョンを管理した上で、BundlerとTest Kitchen、その他必要なライブラリをインストールしていきます。

```
$ cd /path/to/yourRepo
$ vi Gemfile

source 'https://rubygems.org'

gem 'test-kitchen'
gem 'kitchen-ansiblepush'
gem 'kitchen-ec2'
gem 'rake'
gem 'serverspec'

# ライブラリをインストール
$ bundle install --path ./vendor/bundle
```

### Test Kitchenのセットアップ
Test Kitchenの設定ファイルを記述します。
(`driver`の`subnet_id`項目については、変更する可能性があります)

```yaml
# AWS環境techtoolsの設定
driver:
  name: ec2
  region: ap-northeast-1
  shared_credentials_profile: test-kitchen
  aws_ssh_key_id: test-kitchen
  associate_public_ip: true
  interface: public
  instance_type: t2.nano
  security_group_ids:
    - sg-dc1c03bb
  subnet_id: subnet-9cbc9dea
  block_device_mappings:
    - device_name: /dev/xvda
      ebs:
        volume_type: standard
        volume_size: 8
        delete_on_termination: true
  tags:
    Name: Created-by-test-kitchen-at-<%= Time.now.strftime('%Y%m%d') %>

# プロビジョニングツールの設定
# test.ymlでstage変数を設定注入
provisioner:
  name: ansible_push
  playbook: stns-client.yml
  extra_vars: "@test.yml"
  sudo: true
  chef_bootstrap_url: false

# ディストリビューションを指定
# 特定のディストリビューションは、最新のバージョンの公式AMIを指定したことになる
# AMI IDを指定することも可能
platforms:
  - name: debian
  - name: centos
  - name: amazonlinux-201703
    driver:
      image_id: ami-923d12f5
    transport:
      username: ec2-user

# Chefは利用しないが記述は必要
suites:
  - name: default
    run_list:
    attributes:

# EC2へのSSH接続設定
transport:
  ssh_key: ~/.ssh/test-kitchen

# テスト(Serverspec)実行
verifier:
  name: shell
  command: bundle exec rspec -f d
```

## Playbook開発
[Best Practices — Ansible Documentation](http://docs.ansible.com/ansible/playbooks_best_practices.html)
Playbookの構成は、公式ドキュメントのベストプラクティスやチーム現状を尊重、踏襲する方針がベターかと思います。
shiromaneプロジェクトでは、サポートする環境が単純(検証 or 商用)だったので、`group_vars`は環境別に用意しました。
Test Kitchenでのテスト実行の際は、インベントリーファイルが動的に生成されるので、`inventories`配下に`test`は不要です。

```
$ tree /path/to/lib_shiromane/ansible
ansible
├── README.md
├── .kitchen.yml # Test Kitchen設定ファイル
├── group_vars # 環境(ステージ)別にグループ化
│   ├── develop.yml
│   └── test.yml # Test Kitchenで利用する変数
├── inventories # 環境別にホスト情報をグループ化
│   └── develop
├── roles # 目的別にタスクをグループ化
│   ├── stns-client
│   │   ├── files
│   │   │   ├── ca.pem # 中間証明書(バージョン管理対象外)
│   │   │   ├── client.crt # クライアント証明書 SSL証明書(バージョン管理対象外)
│   │   │   └── client.key # クライアント証明書 秘密鍵(バージョン管理対象外)
│   │   ├── handlers
│   │   │   └── main.yml
│   │   ├── tasks
│   │   │   └── main.yml
│   │   └── templates
│   │       ├── libnss_stns.conf.j2
│   │       └── nscd.conf.j2
│   └── stns-common
│       ├── files
│       │   ├── stns.list
│       │   └── stns.repo
│       └── tasks
│           └── main.yml
├── spec # rspc(Serverspec)テストコードを格納
│   ├── spec_helper.rb
│   ├── nscd # nscdに関するテストコード
│   │   └── nscd_spec.rb
│   ├── nss # NSSに関するテストコード
│   │   ├── libnss_stns_spec.rb
│   │   └── nsswitch-conf_spec.rb
│   ├── sshd # sshdに関するテストコード
│   │   └── sshd_spec.rb
│   └── yum # yumに関するテストコード
│       └── yum_spec.rb
├── stns-client.yml # 実行するPlaybook
└── test.yml # Test Kitchenで利用する変数を定義
```

## テストコード

[RSpec](http://rspec.info/)のAPIで比較的簡単にプロビジョニングのテストが出来る[Serverspec](http://serverspec.org/)でテストコードを開発します。
ポイントは、ディレクトリ構成をRSpecのパターンに合わせる点です。

- 例:nscdに関するテストコードを書きたい

```
# Playbookと同一ディレクトリにmkdir
% mkdir -p spec/nscd && cd $_

# nscdのプロビジョニングについてテストしたい内容を記述
% vi nscd_spec.rb

# RSpecのディレクトリパターンを合致させる
% cd ../../
% tree spec/nscd
spec/nscd
└── nscd_spec.rb
```

### specディレクトリのパスを指定するパターン

Test Kitchenの設定ファイル(`.kitchen.yml`)で`rspec`コマンドを設定、実行していますので、`I`オプションと`--default-path`オプションでディレクトリのパスを指定して出来ます。

- specディレクトリと`.kitchen.yml`の位置関係

```
% cd /path/to/YOUR_REPO
% tree ansible
ansible
├── .kitchen.yml

% tree hoge
hoge
└── spec
    └── yum
		└── yum_spec.rb
```

- `.kitchen.yml`を編集し、`rspec`コマンドのオプション指定

```
# テスト(Serverspeb)実行
verifier:
  name: shell
  command: bundle exec rspec -f d -I ../hoge/spec --default-path ../hoge
```

## テスト

### 流れ

Test Kitchenで以下のタスクを担当してくれます。

(0. EC2インスタンスの停止、破棄)
1. プロビジョニング対象のEC2インスタンスを起動、SSH接続が可能となる
1. プロビジョニング(Ansible Playbookの実行)
1. テスト実行(Serverspecのテストコードを実行)
1. EC2インスタンスを停止(プロビジョニングもしくは、テストに失敗した場合は、起動した状態)

### コマンド

```
$ bundle exec kitchen test
```

### 効率化
テストプロセスの並行実行も出来ます。

```
$ bundle exec kitchen test -c 3
```

## 今後の展望

- Travisおじさんに任せたい
- インスタンス管理を安全に行いたい
- AWSのリソース構築とそのテストもプロセスに組み込みたい

## 問い合わせ先

インフラ部DevOpsチーム 北田まで
