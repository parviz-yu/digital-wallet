package main

import (
	"fmt"

	"github.com/parviz-yu/digital-wallet/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init logger

	// init storage

}
