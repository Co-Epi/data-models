.PHONY: coepi

GOBIN = $(shell pwd)/bin
GO ?= latest

coepi:
		go build -o bin/coepi
		@echo "Done building coepi.  Run \"$(GOBIN)/coepi\" to launch coepi."
