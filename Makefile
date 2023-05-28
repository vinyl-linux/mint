default: mint

mint: parser/*.go go.mod go.sum
	go build -o $@ ./cmd/
