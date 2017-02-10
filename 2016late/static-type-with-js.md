# Vimmerが型を意識しながらJavaScript(ES2015)を書くためにやったこと

## はじめに

golangに出会ってから、静的型付けのメリットを改めて実感したVimmerがJavaScriptを書くために0.5人日くらい試行錯誤したお話です。

## 手順

### 用意するもの

- [Vim 8](http://vim-jp.org/blog/2016/09/13/vim8.0-features.html)
- [Node.js(v6.9.4)](https://nodejs.org/)
- [Flow](https://flowtype.org/)
- [ESLint](http://eslint.org/)
- [ALE](https://github.com/w0rp/ale)

### セットアップ

Vim及びNode.jsのセットアップは省略します。
また、Vimは8系を利用しますので、ご注意くだしあ。

#### JavaScript関連

[Flow](https://flowtype.org/)は、JavaScriptのタイプチェッカーで、Facebook社が主導しているプロジェクトです。
[flow-bin](https://www.npmjs.com/package/flow-bin)として、NPMで配信されていいますので、npmコマンドでインストールします。
グローバルにインストールしていますが、各自の環境に合わせてください。

```
$ npm install flow-bin -g
```

最終的には、AWS Lambda(Node.js)の関数としてデプロイしますので、ES2015で開発していきます。
リンターとして、ESLintをインストールします。

```
$ npm install eslint -g
```


#### ALE(Asynchronous Lint Engine)

ALEは、Vim8(NeoVim含む)のjobやtimmerの機能を使って非同期にリントしてくれる素敵なVim Pluginです。
~~面倒くさがり~~保守的な私は、[NeoBundle](https://github.com/Shougo/neobundle.vim)を利用していますので、以下の様にALEをインストールしました。

```
$ vi ~/.vimrc

" ALE(Asynchronous Lint Engine)
NeoBundle 'w0rp/ale'

:NeoBundleInstall
```

## 動作確認

今回は、JavaScript(ES2015)の学習も兼ねて、Flowのリポジトリをチェックアウトし、チュートリアル用のコードベースを利用します。

```
$ git clone git@github.com:facebook/flow.git && cd flow/example
```


### Flowによる型チェック

Flowに型をチェックしてもらうためには、以下の作業が必要となります。

- 設定ファイルを配置
- 対象のJavaScriptファイルに```@flow```をコメントとして記載

チュートリアルでは上記は設定済みですが、自分のプロジェクトでは対応する必要があるので、注意してください。

では、5つのチュートリアルが用意されているので、Vimで最初のものを開いていきます。

![flow](https://qiita-image-store.s3.amazonaws.com/0/27655/5c1d890c-0ea4-08c9-8b8d-aac6afffd19e.gif)

fooという関数は、number型を引数に乗算しますが、7行目でコールする際にstring型を渡しているので、4行目で型が異なる旨を指摘してくれてます。
コールする際に、number型を渡すように修正すると指摘事項がなくなり、クリアになっていますね :+1:

### ESLintによるリント

リントの設定ファイルを生成し、中身を確認してみます。

```
$ eslint --init
? How would you like to configure ESLint? Answer questions about your style
? Are you using ECMAScript 6 features? Yes
? Are you using ES6 modules? Yes
? Where will your code run? Node
? Do you use JSX? No
? What style of indentation do you use? Spaces
? What quotes do you use for strings? Double
? What line endings do you use? Unix
? Do you require semicolons? Yes
? What format do you want your config file to be in? JavaScript

$ vi -R .eslintrc.js

module.exports = {
    "env": {
        "es6": true,
        "node": true
    },
    "extends": "eslint:recommended",
    "parserOptions": {
        "sourceType": "module"
    },
    "rules": {
        "indent": [
            "error",
            4
        ],
        "linebreak-style": [
            "error",
            "unix"
        ],
        "quotes": [
            "error",
            "double"
        ],
        "semi": [
            "error",
            "always"
        ]
    }
};
```

先程と同様のJavaScriptファイルをVimで開いてみます。

![tty.gif](https://qiita-image-store.s3.amazonaws.com/0/27655/d248016d-c893-7267-e500-d6dea588fe0d.gif)

インデントのスペース量で指摘されてますので、4スペースに修正しました。
クリアになりましたね :+1:


## おわりに

VimmerがJavaScript(ES2015)でプログラミングするに際に、少しだけ幸せになれる(でも重要な)コネタを紹介してみました。
JavaScriptはまだ習得している途中ですが、奥深い言語ですので楽しみながらプログラミング出来ればと思います。

## 参考

- [VimでESLintとFlowを使うためにNeomakeからALEに乗り換える - Qiita](http://qiita.com/zaki-yama/items/6bcc24469d06acdf8643)
- [O'Reilly Japan - 初めてのJavaScript](https://www.oreilly.co.jp/books/9784873113227/)
