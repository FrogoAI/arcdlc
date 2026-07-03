# arctool — deterministic companion CLI for the ArcDLC plan.
#
#   make build     # local binary at bin/arctool
#   make install   # copy binary into BINDIR (default ~/.local/bin, usually on PATH)
#   make test      # go test ./...
#   make release   # static cross-compiled binaries in dist/ for all platforms
#
# The binary is pure standard library, so release builds are static
# (CGO_ENABLED=0) and need no runtime on the target host.

BIN       ?= arctool
CMD       := ./cmd/arctool
BINDIR    ?= $(HOME)/.local/bin
DIST      ?= dist
LDFLAGS   := -s -w
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: build install test release clean

build:
	go build -trimpath -ldflags='$(LDFLAGS)' -o bin/$(BIN) $(CMD)

install: build
	@mkdir -p $(BINDIR)
	install -m 0755 bin/$(BIN) $(BINDIR)/$(BIN)
	@echo "installed $(BINDIR)/$(BIN)  (ensure $(BINDIR) is on PATH)"

test:
	go test ./...

release:
	@mkdir -p $(DIST)
	@for p in $(PLATFORMS); do \
	  os=$${p%/*}; arch=$${p#*/}; \
	  echo "  $(DIST)/$(BIN)-$$os-$$arch"; \
	  CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
	    go build -trimpath -ldflags='$(LDFLAGS)' -o $(DIST)/$(BIN)-$$os-$$arch $(CMD) || exit 1; \
	done

clean:
	rm -rf bin $(DIST)
