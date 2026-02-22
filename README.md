# peercast-pcp

Go標準ライブラリのみで実装された、[PeerCast Protocol (PCP)](docs/PCP_SPEC.md) パケットのシリアライズ/デシリアライズライブラリです。

## 概要

PCP は PeerCast が P2P ストリーミングに使用する Tag-Length-Value (TLV) バイナリプロトコルで、TCP 上で動作します。
マルチバイト整数はすべて **リトルエンディアン** でエンコードされます。
プロトコルの基本単位は **Atom** であり、次の 2 種類があります。

- **コンテナ（親）Atom** — 0 個以上の子 Atom を保持する
- **データ（葉）Atom** — バイト列のペイロードを保持する

## 動作要件

- Go 1.22 以降
- サードパーティ依存なし

## インストール

```bash
go get github.com/titagaki/peercast-pcp
```

```go
import "github.com/titagaki/peercast-pcp/pcp"
```

## 使い方

### 接続してマジックアトムを送信する

```go
conn, err := pcp.Dial("127.0.0.1:7144")
if err != nil {
    // エラー処理
}
defer conn.Close()
```

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

conn, err := pcp.DialContext(ctx, "127.0.0.1:7144")
if err != nil {
    // エラー処理
}
defer conn.Close()
```

### ストリームから Atom を読み込む

```go
import "github.com/titagaki/peercast-pcp/pcp"

atom, err := pcp.ReadAtom(conn)
if err != nil {
    // エラー処理
}

switch atom.Tag {
case pcp.PCPHelo:
    // ハンドシェイク処理
case pcp.PCPChan:
    // チャンネル処理
}
```

### Atom を生成してストリームに書き込む

```go
// コンテナ（親）Atom
helo := pcp.NewParentAtom(pcp.PCPHelo,
    pcp.NewStringAtom(pcp.PCPHeloAgent, "MyClient/1.0"),
    pcp.NewShortAtom(pcp.PCPHeloPort, 7144),
    pcp.NewIntAtom(pcp.PCPHeloVersion, 1218),
)

if err := helo.Write(conn); err != nil {
    // エラー処理
}
```

### データ Atom から型付きの値を取得する

```go
// uint32 (INT)
v, err := atom.GetInt()

// uint16 (SHORT)
v, err := atom.GetShort()

// 1 バイト
b, err := atom.GetByte()

// null 終端文字列
s := atom.GetString()

// 16 バイト GnuID
id, err := atom.GetID()
```

### 子 Atom を検索する

```go
// タグが一致する最初の子 Atom
child := parent.FindChild(pcp.PCPHeloAgent)

// タグが一致するすべての子 Atom
children := parent.FindChildren(pcp.PCPBcstFrom)
```

### 未知の Atom をスキップする（前方互換性）

```go
if err := pcp.SkipAtom(r); err != nil {
    // エラー処理
}
```

## パッケージ構成

| ファイル | 内容 |
|---------|------|
| `atom.go` | `Atom` 型・コンストラクタ・`ReadAtom`・`Write`・`SkipAtom` |
| `conn.go` | `Conn` 型・`Dial`・`DialContext`・`NewConn`・`ReadAtom`/`WriteAtom` |
| `id4.go` | `ID4` 型・`GnuID` 型 |
| `tags.go` | 既知タグ変数一覧（`PCPHelo`、`PCPChan` など） |
| `constants.go` | 数値定数（フラグ、エラーコードなど） |
| `packets.go` | 型付きパケット構造体（`HeloPacket`、`ChanPacket` など） |

## プロトコルのワイヤーフォーマット

```
┌─────────────────────────────────────────┐
│  Tag       [4]byte                      │
│  Length    uint32 (リトルエンディアン)   │
│    MSB=1 → 親; 下位 31 ビット = 子の数  │
│    MSB=0 → データ; 値 = ペイロードサイズ │
│  Payload   N 個の子 Atom または N バイト │
└─────────────────────────────────────────┘
```

※ `N` は、親 Atom の場合は「子 Atom の数」、データ Atom の場合は「ペイロードのバイト数」を表します。

## 定数

### ブロードキャストグループフラグ (`grp`)

| 定数 | 値 | 説明 |
|------|----|------|
| `PCPBcstGroupAll` | `0xff` | 全グループ |
| `PCPBcstGroupRoot` | `0x01` | ルートサーバー |
| `PCPBcstGroupTrackers` | `0x02` | トラッカーノード |
| `PCPBcstGroupRelays` | `0x04` | リレーノード |

### ホストフラグ (`flg1`)

| 定数 | 値 | 説明 |
|------|----|------|
| `PCPHostFlags1Tracker` | `0x01` | トラッカーノード |
| `PCPHostFlags1Relay` | `0x02` | リレー中 |
| `PCPHostFlags1Direct` | `0x04` | ダイレクト接続可 |
| `PCPHostFlags1Push` | `0x08` | プッシュ接続が必要 |
| `PCPHostFlags1Recv` | `0x10` | 受信中 |
| `PCPHostFlags1CIN` | `0x20` | コントロールイン |
| `PCPHostFlags1Private` | `0x40` | プライベートネットワーク |

## 開発

```bash
# テストをすべて実行
go test ./pcp/...

# 詳細出力でテストを実行
go test -v ./pcp/...
```

## 参考資料

- [PCP 仕様書](docs/PCP_SPEC.md)
- [pcp.h リファレンス](docs/reference/pcp.h)
- [atom.h リファレンス](docs/reference/atom.h)

## ライセンス

This project is licensed under the GNU General Public License v3.0.
Portions of this software are Copyright (C) 2026 ITAGAKI Takayuki
