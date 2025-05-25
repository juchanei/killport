BINARY=port-killer
OUTDIR=bin

.PHONY: all build clean test integration-test

all: build

build:
	@mkdir -p $(OUTDIR)
	go build -o $(OUTDIR)/$(BINARY)

integration-test:
	@$(MAKE) build
	@echo "[통합테스트] main_integration_test.go 실행"
	go test -v main_integration_test.go
