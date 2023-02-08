FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY .env .
COPY goql .
COPY schema schema
COPY --from=builder /app/main /app/main
EXPOSE 9000
CMD ["/app/main"]