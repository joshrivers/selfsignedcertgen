run:
	go run examples/example_1/main.go

test:
	go test -v .

coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out

godoc:
	open http://localhost:8000/pkg/github.com/joshrivers/selfsignedcertgen/ && Tgodoc -http localhost:8000