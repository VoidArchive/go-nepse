# go-nepse

Type-safe Go client for NEPSE market data.

> **Disclaimer**: Unofficial library using NEPSE's undocumented API. For educational/personal use only. May break without notice.

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
    opts.TLSVerification = false // NEPSE servers have certificate issues

    client, err := nepse.NewClient(opts)
    if err != nil {
        log.Fatalf("Failed to create NEPSE client: %v", err)
    }
    defer client.Close()

    ctx := context.Background()

    // Get market summary
    summary, err := client.MarketSummary(ctx)
    if err != nil {
        log.Fatalf("Failed to get market summary: %v", err)
    }
    fmt.Printf("Total Turnover: Rs. %.2f\n", summary.TotalTurnover)
    fmt.Printf("Total Transactions: %.0f\n", summary.TotalTransactions)

    // Get company details by symbol
    details, err := client.CompanyBySymbol(ctx, "NABIL")
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
| `MarketSummary()` | Overall market statistics (turnover, volume, capitalization) |
| `MarketStatus()` | Current market open/close status |
| `NepseIndex()` | Main NEPSE index with current value and 52-week range |
| `SubIndices()` | All sector sub-indices (Note: API currently returns empty) |
| `LiveMarket()` | Real-time price and volume data |
| `SupplyDemand()` | Aggregate supply and demand data |

### Securities & Companies

| Method | Description |
|--------|-------------|
| `Securities()` | All tradable securities |
| `Companies()` | All listed companies with sector info |
| `Company(id)` | Comprehensive info including price data |
| `CompanyBySymbol(symbol)` | Same as above, by ticker symbol |
| `SectorScrips()` | Securities grouped by sector |
| `FindSecurity(id)` / `FindSecurityBySymbol(symbol)` | Find security by ID or symbol |

### Price & Trading Data

| Method | Description |
|--------|-------------|
| `TodaysPrices(date)` | Price data for all securities on a date |
| `PriceHistory(id, start, end)` | Historical OHLCV data |
| `PriceHistoryBySymbol(symbol, start, end)` | Same as above, by symbol |
| `MarketDepth(id)` / `MarketDepthBySymbol(symbol)` | Order book (bid/ask levels) |
| `FloorSheet()` | All trades for current day |
| `FloorSheetOf(id, date)` / `FloorSheetBySymbol(symbol, date)` | Trades for specific security |

### Top Lists

| Method | Description |
|--------|-------------|
| `TopGainers()` | Securities with highest % gains |
| `TopLosers()` | Securities with highest % losses |
| `TopTenTrade()` | Top by traded share volume |
| `TopTenTransaction()` | Top by transaction count |
| `TopTenTurnover()` | Top by trading turnover |

### Graph Data

| Method | Description |
|--------|-------------|
| `DailyIndexGraph(indexType)` | Intraday graph for any index type |
| `DailyNepseIndexGraph()` | Main NEPSE index chart |
| `DailyScripGraph(id)` | Intraday chart for a security |

### Company Fundamentals

| Method | Description |
|--------|-------------|
| `CompanyProfile(id)` / `CompanyProfileBySymbol(symbol)` | Detailed company profile (contact, address, etc.) |
| `BoardOfDirectors(id)` / `BoardOfDirectorsBySymbol(symbol)` | Board members with designations |
| `Reports(id)` / `ReportsBySymbol(symbol)` | Quarterly & annual reports with PE, EPS, book value |
| `CorporateActions(id)` / `CorporateActionsBySymbol(symbol)` | Bonus shares, rights issues (approval/distribution) |
| `Dividends(id)` / `DividendsBySymbol(symbol)` | Dividend declarations (cash & bonus) |

> **Note**: Corporate Actions tracks when bonus/rights shares were *approved/distributed*, while Dividends tracks when they were *declared*. There may be a 1-year lag between declaration and distribution.

## Configuration

```go
opts := nepse.DefaultOptions()
opts.TLSVerification = false // NEPSE servers have certificate issues
opts.HTTPTimeout = 30 * time.Second
opts.MaxRetries = 3
opts.RetryDelay = time.Second

client, err := nepse.NewClient(opts)
```

## Error Handling

The library provides structured error types:

```go
import "errors"

data, err := client.MarketSummary(ctx)
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

- [ ] **API Risks**: Unofficial API, will break when NEPSE updates infrastructure
- [ ] **TLS Security**: Route through secure proxy to handle certificate issues
- [ ] **Caching**: Cache responses to avoid rate limiting
- [ ] **Monitoring**: Alert on `ErrorTypeNetwork` or `ErrorTypeInternal`
- [ ] **Rate Limiting**: Respect implicit limits to avoid IP blocks
- [ ] **Fallback**: Have alternative data source for outages


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
