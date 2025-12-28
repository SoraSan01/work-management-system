package main

import (
	"workms/internal/server"
)

func main() {
	r := server.NewServer() // create server
	r.Run(":8080")          // start server on port 8080
}
