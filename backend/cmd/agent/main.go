package main

import (
	"fmt"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
)

func main() {
	cfg := config.Get()
	fmt.Println(cfg.Env)
}
