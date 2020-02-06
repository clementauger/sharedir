build:
	go-bindata assets/
	go build -tags prod -o sharedir
clean:
	rm bindata.go
	rm sharedir
run:
	go run .
prod:
	go run -tags prod .
