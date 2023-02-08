# Running App

## 1. Download And Deploy
```bash
git clone https://github.com/danangkonang/cake-store-restful.git

cd cake-store-restful

# linux
./build.sh

# windows/macos/linux
go clean -testcache

go test -v ./...

docker build -t cake .

docker-compose up -d

docker exec app sh -c "./goql up migration --dir schema/migration --db mysql://danang:danang@tcp\(mysql:3306\)/db-cake?parseTime=true&loc=Asia%2FJakarta"
```

## 3. Testing

### - Store cake
```http
POST /api/v1/cake HTTP/1.1
Host: localhost:9000
Content-Type: application/json

{
  "title": "Lemon cheesecake",
  "description": "A cheesecake made of lemon",
  "rating": 7,
  "image": "https://img.taste.com.au/ynYrqkOs/w720-h480-cfill-q80/taste/2016/11/sunny-lemon-cheesecake-102220-1.jpeg"
}
```

### - Update cake
```http
PUT /api/v1/cake HTTP/1.1
Host: localhost:9000
Content-Type: application/json

{
  "id: 1,
  "title": "edit Lemon cheesecake",
  "description": "A cheesecake made of lemon",
  "rating": 7,
  "image": "https://img.taste.com.au/ynYrqkOs/w720-h480-cfill-q80/taste/2016/11/sunny-lemon-cheesecake-102220-1.jpeg"
}
```

### - Find cake
```http
GET /api/v1/cakes HTTP/1.1
Host: localhost:9000
Content-Type: application/json
```

### - Detail cake
```http
GET /api/v1/cake/1 HTTP/1.1
Host: localhost:9000
Content-Type: application/json
```

### - Destroy cake
```http
DELETE /api/v1/cake HTTP/1.1
Host: localhost:9000
Content-Type: application/json

{
  "id: 1
}
```