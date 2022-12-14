#### Run

```shell
$ go run github.com/adzeitor/csvmagi "This is column foo {foo} and column bar {bar}" < file.csv
```

Also you can use column numbers using special variables `{_1}`, `{_2}`, etc in template: 

```shell
$ go run github.com/adzeitor/csvmagi 'UPDATE some_table SET foo="{_2}" WHERE id="{_1}"'< file.csv
```

In strict mode you receive errors on undefined columns:

```shell
$ go run github.com/adzeitor/csvmagi -strict 'Some of the last columns {_999}' < file.csv
```

#### Limitations:
- only support CSV files with headers
- no proper case-insensitive key matching
