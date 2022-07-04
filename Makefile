.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports -l -local github.org/tiedotguy/scrubeer .

.PHONY: build
build:
	go build ./cmd/scrubeer/

.PHONE: docker
docker:
	docker build -t scrubeer:local .
