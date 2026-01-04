package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/voidarchive/go-nepse"
)

func main() {
	opts := nepse.DefaultOptions()
	opts.TLSVerification = false

	client, err := nepse.NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := client.Companies(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Companies count: %d\n", len(data))
	if len(data) > 0 {
		fmt.Printf("First company: %+v\n", data[0])
	}
}
