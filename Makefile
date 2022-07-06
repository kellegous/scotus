ALL: bin/csvtojson bin/analyze

bin/%: cmd/%/main.go $(shell find pkg -type f)
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin