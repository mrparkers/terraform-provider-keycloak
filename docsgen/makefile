MAKEFLAGS += --silent

.PHONY=build docs clean-docs

build:
	 GO111MODULE=on go build -o docsgen

clean-docs:
	rm -rf docs

docs: clean-docs build
	./docsgen
