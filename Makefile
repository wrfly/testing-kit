REPO = github.com/wrfly/testing-kit

.PHONY: build test clean all

BINS = $(shell find * -name 'main.go' | grep -v util)
build:
	@for bin in $(BINS);do \
		NAME=$$(echo $$bin | cut -d"/" -f2); \
		DIR=$$(echo $$bin | cut -d"/" -f1); \
		echo Building $$NAME...;\
		go build -o bin/$$NAME $(REPO)/$$DIR/$$NAME; \
	done

TESTCMD = go test -v -timeout 30s
test:
	@find util/* -maxdepth 0 -type d -exec echo $(REPO)/{} \;\
	| xargs $(TESTCMD)

clean:
	@echo "Clean..."
	@rm -rf bin

all: build test clean