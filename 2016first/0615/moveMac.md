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

- 既存Macで設定をexport

```bash
$ /Applications/Karabiner.app/Contents/Library/bin/karabiner export > ~/dotfiles/karabiner.sh
```

- 移行先Macに設定をimport

```
$ ~/dotfiles/karabiner.sh
```

### google日本語入力

https://tools.google.com/dlpage/japaneseinput/eula.html?platform=mac&hl=ja

### shifit

http://shift-it.en.softonic.com/mac

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

