package main

import (
	"fmt"
	"os"
)

func main() {
	hostIP := os.Getenv("HOST_IP")
	fmt.Println("HOST_IP:", hostIP)
}