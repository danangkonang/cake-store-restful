up:
	./goql up migration --dir schema/migration --db mysql://danang:danang@tcp\(localhost:3306\)/db-cake?parseTime=true&loc=Asia%2FJakarta

down:
	./goql down migration --dir schema/migration --db mysql://danang:danang@tcp\(localhost:3306\)/db-cake?parseTime=true&loc=Asia%2FJakarta
