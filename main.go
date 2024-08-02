package main

import (
	"context"
	"flag"
	"github.com/hina1314/hina/server"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	ctx := context.Background()
	svr := server.NewServer(ctx, *addr)
	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
