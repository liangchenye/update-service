PREFIX ?= $(DESTDIR)/usr
BINDIR ?= $(DESTDIR)/usr/bin

all:
	go build -tags "$(BUILDTAGS)" -o upserver .

.PHONY: man
man:
	go-md2man -in "man/upserver.1.md" -out "upserver.1"

install: man
	install -d -m 755 $(BINDIR)
	install -m 755 upserver $(BINDIR)
	install -d -m 755 $(PREFIX)/share/man/man1
	install -m 644 *.1 $(PREFIX)/share/man/man1

uninstall:
	rm -f $(PREFIX)/share/man/man1/upserver.1 
clean:
	rm -f upserver upserver.1

.PHONY: test .gofmt .govet .golint

test: .gofmt .govet .golint

.gofmt:
	go fmt ./...

.govet:
	go vet -x ./...

.golint:
	golint ./...

