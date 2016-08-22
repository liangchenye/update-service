PREFIX ?= $(DESTDIR)/usr
BINDIR ?= $(DESTDIR)/usr/bin

all:
	go build -tags "$(BUILDTAGS)" -o upclient .

.PHONY: man
man:
	go-md2man -in "man/upclient.1.md" -out "upclient.1"

install: man
	install -d -m 755 $(BINDIR)
	install -m 755 upclient $(BINDIR)
	install -d -m 755 $(PREFIX)/share/man/man1
	install -m 644 *.1 $(PREFIX)/share/man/man1
	install -d -m 755 $(PREFIX)/share/bash-completion/completions
	install -m 644 completions/bash/upclient $(PREFIX)/share/bash-completion/completions

uninstall:
	rm -f $(PREFIX)/share/man/man1/upclient.1 $(PREFIX)/share/bash-completion/completions/upclient
clean:
	rm -f upclient upclient.1

.PHONY: test .gofmt .govet .golint

test: .gofmt .govet .golint

.gofmt:
	go fmt ./...

.govet:
	go vet -x ./...

.golint:
	golint ./...

