package main

import (
	"go-study/day01/internal/web"
)

func main() {
	server := web.RegisterRoutes()
	server.Run(":8080")
}
