default: test vet lint

test:
	go test .

vet:
	go vet ./...

lint:
	golint ./...

doc:
	godoc -http=:6060

cov:
	 go test -coverprofile=cov.out
	 go tool cover -html=cov.out
