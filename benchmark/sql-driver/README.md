# GO ORM Benchmarks

[![Build Status](https://travis-ci.com/derkan/gobenchorm.svg?branch=master)](https://travis-ci.com/derkan/gobenchorm) [![GoDoc](https://godoc.org/github.com/derkan/gobenchorm/benchs?status.svg)](https://godoc.org/github.com/derkan/gobenchorm/benchs) [![Coverage Status](https://coveralls.io/repos/github/derkan/gobenchorm/badge.svg?branch=master)](https://coveralls.io/github/derkan/gobenchorm?branch=master)

## About

ORM benchmarks for GoLang. Originally forked from [orm-benchmark](https://github.com/milkpod/orm-benchmark).
Contributions are wellcome.

## Environment

- go version go1.9 linux/amd64

## PostgreSQL

- PostgreSQL 12.4 for Linux on WSL2

## ORMs

- [dbr](https://github.com/gocraft/dbr/v2)
- [genmai](https://github.com/naoina/genmai)
- [gorm](https://github.com/jinzhu/gorm)
- [gorp](http://gopkg.in/gorp.v3)
- [pg](https://github.com/go-pg/pg/v10)
- [beego](https://github.com/astaxie/beego/tree/master/orm)
- [sqlx](https://github.com/jmoiron/sqlx)
- [xorm](https://github.com/xormplus/xorm)
- [godb](https://github.com/samonzeweb/godb)
- [upper](https://github.com/upper/db/v4)
- [hood](https://github.com/eaigner/hood)
- [modl](https://github.com/jmoiron/modl)
- [qbs](https://github.com/coocood/qbs)
- [pop](https://github.com/gobuffalo/pop)
- [rel](https://github.com/go-rel/rel)

### Notes

#### Hood

- `hood` needs patch for reflecting `string` values:

  `github.com/eaigner/hood/base.go:50` should be patched to:

  `fieldValue.SetString(string(driverValue.Elem().String()))`

- Multi insert is too slow(over 100 seconds), need check/help

#### QBS

- `qbs` needs patch for reflecting `string` values:

  `github.com/coocood/qbs/base.go:69` should be patched to:

  `fieldValue.SetString(string(driverValue.Elem().String()))`

- `BulkInsert` is not working as expected.

#### Gorm

- [No support for multi insert](https://github.com/jinzhu/gorm/issues/255)

#### Genmai

- Fails on reading 10000 rows (err=>sql: expected 4464 arguments, got 70000)

#### Gorp

- `BulkInsert` is not working as expected. It's too slow.

#### Pop

- `BulkInsert` is not working as expected. It's too slow.

### Prepare DB

```sql
CREATE ROLE bench LOGIN PASSWORD 'pass'
   VALID UNTIL 'infinity';
CREATE DATABASE benchdb
  WITH OWNER = bench;
```

### Run

```go
# build:
cd database/gobenchorm/cmd
go build

# run all benchmarks:
./gobenchorm -multi=1 -orm=all

# run given benchmarks:
./gobenchorm -multi=1 -orm=xorm -orm=raw -orm=godb
```

### Reports

```yaml
Reports:

  2000 times - Insert
  pgx:     0.49s       246671 ns/op     285 B/op     10 allocs/op
  raw:     0.51s       256402 ns/op     696 B/op     18 allocs/op
  beego:     0.61s       303611 ns/op    2425 B/op     56 allocs/op
    qbs:     0.66s       327628 ns/op    5679 B/op    123 allocs/op
      pg:     1.79s       894601 ns/op    1560 B/op     11 allocs/op
    gorm:     3.30s      1651142 ns/op    6825 B/op     97 allocs/op
    xorm:     3.36s      1677576 ns/op    3162 B/op     98 allocs/op
    modl:     3.36s      1679712 ns/op    1686 B/op     43 allocs/op
    godb:     3.42s      1709259 ns/op    4729 B/op    115 allocs/op
    hood:     3.43s      1716343 ns/op    7088 B/op    173 allocs/op
      rel:     3.44s      1719134 ns/op    2446 B/op     49 allocs/op
    gorp:     3.46s      1730252 ns/op    1688 B/op     44 allocs/op
    sqlx:     3.55s      1776314 ns/op    2319 B/op     51 allocs/op
      pop:     3.64s      1817719 ns/op   10324 B/op    248 allocs/op
  genmai:     4.18s      2088679 ns/op    4501 B/op    148 allocs/op
    dbr:     4.31s      2152616 ns/op    2983 B/op     74 allocs/op
    upper:     4.99s      2493620 ns/op   27781 B/op   1185 allocs/op

  500 times - BulkInsert 100 row
  beego:     1.58s      3166656 ns/op  196363 B/op   2845 allocs/op
  genmai:     2.24s      4489226 ns/op  204872 B/op   3066 allocs/op
    xorm:     2.25s      4495976 ns/op  319784 B/op   7542 allocs/op
      pg:     2.30s      4605786 ns/op   19113 B/op    214 allocs/op
    godb:     2.33s      4668724 ns/op  289680 B/op   5994 allocs/op
      rel:     2.37s      4744384 ns/op  287076 B/op   4053 allocs/op
    upper:     2.80s      5602112 ns/op  482084 B/op  19820 allocs/op
      modl:     Don't support bulk insert
        pop:     Problematic bulk insert, too slow
        pgx:     0.00s      1.16 ns/op       0 B/op      0 allocs/op
      gorp:     Problematic bulk insert, too slow
      hood:     Problematic bulk insert, too slow
        qbs:     Don't support bulk insert, err driver: bad connection
      gorm:     Don't support bulk insert - https://github.com/jinzhu/gorm/issues/255
        dbr:     Does not support bulk insert
      sqlx:     benchmark not implemeted yet - https://github.com/jmoiron/sqlx/issues/134
        raw:     0.00s      1.19 ns/op       0 B/op      0 allocs/op

  2000 times - Update
  pgx:     0.24s       118040 ns/op     289 B/op     10 allocs/op
  raw:     0.39s       194489 ns/op     712 B/op     19 allocs/op
  beego:     0.58s       288797 ns/op    1801 B/op     47 allocs/op
    pop:     0.74s       371735 ns/op    6794 B/op    197 allocs/op
    dbr:     0.90s       447675 ns/op    2619 B/op     57 allocs/op
    qbs:     2.87s      1437398 ns/op    5899 B/op    149 allocs/op
      pg:     3.12s      1558253 ns/op     992 B/op     13 allocs/op
    gorm:     3.26s      1627977 ns/op    7467 B/op     93 allocs/op
    godb:     3.35s      1674696 ns/op    5377 B/op    154 allocs/op
    modl:     3.37s      1682538 ns/op    1296 B/op     40 allocs/op
      rel:     3.42s      1709182 ns/op    2608 B/op     50 allocs/op
    xorm:     3.43s      1716999 ns/op    3217 B/op    126 allocs/op
    sqlx:     3.44s      1719007 ns/op    1016 B/op     21 allocs/op
    gorp:     3.44s      1720255 ns/op    1344 B/op     39 allocs/op
  genmai:     3.50s      1748804 ns/op    3520 B/op    146 allocs/op
    hood:     3.53s      1765963 ns/op   13482 B/op    324 allocs/op
    upper:     5.09s      2543716 ns/op   33053 B/op   1491 allocs/op

  4000 times - Read
  pgx:     0.46s       113938 ns/op    1022 B/op      8 allocs/op
  raw:     0.49s       122029 ns/op     888 B/op     24 allocs/op
  beego:     0.53s       133134 ns/op    2112 B/op     75 allocs/op
    gorm:     0.61s       152176 ns/op    4612 B/op     93 allocs/op
      pop:     0.63s       156296 ns/op    3668 B/op     72 allocs/op
    sqlx:     0.68s       169483 ns/op    1744 B/op     38 allocs/op
      pg:     0.71s       178166 ns/op    1262 B/op     14 allocs/op
    gorp:     0.75s       186883 ns/op    3952 B/op    188 allocs/op
      rel:     0.77s       192959 ns/op    1616 B/op     44 allocs/op
    modl:     1.31s       326696 ns/op    1776 B/op     41 allocs/op
      dbr:     1.35s       336486 ns/op    2176 B/op     36 allocs/op
    godb:     1.39s       348661 ns/op    4193 B/op    102 allocs/op
  genmai:     1.69s       422246 ns/op    3312 B/op    171 allocs/op
    upper:     1.86s       463833 ns/op    7402 B/op    293 allocs/op
      xorm:     1.86s       465954 ns/op    8796 B/op    252 allocs/op
      hood:     reflect: call of reflect.Value.Bytes on string Value
                qbs:     reflect: call of reflect.Value.Bytes on string Value

  2000 times - MultiRead limit 1000
  raw:     3.99s      1995605 ns/op  272016 B/op  11657 allocs/op
  pgx:     4.26s      2129283 ns/op  441141 B/op   5020 allocs/op
    pg:     4.48s      2241461 ns/op  321129 B/op   5027 allocs/op
  genmai:     5.56s      2777999 ns/op  420669 B/op  12844 allocs/op
    modl:     5.57s      2787467 ns/op  514071 B/op  16676 allocs/op
    sqlx:     5.62s      2810328 ns/op  499425 B/op  13691 allocs/op
    gorp:     6.15s      3076917 ns/op  737103 B/op  15861 allocs/op
      dbr:     6.50s      3247816 ns/op  514960 B/op  16705 allocs/op
      pop:     6.77s      3383975 ns/op  695013 B/op  14756 allocs/op
    upper:     6.89s      3445156 ns/op  648037 B/op  14047 allocs/op
    beego:     7.75s      3873104 ns/op  746747 B/op  32474 allocs/op
      rel:     8.78s      4392132 ns/op 1010637 B/op  24674 allocs/op
      godb:     8.86s      4429233 ns/op  997659 B/op  31738 allocs/op
      gorm:    11.53s      5766083 ns/op  876107 B/op  36740 allocs/op
      xorm:    18.85s      9424722 ns/op 1447649 B/op  55858 allocs/op
      hood:     reflect: call of reflect.Value.Bytes on string Value
                qbs:     reflect: call of reflect.Value.Bytes on string Value
```
