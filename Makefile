BINARY=port-killer
OUTDIR=bin

.PHONY: all build clean test

all: build

build:
	@mkdir -p $(OUTDIR)
	go build -o $(OUTDIR)/$(BINARY)
