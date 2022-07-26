ALL: bin/csvtojson bin/assemble

bin/%: cmd/%/main.go $(shell find pkg -type f)
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin