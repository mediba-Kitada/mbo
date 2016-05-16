# vimでgoする

メディアシステム開発部の北田です。
auスマートパスのサーバサイド開発を担当しております。

担当サービスのQCDと個人的な問題意識から、[Go言語](https://golang.org/)の導入に取り組んでいます。
今回は、[Vimmer](https://www.google.co.jp/trends/explore#q=vimmer%2C%20emacser&cmpt=q&tz=Etc%2FGMT-9)がGo言語と向き合う際のtipsについてまとめました。

## はじめに

無用な混乱を避けるため、注意事項及び前提条件を設けました。
はじめにお読みください。

### ソースコード

本稿で利用するソースコードは、[GitHubの公開リポジトリ](https://github.com/mediba-Kitada/mbo/tree/master/2016first/blog_vim-go)に格納してあります。
断りがない場合、カレントディレクトリは、以下の通りです。

    $ pwd
    /path/to/mediba-kitada/mbo/2016first/blog_vim-go

### 省略したこと

- Vim及びGo言語の概要については省略しています。
    - 参考にした書籍
        - [Software Design 2016年5月号](http://gihyo.jp/magazine/SD/archive/2016/201605)
        - [スターティングGo言語](http://www.shoeisha.co.jp/book/detail/9784798142418)
- 以下のツールを利用しますが、概要については省略します。
    - [direnv](http://direnv.net/)
    - [neobundle](https://github.com/Shougo/neobundle.vim)
    - [neocomplete](https://github.com/Shougo/neocomplete.vim)

### サポートする環境

- 以下のVimをサポートしています。

|ディストリビューション|バージョン|備考|
|---|---|---|
|MacVim|7.4|筆者環境は、luaサポートをonにしてあります|

- 利用するVim pluginは、以下の通りです。

|名称|バージョン|備考|
|---|---|---|
|vim-go|master|-|

- 以下のGo言語をサポートしています。
    - 1.1以上

- 以下のshellをサポートしています。
    - zsh

## 想定読者

上記を踏まえ本稿の想定読者は、以下の通りです。

- __[Go言語のチュートリアル](https://go-tour-jp.appspot.com/welcome/1)をやってみて、何かを見出した__
- __Go言語と長くお付き合いしたい一般的なVimmer__

私です。

## 導入

[vim-go](https://github.com/fatih/vim-go)をインストールしていきます。  
vim-goは、Go言語向けVim pluginで他のpluginに比べ[圧倒的な支持](https://github.com/search?o=desc&q=vim+go+&ref=searchresults&s=stars&type=Repositories&utf8=%E2%9C%93)を得ています。
neobundleを用いて以下の様に導入していきましょう。

    $ vi ~/.vimrc  
    
    " golang  
    NeoBundle 'fatih/vim-go'  
    
    " コマンドモードで.vimrcをリロード  
    :<C-u>source $MYVIMRC  
    
    " コマンドモードでNeoBundleInstall  
    :NeoBundleInstall  

試運転をしてみましょう。  
フォーマットが出鱈目なエントリーポイントを用意し、セーブしてみます。


    # フォーマットが出鱈目なエントリーポイントを開く  
    $ vi src/main/main.go  
    
    " コマンドモードでGoFmt  
    :GoFmt  

[![GoFmt](https://i.gyazo.com/545cba71b9999d8dd20eaddfdddea3c9.gif)](https://gyazo.com/545cba71b9999d8dd20eaddfdddea3c9)

vim-goの    :GoFmt    コマンドでフォーマットされていることがわかります。

## カスタマイズ

### ツール群のインストール

vim-goは/、    :GoInstallBinaries    コマンドを用いてGo言語での開発に便利なツール群をインストールしてくれます。ただし、インストールパスの解決には、    $GOPATH/bin    又は    $GOBIN    が参照されます。この際、問題になるのは    $GOPATH    を用いたパスの解決です。envdirは、    $GOPATH    問題をカジュアルに解決することが出来ますが、パッケージが混在する問題が生じます。

    # プロジェクト固有のパッケージ  
    $ tree ./  
    └── src  
        ├── animals  
        │   ├── elephant.go  
        │   └── monkey.go  
        └── main  
            └── main.go  
    # vim-go等で利用する共有するパッケージ(srcディレクトリの一部を抜粋)  
    $ tree $GOPATH/src  
    ├── github.com  
    │   ├── alecthomas  
    │   │   ├── gometalinter  
    │   │   │   ├── CONTRIBUTING.md  
    │   │   │   ├── COPYING  
    │   │   │   ├── README.md  
    │   │   │   ├── main.go  
    │   │   │   └── regressiontests  

この問題に対処するためにシンプルにshellの設定ファイルを用いて共通のGo言語環境を用意してしまいましょう。

    # 共有のGo言語環境の各ディレクトリを作成  
    $ mkdir -p {$HOME/go,$HOME/go/bin,$HOME/go/pkg,$HOME/go/src}  
    
    # shellの設定ファイルにGOPATHを設定  
    $ vi $HOME/.zshrc  
    
    # 共通のGo言語環境変数とそのパス  
    export $GOPATH="$HOME/go"  
    # ツール群のバイナリファイルが格納される  
    export $GOBIN="$GOPATH/gin"  
    
    # shellの設定を更新  
    $ source $HOME/.zshrc  

次に    :GoInstallBinaries    コマンドでツール群をshellの設定ファイルで指定したパスにインストールします。

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

ツール群の試運転をしてみましょう。animalsパッケージのmonkey.goを修正します。

    $ vi src/animals/monkey.go  
    
    " コマンドモードでGoImports  
    :GoImports  

[![GoImports](https://i.gyazo.com/ff6a7e651889e798892347f9dc460fa2.gif)](https://gyazo.com/ff6a7e651889e798892347f9dc460fa2)

vim-goの    :GoImports    コマンドでfmtパッケージがインポートされていることがわかります。

### vim-goの設定

デフォルトではシンタックスハイライトがオフになっていますので、設定ファイルを用いて設定をカスタマイズしていきましょう。    :GoFmt    コマンドを実行した際の挙動もカスタマイズしてみます。

    $ vi $HOME/.vimrc  
    " vim-go  
    "" mapping  
    """ go runのキーマッピング  
    au FileType go nmap <Leader>gr <Plug>(go-run)  
    """ go testのキーマッピング  
    au FileType go nmap <Leader>gt <Plug>(go-test)  
    "" highlight  
    let g:go_hightlight_functions = 1  
    let g:go_hightlight_methods = 1  
    let g:go_hightlight_structs = 1  
    let g:go_hightlight_interfaces = 1  
    let g:go_hightlight_operators = 1  
    let g:go_hightlight_build_constraints = 1  
    "" GoFmt時にインポートするパッケージを整理(GoFmtはファイル書き込み時に自動的に実行される)  
    let g:go_fmt_command = "goimports"  

試運転をしてみましょう。先ほどと同様にanimalsパッケージのmonkey.goを修正します。

    $ vi src/animals/monkey.go  
    
    " バッファの内容をファイルに書込  
    :w  

[![write](https://i.gyazo.com/6e10844988f2f296f4f350813022fec6.gif)](https://gyazo.com/6e10844988f2f296f4f350813022fec6)

ファイル書き込み時に    :GoFmt    によるフォーマットと、    :GoImports    コマンドによるパッケージのインポートが実行されていることがわかります。

テストも実行してみましょう。


    $ vi src/animals/elephant.go  
    
    " テストの実行  
    <Leader>gt  

[![go test](https://i.gyazo.com/f75ff620de8ff02e59fddc031d5b0afe.gif)](https://gyazo.com/f75ff620de8ff02e59fddc031d5b0afe)

### 自動補完

neocompleteをインストール、設定していれば自動補完もかっこ良くなります。

[![https://gyazo.com/659531f90a07c84c1c42c28492de63f9](https://i.gyazo.com/659531f90a07c84c1c42c28492de63f9.gif)](https://gyazo.com/659531f90a07c84c1c42c28492de63f9)

## 課題

本稿では対応出来ていないvim-goのカスタマイズについて簡単に触れておきます。

### [Tagbar](http://majutsushi.github.io/tagbar/)

- ctagsを適切に設定することで関数の一覧バッファをサイドビューとして表示し、タグジャンプにも利用出来ます。

### [Neovim](https://neovim.io/)対応

- [vim-go開発者のfaithさんもNeovimを利用しだした](https://github.com/fatih/vim-go#using-with-neovim-beta)とのことですが、まだうまく動いておらず、ベータとして提供している様です。

## おわりに

Vimmer向けにGo言語開発環境の実践的な内容をまとめてみました。ここでは触れていませんが、ビルドコマンドも提供されており、vim-goさえあればVimmerのGo言語開発は何とかなるのではと見込んでおります。これは、シンプルな規約に沿ってデザインされたGo言語の恩恵のひとつなのかもしれません。  
弊社では、バックエンド(サーバサイド)エンジニアを[鋭意募集](https://rec-log.jp/site/jobVw.aspx?Oy24gBu7IDWaPGodjJQfOM2igOulIRWoaUoqCXQtO02wg2uzI5WCa8oECbQHOe2KgguNIjWPamoSCpQVOr2Yguu1IxW3aAo6CDQ9OF2cgIueeL4hGOtkKQ1n)しています。Go言語の他にもモダンな技術の導入には前向きに取り組んでおりますので、腕自慢の方はぜひご応募ください。
