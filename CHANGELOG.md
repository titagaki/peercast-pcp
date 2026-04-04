# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## v0.3.0

### Breaking Changes

以下の関数のシグネチャが変更されました。戻り値に `error` が追加されています。

| 関数 | v0.2.0 | 現在 |
|------|--------|------|
| `IPv4ToUint32` | `func(net.IP) uint32` | `func(net.IP) (uint32, error)` |
| `ParseChanInfo` | `func(*Atom) ChanInfo` | `func(*Atom) (ChanInfo, error)` |
| `ParseChanTrack` | `func(*Atom) ChanTrack` | `func(*Atom) (ChanTrack, error)` |
| `ParseHostPacket` | `func(*Atom) HostPacket` | `func(*Atom) (HostPacket, error)` |

**移行方法:** 戻り値を 2 値で受け取り、`error` をハンドリングしてください。

```go
// v0.2.0
info := pcp.ParseChanInfo(atom)
v := pcp.IPv4ToUint32(ip)

// v0.3.0 以降
info, err := pcp.ParseChanInfo(atom)
v, err := pcp.IPv4ToUint32(ip)
```

### New Features

- **全パケット型の Build/Parse 関数を追加**: `HeloPacket`, `RootPacket`, `ChanPktData`, `ChanPacket`, `BcstPacket`, `PushPacket`, `GetPacket`, `MesgPacket`
  - v0.2.0 では `ChanInfo`, `ChanTrack`, `HostPacket` のみ対応していました
- **`MaxAtomDataSize` 定数 (16 MiB)**: `ReadAtom`/`SkipAtom` がデータ Atom のペイロードサイズを検証し、過大なメモリ割り当てを防止します
- **Atom ネスト深さ制限 (64)**: `ReadAtom`/`SkipAtom` が再帰深さを制限し、悪意あるストリームによるスタックオーバーフローを防止します
- **`NewConn` のエラーハンドリング改善**: マジック Atom の書き込み失敗時に `errors.Join` で書き込みエラーと `Close` エラーの両方を返します

## v0.2.0

- `ChanInfo`, `ChanTrack`, `HostPacket` の Build/Parse 関数を追加 (`packets_codec.go`)
- IPv4 アドレス変換ユーティリティを追加 (`ip.go`): `IPv4ToUint32`, `IPv4FromUint32`, `DecodeIPv4`
- ライセンスを GPLv3 に変更

## v0.1.0

- 初回リリース
- `Atom` 型・コンストラクタ・`ReadAtom`・`Write`・`SkipAtom`
- `Conn` 型・`Dial`・`DialContext`・`NewConn`
- `ID4` 型・`GnuID` 型
- 既知タグ変数・数値定数
- 型付きパケット構造体定義 (`packets.go`)
