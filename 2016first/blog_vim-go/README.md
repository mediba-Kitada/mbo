## はじめに

無用な混乱を避けるため、注意事項及び前提条件を設けました。
はじめにお読みください。

### ソースコード

本稿で利用するソースコードは、GitHubの公開リポジトリに格納してあります。
断りがない場合、カレントディレクトリは、以下の通りです。
```zsh
$ pwd
/path/to/mediba-kitada/mbo/2016first/blog_vim-go
```

### 省略したこと

- Vim及びGo言語の概要については省略しています。
    - 参考にした書籍
        - Software Design 2016年5月号
        - スターティングGo言語
- 以下のツールを利用しますが、概要については省略します。
    - direnv
    - glide
    - neobundle
    - neocomplete

### サポートする環境

- 以下のVimをサポートしています。

|ディストリビューション|バージョン|備考|
|---|---|---|

- 利用するVim pluginは、以下の通りです。

|名称|バージョン|備考|
|---|---|---|
|vim-go|||
|neocomplete|||

- 以下のGo言語をサポートしています。
    - 1.1以上

- 以下のshellをサポートしています。
    - zsh

## 想定読者

上記を踏まえ本稿の想定読者は、以下の通りです。

- Go言語のチュートリアルをやってみて、何かを見出した
- __Go言語と長くお付き合いしたい一般的なVimmer__

私です。

## 導入

[vim-go](https://github.com/fatih/vim-go)をインストールしていきます。  
vim-goは、Go言語向けVim pluginで他のpluginに比べ[圧倒的な支持](https://github.com/search?o=desc&q=vim+go+&ref=searchresults&s=stars&type=Repositories&utf8=%E2%9C%93)を得ています。
豊富な機能を提供してくれるvim-goですが、他のpluginを公式にサポートしており、機能の提供の仕方がunix的で好印象です。  
neobundleを用いて以下の様に導入していきましょう。

```zsh
$ vi ~/.vimrc

" golang
NeoBundle 'fatih/vim-go'

# コマンドモードで.vimrcをリロード
:<C-u>source $MYVIMRC

# コマンドモードでNeoBundleInstall
:NeoBundleInstall
```

試運転をしてみましょう。  
フォーマットが出鱈目なエントリーポイントを用意し、セーブしてみます。

```zsh
# フォーマットが出鱈目なエントリーポイントを開く
$ vi src/main.go
```

vim-goの```:GoFmt```コマンドでフォーマットされていることがわかります。

## カスタマイズ

### ツール群のインストール

vim-goは、```:GoInstallBinaries```コマンドを用いてGo言語での開発に便利なツール群をインストールしてくれます。ただし、インストールパスの解決には、```$GOPATH/bin```又は```$GOBIN```が参照されます。この際、問題になるのは```$GOPATH```を用いたパスの解決です。envdirは、```$GOPATH```問題をカジュアルに解決することが出来ますが、パッケージが混在する問題が生じます。

```
- プロジェクト固有のパッケージ
- vim-go等で利用する共有するパッケージ
```

この問題に対処するためにシンプルにshellの設定ファイルを用いて共通のGo言語環境を用意してしまいましょう。

```zsh
# 共有のGo言語環境の各ディレクトリを作成
$ mkdir -p {$HOME/go,$HOME/go/bin,$HOME/go/pkg,$HOME/go/src}

# shellの設定ファイルにGOPATHを設定
$ vi $HOME/.zshrc

export $GOPATH="$HOME/go"

# shellの設定を更新
$ source $HOME/.zshrc
```

次に```:GoInstallBinaries```コマンドでツール群をshellの設定ファイルで指定したパスにインストールします。

```zsh
# プロジェクトのルートディレクトリ(.envrcが格納されているディレクトリ)から離れる
$ cd $HOME
# Vimを起動しコマンドモードで GoInstallBinaries を実行
$ vi

:GoInstallBinaries

# ツール群が共通のGo言語環境のパスにインストールされていることを確認する
$ tree $GOPATH/bin
/Users/mediba-kitada/go/bin
├── asmfmt
├── errcheck
├── gocode
├── gogetdoc
├── goimports
├── golint
├── gometalinter
├── gorename
├── gotags
├── guru
└── motion
```

エントリーポイントを用いて、ツール群の試運転をしてみましょう。```:GoImports```は魔法ぽくて好きです。


### vim-goの設定

デフォルトではシンタックスハイライトがオフになっていますので、設定ファイルを用いて設定をカスタマイズしていきましょう。```:GoFmt```コマンドを実行した際の挙動もカスタマイズしてみます。

```zsh
$ vi $HOME/.vimrc
" vim-go
let g:go_hightlight_functions = 1
let g:go_hightlight_methods = 1
let g:go_hightlight_structs = 1
let g:go_hightlight_interfaces = 1
" :GoImportsをフォーマットの際に自動実行
let g:go_fmt_command = "goimports"
```

neocompleteをインストール、設定していれば自動補完もかっこ良くなります。


## 課題

本稿では対応出来ていないvim-goのカスタマイズについて簡単に触れておきます。

- Tagbar
    - ctagsを適切に設定することで関数の一覧バッファをサイドビューとして表示し、タグジャンプにも利用出来ます。
- Neovim対応
    - vim-go開発者のfaithさんもNeovimを利用しだしたとのことですが、まだうまく動いておらず、ベータとして提供している様です。
        - そもそもNeovimをキャッチアップ出来ていないので、neobundleの開発が終了するまでにはチャレンジしたいと思います。

## おわりに

Vimmer向けにGo言語開発環境の実践的な内容をまとめてみました。ここでは触れていませんが、テストの実行コマンドも提供されており、vim-goさえあればVimmerのGo言語開発は何とかなるのではと見込んでおります。これは、シンプルな規約に沿ってデザインされたGo言語の恩恵のひとつなのかもしれません。  
弊社では、~~Vimmer~~バックエンド(サーバサイド)エンジニアを[鋭意募集](https://rec-log.jp/site/jobVw.aspx?Oy24gBu7IDWaPGodjJQfOM2igOulIRWoaUoqCXQtO02wg2uzI5WCa8oECbQHOe2KgguNIjWPamoSCpQVOr2Yguu1IxW3aAo6CDQ9OF2cgIueeL4hGOtkKQ1n)しています。Go言語の他にもモダンな技術の導入には前向きに取り組んでおりますので、腕自慢の方はぜひご応募ください。
