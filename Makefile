# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2023 Dell Inc, or its subsidiaries.

ROOT_DIR='.'
PROJECTNAME=$(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

compile: get inventory

base:
	@echo "  >  Building base binaries..."
	@CGO_ENABLED=0 go build -o ${PROJECTNAME} ./cmd/...

get:
	@echo "  >  Checking if there are any missing dependencies..."
	@CGO_ENABLED=0 go get ./...

inventory:
	@echo "  >  building binaries for inventory ..."
	@CGO_ENABLED=0 go build -o ${PROJECTNAME} -tags=inventory ./cmd/...

all:
	@echo "  >  building full capabilities binaries..."