# yii2で高性能なWebAPIを実現する

この記事は[mediba advent calendar 2016](http://qiita.com/advent-calendar/2016/mediba) 4日目です。

## はじめに

とあるプロジェクトで[yii](http://www.yiiframework.com/)のバージョン2系で実装したWebAPIを開発、運用しています。  
本エントリーでは、yiiのバージョン2系(以下yii2)でWebAPIを開発する際のtipsをまとめました。  
ご査収ください。

### 想定読者

相変わらず、[ターゲットが狭い](http://ceblog.mediba.jp/post/144441036997/vim-de-go)です。

- PHPer
- yii2で実装されたWEBアプリケーションを開発、運用している
- 高性能なWebAPIを要求されている

### yii2とは

PHPのWeb Application Frameworkです。[多機能](https://en.wikipedia.org/wiki/Comparison_of_web_frameworks#PHP_2)で[パフォーマンスも悪くない](https://github.com/kenjis/php-framework-benchmark#hello-world-benchmark)とバランスの取れたフレームワークと言えるでしょう。  
また、弊社の開発現場では実績の多いフレームワークです。

### 高性能とは

このエントリーでは、以下の通りとします。

- 複数の機能を提供しており、その拡張性を担保出来ている
- 費用対効果の高い性能を担保出来ている

### WebAPIとは

JSON over HTTPなWebアプリケーションとします。

## 3行で

プロジェクトを経て気づいたことを3行でまとめます。

- RESTfulに設計する
- ガイドを熟読する
- 計測する

出来たことも出来なかったこともあったプロジェクトでしたが、気づきの多い良いプロジェクトだっと振り返って思います。  
気づきを得たいエンジニアは、[奮ってご応募下さい。](http://www.mediba.jp/recruit/)


## RESTfulに設計する

フレームワークが課題の解決に寄与できる割合というのは、想像よりも多くありません。可能な限り[RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer)に設計しましょう。多くのフレームワークは、RESTfulな設計を歓迎しています。yii2も例に漏れず、RESTfulな設計に応えてくれます。  
...   
とはいえ、現実のプロジェクトでRESTfulな設計を実現することは難しいものです。現実と向き合い、真摯な姿勢で課題の解決を図るべきでしょう。  
今回のプロジェクトでは、インターフェイスの変更も実施しましたが、それはRESTfulとは呼べないものでした。RESTfulに設計が出来れば、生産性は高くなるなと気づいた程度です。  
yii2のRESTfulについては、ガイドに詳しくまとまっています。ActiveRecordとの密な連携を提供していますので、RESTfulに設計したかったと悔しく思います。

- [Resources - RESTful Web Services - The Definitive Guide to Yii 2.0](http://www.yiiframework.com/doc-2.0/guide-rest-resources.html)

## ガイドを熟読する

開発チームには、迷ったら[Stack Overflow](http://stackoverflow.com/questions/tagged/yii)や他チームの実装ではなく、まずはガイドを読もうと口を酸っぱくして言い続けました。

- [The Definitive Guide to Yii 2.0](http://www.yiiframework.com/doc-2.0/guide-README.html)

[日本語訳](https://github.com/yiisoft/yii2/tree/master/docs/guide-ja)もあります。  
今回のプロジェクトで、参考になった章を掘り下げてみます。

### [Quick Start](http://www.yiiframework.com/doc-2.0/guide-rest-quick-start.html)

まずは、この章を読んでyiiが提供しているWebAPIな機能の概要を把握しましょう。通読すれば、実装のスリム化を図れる場合がほとんどでしょう。

### [Controllers](http://www.yiiframework.com/doc-2.0/guide-rest-controllers.html)

コントローラの実装について、ガイドしてくれます。RESTfulに設計出来なかった場合でも、[yii\rest\Controller](http://www.yiiframework.com/doc-2.0/yii-rest-controller.html)を継承しておくと幸せになれます。WebAPIで重要視されるフィルターの機能が提供されているからです。

```php
/**
 * behaviors
 * エンドポイントの振る舞いを定義
 *
 * @return array behaviors
 */
public function behaviors()
{
    $behaviors = parent::behaviors();
    $behaviors['contentNegotiator'] = [
        'class' => ContentNegotiator::className(),
        'formats' => [
            'application/json' => Response::FORMAT_JSON,
        ],
    ];
    $behaviors['verbs'] = [
        'class' => VerbFilter::className(),
        'actions' => [
            'index'  => ['get'],
        ]
    ];
    return $behaviors;
}
```

上記の様に、コントローラクラス内で[yii\base\Behavior](http://www.yiiframework.com/doc-2.0/yii-base-behavior.html)[]型のbehaviors関数をオーバーライドし、エンドポイントの振る舞いを定義します。エンドポイントの振る舞いを明示的に指定し、見通しを確保することが出来ます。
上記のコードでは、以下を指定しています。
- レスポンスボディのフォーマットをJSONとする
- indexアクション(indexAction関数)へのリクエストは、GETメソッドのみとする

また、ユーザ認証やレートリミット等のアクセス制限も明示的に指定することも出来ます。Webサーバとの担当範囲の問題もありますが、アプリケーションの開発者としてはアクセス制限をコードで管理出来ることのメリットはあると思います。

### [Error Handling](http://www.yiiframework.com/doc-2.0/guide-rest-error-handling.html)

エラーハンドリングは、ハマりどころでした...前述のBehaviorと同じくらいカジュアルに指定出来るといいのですが、アプリケーションコンポネントをグローバルに変更するしかありませんでした。

```php
return [
    // 中略
	],
	'components' => [
		// アプリケーション全体のレスポンスを設定
		'response' => [
			// レスポンスボディのフォーマットはJSONとする
			'format' => yii\web\Response::FORMAT_JSON,
			'charset' => 'UTF-8',
			// 異常系のレスポンスは、レスポンスボディを空配列とする
			'on beforeSend' => function($event) {
				$response = $event->sender;
				if ($response->statusCode != 200) {
					$response->data = [];
				}
			},
		],
    // 以下略
```

yii2のコアなアプリケーションコンポーネントとして、[errorHandler](http://www.yiiframework.com/doc-2.0/yii-base-errorhandler.html)があります。

```php
'errorHandler' => [
    'errorAction' => 'site/error',
],
```

このコンポーネントの[errorAction](http://www.yiiframework.com/doc-2.0/guide-runtime-handling-errors.html#using-error-actions)プロパティにアクション(コントローラで実装した関数)をしてすれば、エラー時のレスポンスを定義出来ると見込んでいたのですが、結果は上述の通り、[response](http://www.yiiframework.com/doc-2.0/yii-web-response.html)コンポーネントのbeforeSendイベントをオーバーライドする必要がありました。  
今回のプロジェクトでは、この実装で事足りましたが、エンドポイント(リソース)毎にエラーレスポンスが異なる場合は、そうもいきません。引き続き、調査をしたいと思います。

## 計測する

これまでは、機能とその拡張の担保の観点で見てきましたが、今回のプロジェクトでは性能管理(向上あるいは維持)もスコープにしておりました。
性能管理は、とにもかくにも計測することには始まりません。[Rob Pike先生のお言葉](http://users.ece.utexas.edu/~adnan/pike.html)ですね。我々のプロジェクトでは、以下について計測しました。  

- プロファイリング
    - [Xhprof](http://php.net/manual/en/book.xhprof.php)を用いてアプリケーションをプロファイリングします。ボトルネックとなりうるコンポーネントや関数を把握しておきます。

- 性能試験
    - [gatling](http://gatling.io/#/)を用いてWebAPIに対し一定の負荷を掛け、秒間リクエスト数を把握します。この際の条件は後の性能試験でも利用するので、wiki等にまとめておきます。

今回のプロジェクトで性能(秒間リクエスト数)を1.7倍程度上げることが出来ました。パフォーマンス・チューニングの勘所もお伝えできるとよいのですが、眠くなってきたので、またの機会に。結構地道なアプローチでしたが、結果を出すことが出来て、良かったです。

計測の理想は、改修の度に(コミットが発生したら)実施することなのでしょうが、今回はそこまで至りませんでした。今後の課題としたいです。  
また、gatlingを用いた性能試験については、[次回](http://qiita.com/advent-calendar/2016/mediba)に譲ります。[TravisCI](https://travis-ci.org/)を利用している開発チームは少し幸せになれるかもしれません。

## おわりに

大仰なタイトルの割に、お伝えできることが小振りなのは[いつものこと](http://qiita.com/hikarut/items/68d102489a5e936467bb#phpアプリケーションに関する12の事実)ですが、プロジェクトを通して、yii2というフレームワークときちんと向き合うことが出来ましたので、その一部をお伝えしてきました。WebAPIを高いクオリティで開発するには、まだまだ課題が多いのは事実ですが、フレームワークを武器に継続的に課題に取り組んでいきたいと思います。  
また、yii2を用いた開発プロジェクトの生産性に少しでも寄与できていたら、とても光栄です。  
明日は、制作部平尾さんの「JavaScript初級者がVue.jsで幸せになれたお話」です。


