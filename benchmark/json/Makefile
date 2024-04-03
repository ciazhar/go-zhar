install:
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest

easy:
	easyjson -all data/data.go

bench-jts:
	go test -bench=. -benchmem data.go data_easyjson.go json_to_struct_test.go

bench-stj:
	go test -bench=. -benchmem data.go data_easyjson.go struct_to_json_test.go