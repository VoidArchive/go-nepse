package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Find AHPC security ID
	companies, err := client.Companies(ctx)
	if err != nil {
		log.Fatalf("Failed to get companies: %v", err)
	}

	var ahpcID int32
	for _, c := range companies {
		if strings.EqualFold(c.Symbol, "AHPC") {
			ahpcID = c.ID
			fmt.Printf("Found AHPC: ID=%d, Name=%s\n", c.ID, c.CompanyName)
			break
		}
	}

	if ahpcID == 0 {
		log.Fatal("AHPC not found")
	}

	// Get the auth token from the library
	token, err := client.DebugDecodedToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Printf("Got auth token (length: %d)\n", len(token))

	// Make direct POST request
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Try POST with no body
	url := fmt.Sprintf("https://nepalstock.com/api/nots/security/%d", ahpcID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Salter "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	req.Header.Set("Origin", "https://nepalstock.com")
	req.Header.Set("Referer", "https://nepalstock.com/")

	fmt.Printf("\n=== POST %s ===\n", url)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Pretty print if JSON
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Printf("Raw response: %s\n", string(respBody))
	} else {
		// Check for share structure fields
		if _, ok := result["stockListedShares"]; ok {
			fmt.Println("\nFound share structure data!")
			fmt.Printf("stockListedShares: %v\n", result["stockListedShares"])
			fmt.Printf("promoterShares: %v\n", result["promoterShares"])
			fmt.Printf("publicShares: %v\n", result["publicShares"])
			fmt.Printf("promoterPercentage: %v\n", result["promoterPercentage"])
			fmt.Printf("publicPercentage: %v\n", result["publicPercentage"])
			fmt.Printf("paidUpCapital: %v\n", result["paidUpCapital"])
			fmt.Printf("marketCapitalization: %v\n", result["marketCapitalization"])
		} else {
			pretty, _ := json.MarshalIndent(result, "", "  ")
			output := string(pretty)
			if len(output) > 3000 {
				output = output[:3000] + "\n... (truncated)"
			}
			fmt.Println(output)
		}
	}
}
