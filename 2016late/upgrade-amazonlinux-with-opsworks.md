# OpsWorks管理下のAmazon LinuxをCLIでアップグレードする

## はじめに

何かにつけ、GUIが苦手です。WEBブラウザを開くと、トンデモなバナーをクリックしてしまうし、そもそもトンデモなバナーが出来てきてしまうネットワーク広告が怖い。
OpsWorks(AWSマネジメントコンソール)では、広告は表示されないのですが、いつか広告モデルになるかもしれない...
そんな不安から日々のオペレーションは、[awscli](http://docs.aws.amazon.com/cli/latest/reference/index.html#cli-aws)を用いたCLIで完結したいというのがモチベーションです。
今回は、OpsWorksで唯一GUIなオペレーション(awscliにサブコマンドが用意されていない)となる```Upgrade Operating System```をCLIで実施するコネタを紹介します。

## バージョン

動作を確認している環境とそのバージョンは、以下の通りです。

- OpsWorks(リモート環境)
	- Cehf 11.10
	- Berkshelf 3.2.0
	- Amazon Linux 2016.03 → 2016.09
- Mac(ローカル環境)
	- aws-cli 1.11.28 Python/2.7.10 Darwin/15.6.0 botocore/1.4.85 

## Custom JSONの準備

下ごしらえとして、OpsWorksに渡すattributeをCustom JSONとしてJSONフォーマットのファイルにまとめておきます。
OpsWorksのビルドインcookbookであるdependenciesのattributeを上書きします。詳細は、後述します。

```
$ mkdir ./aws/opsworks && vi $_/upgrade-amazonlinux.json

{
  "dependencies": {
    "os_release_version": "2016.09",
    "allow_reboot": false,
    "upgrade_debs": true
  }
}

```

## コマンド

CLIでAmazon Linuxをアップグレードしていきます。  
幸せなchefになるために、独自の味付けをしない、つまり可能な限りビルドインのcookbookを利用するべきです。
```Upgrade Operating System```で実施しているのは、```update_dependencies```というdeployになり、```dependencies::update```レシピを実行しています。[該当のレシピ](https://github.com/aws/opsworks-cookbooks/blob/release-chef-11.10/dependencies/recipes/update.rb)を一読すると、```yum -y update```コマンドを実行していました。前述のCustom JSONは、```yum -y update```するために必要なattributeとなります。  

- ```os_release_version```は、アップグレードしたいAmazon Linuxのバージョンを指定します。
- ```allow_reboot``` は、パッケージのアップデートの後に再起動するか指定します。今回は、インスタンスの停止を明示的に実施しますので、```false```としておきます。
- ```upgrade_debs```は、一見Debianぽいですがパッケージをアップデートするか否かのフラグとして実装されてます。今回は、アップデートするので```true```としておきます。

```Upgrade Operating System```の正体を把握できたので、awscliで以下のような一連コマンドを実行していきます。

```
# 1. Stackのリビジョンを指定
$ aws opsworks --region us-east-1 update-stack --stack-id STACK_ID --custom-cookbooks-source "{\"Revision\":\"UgradeAmazonLinux\"}"

# 2. Stackで管理している全EC2インスタンスに対して、update_custom_cookbooksの実行(最新版cookbookを配置)
$ aws opsworks --region us-east-1 create-deployment --stack-id STACK_ID --command "{\"Name\":\"update_custom_cookbooks\"}"

# 3. opsworks agentのバージョンアップ(最新版を利用する)
$ aws opsworks --region us-east-1 update-stack --stack-id STACK_ID --agent-version LATEST

# 4. Custom JSONとレシピを指定して、全パッケージをアップデート
$ aws opsworks --region us-east-1 create-deployment --stack-id STACK_ID --instance-ids INSTANCE_ID01 INSTANCE_ID02 --command "{\"Name\":\"execute_recipes\",\"Args\":{\"recipes\":[\"dependencies::update\"]}}" --custom-json file://./aws/opsworks/upgrade-amazonlinux.json

# 5. EC2インスタンスの停止
$ aws opsworks --region us-east-1 stop-instance --instance-id INSTANCE_ID01
$ aws opsworks --region us-east-1 stop-instance --instance-id INSTANCE_ID02

# 6. OpsWorksで保持しているOSのバージョン情報を更新
$ aws opsworks --region us-east-1 update-instance --instance-id INSTANCE_ID01 --os "Amazon Linux 2016.09"

# 7. EC2インスタンスの起動
$ aws opsworks --region us-east-1 start-instance --instance-id INSTANCE_ID01
$ aws opsworks --region us-east-1 start-instance --instance-id INSTANCE_ID02
```

4でビルドインのcookbookにCustom JSONでattributeを渡し、全パッケージのアップデートを実施します。
5でEC2インスタンスを停止するのは、以下の2つの理由があります。

- OpsWorksが保持しているEC2インスタンスの情報を更新するためには、該当のEC2インスタンスを停止する必要がある
- OSアップグレード後は、```setup```ライフサイクルイベントを実施することを[推奨](http://docs.aws.amazon.com/opsworks/latest/userguide/workingstacks-commands.html)されている

```setup```ライフサイクルイベントは、7の起動時に実行されます。

## おわりに

AWSが広告モデルになったら嫌ですね。
Enjoy CLI!

## 参考

- [AWS OpsWorks とは - AWS OpsWorks](http://docs.aws.amazon.com/ja_jp/opsworks/latest/userguide/welcome.html)
