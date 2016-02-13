# bouSaga2.0
gobou plugin

佐賀配信2.0を真似たレスに反応するgobouプラグイン。
スレッドに書かれた特定の単語などに反応して音声を発する。

# インストール方法
[gobou](https://github.com/DoG-peer/gobou)からgobouをインストール。（詳細はリンクから）
gobouコマンドが動くようになったら、
`gobou install DoG-peer/bouSaga2.0`

`go get`コマンドでは正しくインストールできないので注意。

# 設定
`gobou config bouSaga2.0`
で設定ファイルを編集する。

```config.json
{
  "Cache": "/home/user/.cache/gobou/bouSaga2.0",
  "Res": 500,
  "Url": "http://jbbs.shitaraba.net/bbs/rawmode.cgi/game/57358/1389905050/",
  "Voice": [
    {"Name": "sata", "Words":"ぷえー", "Condition":"star"},
    {"Name": "sata2", "Words":"ぷえーーーー", "Condition":"choo"},
    {"Name": "tada", "Words":"ただの香車を", "Condition":"タダ"},
    {"Name": "gobou", "Words":"ごぼう、かわいい", "Condition":"gobou"},
    {"Name": "saga", "Words":"佐賀さん、はよ配信","Condition":"佐賀"}
  ]
}
```

CacheはgobouのCacheの設定するところ。
Urlに掲示板のURL。
Resに開始レス番号。
Voiceに音声の設定
VoiceのNameに識別用の名前（重複してはいけない）
VoiceのWordsにセリフ
VoiceのConditionに反応する基準になる正規表現
