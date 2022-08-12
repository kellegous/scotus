ALL: bin/csvtojson bin/assemble bin/ot21 bin/server

bin/%: cmd/%/main.go $(shell find pkg -type f)
	go build -o $@ ./cmd/$*

serve-dev: bin/server
	bin/server --http.assets-dir=pkg/web/dist

clean:
	rm -rf bin