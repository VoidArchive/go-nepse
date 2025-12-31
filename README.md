# go-nepse

A modern, type-safe Go client library for the NEPSE (Nepal Stock Exchange) API. This library provides comprehensive access to NEPSE market data with clean architecture, proper error handling, and full type safety.

> **Disclaimer**: This is an **unofficial** library that interacts with NEPSE's undocumented internal API. It is intended for educational and personal use only. **Do not use this library for commercial projects.** The API may change without notice, and there are no guarantees of accuracy, reliability, or availability. Use at your own risk.

![NEPSE Go Demo](assets/market_overview.png)

<p align="center">
  <img src="assets/top_lists.png" width="49%" />
  <img src="assets/sector_distribution.png" width="49%" />
</p>

## Features

- **Type Safety** - All responses are properly typed structs
- **Automatic Authentication** - Token management handled transparently
- **Retry Logic** - Built-in retry with exponential backoff
- **Context Support** - Full `context.Context` support for cancellation and timeouts
- **Error Handling** - Structured error types with proper error chains

## Installation

```bash
go get github.com/voidarchive/go-nepse
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/voidarchive/go-nepse"
)

func main() {
    // Create a new NEPSE client
    opts := nepse.DefaultOptions()
    opts.TLSVerification = false // Required due to NEPSE server TLS issues

    client, err := nepse.NewClient(opts)
    if err != nil {
        log.Fatalf("Failed to create NEPSE client: %v", err)
    }
    defer client.Close()

    ctx := context.Background()

    // Get market summary
    summary, err := client.GetMarketSummary(ctx)
    if err != nil {
        log.Fatalf("Failed to get market summary: %v", err)
    }
    fmt.Printf("Total Turnover: Rs. %.2f\n", summary.TotalTurnover)
    fmt.Printf("Total Transactions: %.0f\n", summary.TotalTransactions)

    // Get company details by symbol
    details, err := client.GetCompanyDetailsBySymbol(ctx, "NABIL")
    if err != nil {
        log.Fatalf("Failed to get company details: %v", err)
    }
    fmt.Printf("Company: %s, LTP: Rs. %.2f\n", details.SecurityName, details.LastTradedPrice)
}
```

## API Coverage

### Market Data

| Method | Description |
|--------|-------------|
| `GetMarketSummary()` | Overall market statistics (turnover, volume, capitalization) |
| `GetMarketStatus()` | Current market open/close status |
| `GetNepseIndex()` | Main NEPSE index with current value and 52-week range |
| `GetNepseSubIndices()` | All sector sub-indices (Note: API currently returns empty) |
| `GetLiveMarket()` | Real-time price and volume data |
| `GetSupplyDemand()` | Aggregate supply and demand data |

### Securities & Companies

| Method | Description |
|--------|-------------|
| `GetSecurityList()` | All tradable securities |
| `GetCompanyList()` | All listed companies with sector info |
| `GetCompanyDetails(id)` | Comprehensive info including price data |
| `GetCompanyDetailsBySymbol(symbol)` | Same as above, by ticker symbol |
| `GetSectorScrips()` | Securities grouped by sector |
| `FindSecurity(id)` / `FindSecurityBySymbol(symbol)` | Find security by ID or symbol |

### Price & Trading Data

| Method | Description |
|--------|-------------|
| `GetTodaysPrices(date)` | Price data for all securities on a date |
| `GetPriceVolumeHistory(id, start, end)` | Historical OHLCV data |
| `GetPriceVolumeHistoryBySymbol(symbol, start, end)` | Same as above, by symbol |
| `GetMarketDepth(id)` / `GetMarketDepthBySymbol(symbol)` | Order book (bid/ask levels) |
| `GetFloorSheet()` | All trades for current day |
| `GetFloorSheetOf(id, date)` / `GetFloorSheetBySymbol(symbol, date)` | Trades for specific security |

### Top Lists

| Method | Description |
|--------|-------------|
| `GetTopGainers()` | Securities with highest % gains |
| `GetTopLosers()` | Securities with highest % losses |
| `GetTopTenTrade()` | Top by traded share volume |
| `GetTopTenTransaction()` | Top by transaction count |
| `GetTopTenTurnover()` | Top by trading turnover |

### Graph Data

| Method | Description |
|--------|-------------|
| `GetDailyIndexGraph(indexType)` | Intraday graph for any index type |
| `GetDailyNepseIndexGraph()` | Main NEPSE index chart |
| `GetDailyScripPriceGraph(id)` | Intraday chart for a security |

> **Note**: Graph endpoints currently return empty data. Use `GetPriceVolumeHistory` for charting.

## Configuration

```go
opts := nepse.DefaultOptions()
opts.TLSVerification = false  // Required due to NEPSE server TLS issues
opts.HTTPTimeout = 30 * time.Second
opts.MaxRetries = 3
opts.RetryDelay = time.Second

client, err := nepse.NewClient(opts)
```

### TLS Verification

The `TLSVerification: false` option is required due to TLS configuration issues on NEPSE's servers. This is a known limitation of the NEPSE API infrastructure.

## Error Handling

The library provides structured error types:

```go
import "errors"

data, err := client.GetMarketSummary(ctx)
if err != nil {
    var nepseErr *nepse.NepseError
    if errors.As(err, &nepseErr) {
        switch nepseErr.Type {
        case nepse.ErrorTypeNotFound:
            // Handle not found
        case nepse.ErrorTypeNetwork:
            // Handle network issues
        case nepse.ErrorTypeRateLimit:
            // Handle rate limiting
        }
    }
}
```

## Production Checklist

Before using this library in any production-like environment, ensure you have addressed the following:

- [ ] **Accept Undocumented API Risks**: Acknowledge that this library uses an unofficial API. It **will** break if NEPSE updates their infrastructure.
- [ ] **TLS Security**: Address the `TLSVerification: false` requirement. In production, consider routing requests through a secure proxy that handles the connection to NEPSE.
- [ ] **Caching Strategy**: Implement application-level caching. The NEPSE API is sensitive to high traffic and may block requests if hit too frequently.
- [ ] **Error Monitoring**: Listen for `nepse.NepseError` and set up alerts for `ErrorTypeNetwork` or `ErrorTypeInternal` which often indicate API changes.
- [ ] **Rate Limiting**: Ensure your application respects NEPSE's implicit rate limits to avoid IP blocks.
- [ ] **Fallback Plan**: Have a manual or alternative data source strategy for when the API is unavailable during trading hours.


## Examples

See `_examples/` for complete usage examples:

```bash
# Run the basic example
go run _examples/basic/main.go

# Include graph endpoints
go run _examples/basic/main.go --with-graphs

# Include floor sheet data
go run _examples/basic/main.go --with-floorsheet
```

## Requirements

- Go 1.25+

## License

MIT License - see [LICENSE](LICENSE) for details.
