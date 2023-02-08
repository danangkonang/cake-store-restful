go clean -testcache
sleep 2

go test -v ./...
sleep 2

docker build -t cake .
sleep 2

docker-compose up -d
echo "waiting for database running"
sleep 5

docker exec app sh -c "./goql up migration --dir schema/migration --db mysql://danang:danang@tcp\(mysql:3306\)/db-cake?parseTime=true&loc=Asia%2FJakarta"

exit 1