# Json Benchmark

> Tujuan dari benchmark ini adalah untuk mencari alternative dari package encoding/json. Benchmark ini hanya akan
> menggunakan 2 usecase yaitu marshal dari struct ke json dan unmarshal dari json ke struct. Library yang tidak mengcover
> salah satu dari use case itu tidak akan di benchmark.

## Tools
- MacbookPro M1 2020
- Go (go1.21.4 darwin/arm64)

## Library List

- [DJSON](https://github.com/a8m/djson) 
- [easyjson](https://github.com/mailru/easyjson)
- [ffjson](https://github.com/pquerna/ffjson)
- [Gabs](https://github.com/Jeffail/gabs)
- [gojay](https://github.com/francoispqt/gojay)
- [go-codec](https://github.com/ugorji/go)
- [go-json](https://github.com/goccy/go-json)
- [go-simplejson](https://github.com/bitly/go-simplejson)
- [go-ujson](https://github.com/mreiferson/go-ujson)
- [Jason](https://github.com/antonholmquist/jason) 
- [jettison](https://github.com/wI2L/jettison)
- [Jingo](https://github.com/bet365/jingo)
- [jsonparser](https://github.com/buger/jsonparser)
- [json-iterator](https://github.com/json-iterator/go)
- [segmentio](https://github.com/segmentio/encoding/tree/master/json)
- [simdjson-go](https://github.com/minio/simdjson-go)
- [Stdlib](https://golang.org/pkg/encoding/json)


## Run Benchmark
### Json To Struct
```bash
    go test -bench=. -benchmem data.go data_easyjson.go json_to_struct_test.go 
```
### Struct To Json
```bash
    go test -bench=. -benchmem data.go data_easyjson.go struct_to_json_test.go 
```

## Result
### Struct to JSON
```text
goos: darwin
goarch: arm64
BenchmarkFFjsonSTJ-8         	 7635714	       157.3 ns/op	     272 B/op	       3 allocs/op
BenchmarkEasyJsonSTJ-8       	 6101962	       173.3 ns/op	     384 B/op	       4 allocs/op
BenchmarkGojayMarshalSTJ-8   	 6163507	       194.7 ns/op	     592 B/op	       2 allocs/op
BenchmarkGojayEncodeSTJ-8    	 4848584	       246.0 ns/op	     720 B/op	       4 allocs/op
BenchmarkJsonIteratorSTJ-8   	 5467234	       216.9 ns/op	     352 B/op	       4 allocs/op
BenchmarkGoJsonSTJ-8         	 3319440	       362.5 ns/op	     352 B/op	       4 allocs/op
BenchmarkStdlibSTJ-8         	 2286153	       529.3 ns/op	     352 B/op	       4 allocs/op
BenchmarkGabsSTJ-8           	 2284868	       525.6 ns/op	     352 B/op	       4 allocs/op
BenchmarkSegmentioSTJ-8      	 2534533	       470.0 ns/op	     416 B/op	       5 allocs/op
BenchmarkGoCodecSTJ-8        	 2125832	       557.0 ns/op	    1888 B/op	      12 allocs/op
BenchmarkJettisonSTJ-8       	 1919134	       626.1 ns/op	     352 B/op	       4 allocs/op
BenchmarkJingoSTJ-8          	  949704	      1268 ns/op	    2432 B/op	      22 allocs/op
```
### JSON To Struct
```text
BenchmarkGojayUnmarshal-8     	 6681082	       179.3 ns/op	     224 B/op	       3 allocs/op
BenchmarkEasyJson-8           	 4610371	       260.3 ns/op	     104 B/op	       4 allocs/op
BenchmarkFFjson-8             	 4058197	       295.1 ns/op	     168 B/op	       5 allocs/op
BenchmarkDJSONAllocString-8   	 3162339	       380.1 ns/op	     560 B/op	       9 allocs/op
BenchmarkGoJson-8             	 2999654	       399.5 ns/op	     328 B/op	       7 allocs/op
BenchmarkJsonParser-8         	 3010540	       397.1 ns/op	     424 B/op	       8 allocs/op
BenchmarkSegmentio-8          	 2988072	       402.3 ns/op	     168 B/op	       5 allocs/op
BenchmarkDJSON-8              	 2551957	       449.3 ns/op	     536 B/op	      16 allocs/op
BenchmarkGoUJson-8            	 1778368	       677.3 ns/op	     656 B/op	      22 allocs/op
BenchmarkGoCodec-8            	 1531711	       776.6 ns/op	    1656 B/op	      15 allocs/op
BenchmarkStdlib-8             	 1550151	       777.3 ns/op	     320 B/op	       7 allocs/op
BenchmarkGojayDecode-8        	 4633572	       258.6 ns/op	     800 B/op	       5 allocs/op
BenchmarkGabs-8               	 1209264	       989.5 ns/op	     784 B/op	      26 allocs/op
BenchmarkGoSimpleJson-8       	 1000000	      1063 ns/op	    1464 B/op	      24 allocs/op
BenchmarkJason-8              	  412992	      2904 ns/op	    3776 B/op	      67 allocs/op
```