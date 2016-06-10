# 開発用Macお引越し

## 目標

- 再現可能なお引越し手順を構築する
- 再現可能なローカル環境を構築する

## セットアップ

- まずは、ヘルプデスク苅谷さんにADの設定とかしてもらう
    - ADログイン確認
    - プロファイルのインストール

### AppleID設定

- mediba.jpのアカウントを作成
- クレジットカードの情報は、アカウント作成後に削除

### システムアップデート

- App Storeの各アプリ、OSをアップデート

### Xcode

- App StoreからDL、インストール
- Command Line Tool for Xcodeをインストール

```bash
$ xcode-select --install
$ sudo xcrun cc
agree
```

### システム環境設定

#### トラックパッド

- ポイントとクリック
    - タップでクリック(1本指でタップ)
    - 軌跡の速さ(目盛りを真ん中に)
- その他のジェスチャ
    - Mission Controll
        - 利用しない
    - アプリケーションExpose
        - 利用しない
    - Launchpad
        - 利用しない

#### Dock

- サイズ
    - 小に目盛りを寄せる
- 画面上の位置
    - 左
- Dockを自動的に隠す/表示
    - チェックする
- 起動中のアプリケーションをアニメーションで表示
    - チェックを外す
- 起動済みのアプリケーションにインジケータを表示
    - チェックを外す

#### キーボード

- キーボード
    - F1,F2などのすべてのキーを標準のファンクションキーとして使用
- ショートカット
    - Spotlight
        - Spotlight検索を表示
            Option + Space
        - Finderの検索ウィンドウを表示
            Controll + Option + Space
    - 入力ソース
        - 前の入力ソースを選択
            - チェックを外す
        - 入力メニューの次のソースを選択
            - チェックを外す
### terminal

- 一般
    - Pro
- プロファイル
    - Proを選択
    - テキスト
        - テキスト
            - テキストをアンチエイリアス処理
            - 点滅テキストの使用を許可
            - ANSIカラーを表示
        - カーソル
            - ブロック
                - チェックを入れる
            - 点滅カーソル
                - チェックを入れる
    - ウィンドウ
        - ウィンドウサイズ
            - 列
                - 210
            - 行
                - 60
    - シェル
        - コマンドを実行
            - チェックを入れる
            - tmux を指定
        - シェルの終了時
            - シェルが正常に終了した場合は閉じる

### Google Chrome

- [Share Extensions](https://chrome.google.com/webstore/detail/share-extensions/chdafcbnfkfenoeejpaeenpdamhmalhe)でExtensions一覧をエクスポート
- safariを起動して手動でインストール

### SSH

```
$ ssh-keygen -t rsa -C kitada@mediba.jp
```

- 公開鍵をgithubとかbitbucketに登録する

### Karabiner

- Chromeを起動し、アプリケーションをDL
    https://pqrs.org/osx/karabiner/index.html.ja
- システム環境設定
    - セキュリティとプライバシー
        - プライバシー
            - Karabinerのパスを指定し、制御を許可する
            - Karabiner_AZNotfilterを指定し、制御を許可する

- 既存Macで設定をexport

```bash
$ /Applications/Karabiner.app/Contents/Library/bin/karabiner export > ~/dotfiles/karabiner.sh
```

- KarabinerをGUIで開く
- 移行先Macに設定をimport

```
$ ~/dotfiles/karabiner.sh
```

### google日本語入力

- https://tools.google.com/dlpage/japaneseinput/eula.html?platform=mac&hl=ja
- システム環境設定
    - キーボード
        - 入力ソース
            - ひらがな(Google)
            - 英数(Google)

### shifit

- http://shift-it.en.softonic.com/mac
- システム環境設定
    - セキュリティとプライバシー
        - プライバシー
            - shiftitのパスを指定し、制御を許可する

### homebrew

```bash
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
$ brew tap b4b4r07/brionac
$ brew install brionac
```

### CLI

- homebrew
    - brionac
- dotfiles

### GUI

- GUIなアプリは、手動でインストール
- shift
- google日本語入力
- Chrome
- Karabina

