LDFLAGS := -ldflags '-X main.url=$(URL)'

build: env
	go build -race $(LDFLAGS) -v .

help:
	@awk '/^[^ ]*:/ { gsub(":.*", ""); print }' Makefile

env:
	@test $(URL) || { echo 'err: $$URL is unset'; exit 1; }

rpi: env
	@test $(URL) || { echo 'err: $$URL is unset'; exit 1; }
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -v .

oldrpi: env
	@test $(URL) || { echo 'err: $$URL is unset'; exit 1; }
	GOOS=linux GOARCH=arm GOARM=6 go build $(LDFLAGS) -v .

macos: env
	@test $(URL) || { echo 'err: $$URL is unset'; exit 1; }
	GOOS=darwin GOARCH=amd64 go build -race $(LDFLAGS) -v .

linux: env
	GOOS=linux GOARCH=amd64 go build -race $(LDFLAGS) -v .

test:
	go test -race -v .

bench:
	go test -race -bench=. -benchmem

clean:
	go clean

.PHONY: bench build clean env help linux macos oldrpi rpi test
