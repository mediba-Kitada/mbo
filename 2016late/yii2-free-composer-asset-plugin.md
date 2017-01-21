# yii2をcomposer-asset-pluginから自由にする

## はじめに

PHPのアプリケーションフレームワーク[yii(以下yii2)](http://www.yiiframework.com/)は、FastでSecureでプロユースですが、ビルドに一手間掛かるのが玉に瑕です。
その手順は、[公式ドキュメント](http://www.yiiframework.com/doc-2.0/index.html)のファーストビューに飛び込んできてしまうほどです。
今回は、yii2をFastにビルドするコネタを紹介します。

## 課題

ビルドの一手間を具体的に確認してみます。

>[公式ドキュメント](http://www.yiiframework.com/doc-2.0/index.html)より
```
composer global require "fxp/composer-asset-plugin:^1.2.0"
composer require yiisoft/yii2
```

composerを用いて、2つのコンポーネントをインストールしています。2行目は、我らがyii2の本体となります。
1行目でインストールしている```fxp/composer-asset-plugin```の[リポジトリ](https://github.com/fxpio/composer-asset-plugin)を一読してみると

>This plugin works by transposing package information from NPM or Bower to a compatible version for Composer. This allows you to manage asset dependencies in a PHP based project very easily.

とあります。
NPMやBowerで管理されているライブラリをcomposer互換に整理してくれるツールの様です。yii2では、このツールをフロントエンドライブラリのパッケージマネージャとして利用するわけですが、以下の問題が生じます。

1. グローバルにインストールする
1. フロントエンドを必要としないアプリケーション(WebAPI等)には不要
1. インストールの所要時間が長い

1は言わずもがなですが、サーバの環境は可能な限りクリーンにしておきたいものです。  
2については、ビューを提供するWEBアプリケーションであれば必要コストなので問題ありませんが、フロントエンドを必要としないWebAPI(JSON over HTTP)なアプリケーションの場合は、面倒なだけです。  
3については、後ほど記載しますが、所要時間が想像以上に長く、yii2の価値を毀損しているのではないかと思うほどです。  

私の所属するチームでは、WebAPIの開発も担当しています。
実現したいことは、__```fxp/composer-asset-plugin```をインストールすることなくyii2で実装したアプリケーションをビルドする__ことです。

## 方法

簡単です。同じような課題意識を持った先輩の肩に乗るだけです。

- [cebe/assetfree-yii2](https://github.com/cebe/assetfree-yii2)

composer.jsonに上記のコンポーネントを追加します。注意点としては、```cebe/assetfree-yii2```とyii2本体のバージョンを合わせることです。

- yii2のバージョンを2.0.7とする場合
```
"yiisoft/yii2": "2.0.7",
"cebe/assetfree-yii2": "2.0.7"
```

追加後は、```composer install```でインストールしていきます。事前に```composer global require "fxp/composer-asset-plugin:^1.2.0"```コマンドを実行する必要はありません。

## 成果

私の所属するチームでは、OpsWorksを利用しており、アプリケーションのデプロイもOpsWorksが担当します。
```cebe/assetfree-yii2``` 導入前後デプロイの所要時間をOpsWorksのログファイルから拾ってみました。

- before : 2分37秒
```
[2017-01-21T01:59:56+09:00] INFO: Processing script[install_composer] action run (/srv/www/apollo_app_stg/releases/20170120165934/deploy/before_restart.rb line 24)
[2017-01-21T02:02:33+09:00] INFO: script[install_composer] ran successfully
```

- after : 1分10秒
```
[2017-01-21T01:54:30+09:00] INFO: Processing script[install_composer] action run (/srv/www/apollo_app_stg/releases/20170120165407/deploy/before_restart.rb line 24)
[2017-01-21T01:55:40+09:00] INFO: script[install_composer] ran successfully
```

所要時間が半分以下になりました。デプロイ時のコマンドも省略することが出来ました。

## おわりに

yii2をFastにビルドする方法を紹介致しました。次期バージョンでは、```fxp/composer-asset-plugin```の利用がオプションになるかもとのことですが、当面はお付き合いする必要がありそうです。
Fastといえば、PHP7の[事例](http://www.yiiframework.com/forum/index.php/topic/63662-yii2-and-php-7/page__view__findpost__p__299411)もぼちぼち挙がっている様です。こちらも機会を伺ってチャレンジしてみたいと思います。

## 参考

- [Yii 2.0.x で Composer Asset Plugin を使わないでアプリを作っていく方法 - Qiita](http://qiita.com/livejam_db/items/70d674d0d735038ef93f)
