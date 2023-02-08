go test -v ./...
slep 2

docker build -t cake .
slep 2

docker-compose up -d
slep 2

make up
slep 2

exit 1