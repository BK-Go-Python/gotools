# base64：Base64 编解码

`gotools/base64` 提供 Python `base64` 模块的高频函数。

```go
encoded := base64.B64encode([]byte("hello")) // aGVsbG8=
data, err := base64.B64decode(encoded)

urlsafe := base64.Urlsafe_b64encode([]byte{0xfb, 0xff}) // -_8=
data, err = base64.Urlsafe_b64decode(urlsafe)
```

| Python | gotools |
|---|---|
| `base64.b64encode` | `B64encode` |
| `base64.b64decode` | `B64decode` |
| `base64.urlsafe_b64encode` | `Urlsafe_b64encode` |
| `base64.urlsafe_b64decode` | `Urlsafe_b64decode` |

## DIFF

- 标准编码要求正确的 `=` padding；Python 的 `b64decode` 对缺失 padding 也可能报错，行为一致。
- 未提供 `b32encode`、`b85encode`、`b16encode`、`encodebytes`。
- 解码失败返回 error；Python 抛 `binascii.Error`。
