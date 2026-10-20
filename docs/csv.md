# csv：CSV 读写

`gotools/csv` 包装 `encoding/csv`，默认使用 CRLF 行尾，与 Python `csv` 默认一致。

```go
records, err := csv.Read(strings.NewReader("a,b\n1,2\n"))
err = csv.WriteFile("data.csv", [][]string{{"a", "b"}, {"1", "2"}})

file, _ := os.Open("data.csv")
rows, err := csv.ReadDicts(file)
fmt.Println(rows[0]["a"]) // 1
```

| Python | gotools |
|---|---|
| `csv.reader(f)` | `csv.NewReader(r)` / `csv.Read(r)` |
| `csv.writer(f)` | `csv.NewWriter(w)` / `csv.Write(w, records)` |
| `csv.DictReader(f)` | `csv.ReadDicts(r)` / `csv.ReadDictsFile(path)` |
| 文件读写 | `csv.ReadFile` / `csv.WriteFile` |

## DIFF

- `DictReader` 的额外字段被忽略，Python 默认放到 `None` 键下。
- 没有 `DictWriter`、`Sniffer`、`Dialect` 注册表。
- `Reader` 的 `FieldsPerRecord` 默认与 `encoding/csv` 一致：以第一行为准。
- `Writer` 默认 CRLF；需要 LF 时调用 `SetUseCRLF(false)`。
- 不支持 Python 的 `lineterminator`、`quoting` 等完整方言配置。
