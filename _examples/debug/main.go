package main

import (
	"bytes"
	"context"
	"encoding/json"
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

	// Endpoints to test with security ID 2781 (example company)
	endpoints := []struct {
		name     string
		endpoint string
	}{
		{"Profile", "/api/nots/security/profile/2781"},
		{"BoardOfDirectors", "/api/nots/security/boardOfDirectors/2781"},
		{"CorporateActions", "/api/nots/security/corporate-actions/2781"},
		{"Reports", "/api/nots/application/reports/2781"},
		{"Dividend", "/api/nots/application/dividend/2781"},
	}

	for _, ep := range endpoints {
		fmt.Printf("\n=== %s ===\n", ep.name)
		fmt.Printf("Endpoint: %s\n\n", ep.endpoint)

		data, err := client.DebugRawRequest(ctx, ep.endpoint)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}

		var pretty bytes.Buffer
		if err := json.Indent(&pretty, data, "", "  "); err != nil {
			// Not valid JSON, print raw
			fmt.Printf("Raw response:\n%s\n", string(data))
		} else {
			fmt.Println(pretty.String())
		}
	}
}
