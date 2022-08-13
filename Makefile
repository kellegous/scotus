
SHA := $(shell git rev-parse HEAD)

ASSETS_DIR := pkg/web/dist

ASSETS := \
	$(ASSETS_DIR)/a/index.js \
	$(ASSETS_DIR)/a/index.html \
	$(ASSETS_DIR)/b/index.js \
	$(ASSETS_DIR)/b/index.html

ALL: bin/csvtojson bin/assemble bin/ot21 bin/server

.PRECIOUS: $(ASSETS)

bin/%: cmd/%/main.go $(ASSETS) $(shell find pkg -type f)
	go build -o $@ ./cmd/$*

serve-dev: bin/server
	bin/server --http.assets-dir=pkg/web/dist

$(ASSETS_DIR)/%.js: node_modules/.build $(shell find src -type f \( -name '*.ts' -or -name '*.scss' \))
	npx webpack build --mode=production

$(ASSETS_DIR)/%.html: src/%.html bin/render_html
	bin/render_html -v sha=$(SHA) $< $@

node_modules/.build:
	npm install
	date > $@

bin/render_html:
	go build -o $@ github.com/kellegous/render_html

nuke: clean
	rm -rf node_modules

clean:
	rm -rf bin $(ASSETS)