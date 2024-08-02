package main

import (
	"context"
	"fmt"
	client "github.com/hina1314/hina/client/api"
)

func main() {
	addr := ":8080"
	ctx := context.Background()
	hina, err := client.NewHinaClient(ctx, addr)
	if err != nil {
		return
	}
	defer hina.Close()

	if err := hina.Set("test", "123"); err != nil {
		fmt.Println(err)
	}

}
