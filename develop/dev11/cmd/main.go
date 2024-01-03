package main

import (
	"fmt"
	"main/develop/dev11/internal/config"
)

func main() {
	cfg := config.ReadConfig()
	fmt.Println(cfg)
}
