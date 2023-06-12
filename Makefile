default: mint

mint: parser/*.go go.mod go.sum
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@ ./cmd/
