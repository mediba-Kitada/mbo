# VimでTwitter

## TwitVim

[TwitVim - Twitter client for Vim](https://github.com/twitvim/twitvim)

日がな一日黒い画面と対峙するVimmerには、必携のVim pluginです。
お好きなプラグインマネージャーでインストールしましょう。

## Tips

部署異動に伴って、生活環境を見直していたのですがTwitVimを更新するとユーザ認証でコケるというトラブルに直面しました。

- [Encoding '&' on Zsh · Issue #7 · twitvim/twitvim · GitHub](https://github.com/twitvim/twitvim/issues/7#issuecomment-289390267)

issueを投げたところ、あの[mattnさん](https://github.com/mattn)にサポート頂いたので浮かれて投稿してます。
issueに記載してますが、以下の対応でVimでTwitter出来るようになりました。

- URIをshell用に置換する際の関数がよろしくなかった
	- mattnさんに最新版を ```git push``` して頂きました。
- Python 3だとタイムアウトしたので、Python 2.7を利用するように設定した

参考までにvimrcの一部を貼っておきます。

```vim
" TwitVim
NeoBundle 'twitvim/twitvim'
let twitvim_enable_python = 1
```
