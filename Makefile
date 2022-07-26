ALL: bin/csvtojson bin/analyze bin/snoop bin/assemble

bin/%: cmd/%/main.go $(shell find pkg -type f)
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin