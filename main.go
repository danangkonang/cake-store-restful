package main

import (
	"github.com/danangkonang/cake-store-restful/app"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	app.Run()
}
