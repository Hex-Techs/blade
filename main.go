package main

import (
	"fmt"

	"github.com/fize/go-ext/log"
	"github.com/hex-techs/blade/cmd"
	"github.com/hex-techs/blade/pkg/utils/config"
)

func main() {
	r := cmd.Run()
	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", config.Read().Service.ServerPort)); err != nil {
		log.Fatalf("run error: %v", err)
	}
}
