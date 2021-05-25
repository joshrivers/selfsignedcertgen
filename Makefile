run:
	go run examples/example_1/main.go

test:
	go test -v .

coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out

readme:
	goreadme -credit=false -badge-codecov > README.md