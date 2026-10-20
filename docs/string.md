# string：模板与字符常量

`gotools/string` 提供 Python `string` 模块的高频工具。

```go
string.Capwords("hello   world") // "Hello World"

tmpl := string.NewTemplate("$name is ${age}")
text, err := tmpl.Substitute(map[string]string{"name": "Tom", "age": "18"})
// Tom is 18

safe := tmpl.Safe_substitute(map[string]string{"name": "Tom"})
// Tom is ${age}
```

## API

- `Capwords`：按空白切词，首字母大写、其余小写
- `NewTemplate`、`Template`、`Substitute`、`Safe_substitute`
- 常量：`Ascii_lowercase`、`Ascii_uppercase`、`Ascii_letters`、`Digits`、`Hexdigits`、`Octdigits`、`Punctuation`、`Whitespace`、`Printable`

## Template 规则

- `$name` / `${name}` 替换
- `$$` 输出字面 `$`
- `Substitute` 对缺失变量/非法 `$` 返回 error
- `Safe_substitute` 对缺失变量/非法 `$` 保留原样

## DIFF

- 占位符标识符只支持 ASCII 字母/数字/下划线；Python 支持 Unicode 标识符。
- 没有 `Formatter`、`Template.get_identifiers`、`capwords` 的 `sep` 参数。
- 包名 `string` 会遮蔽 Go 内建 `string` 类型；使用时建议别名导入：
  ```go
  pystring "github.com/BK-Go-Python/gotools/string"
  ```
