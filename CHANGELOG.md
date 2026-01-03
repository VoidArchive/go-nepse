# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Unit tests for core functionality
- Integration tests
- Rate limiting improvements

## [0.2.0] - 2026-01-03

### Added
- **Company Fundamentals**: `CompanyProfile`, `BoardOfDirectors`, `Reports`, `CorporateActions`, `Dividends` endpoints (with BySymbol variants)
- **New Types**: `CompanyProfile`, `BoardMember`, `Report`, `FiscalReport`, `CorporateAction`, `Dividend`, `DividendNotice`
- **Helper Methods**: `FullName()`, `IsAnnual()`, `IsQuarterly()`, `HasCashDividend()`, `HasBonusDividend()`
- **Debug Utilities**: `_examples/debug/main.go` for API discovery

### Changed
- Updated examples and docs with Company Fundamentals section
- Added note on Corporate Actions vs Dividends timing difference

### Closes
- Issue #1: Dividends, reports, news endpoints
- Issue #2: Company fundamentals (PE, EPS, book value)

## [0.1.2] - 2026-01-01

### Added
- **Graph Endpoints**: Dynamic POST payload ID generation
- **Data Unmarshaling**: Robust graph data point handling
- **Authentication**: Token management with retry mechanisms

### Changed
- **BREAKING**: Removed `Get` prefix from methods (`MarketSummary()` instead of `GetMarketSummary()`)
- Improved docs for `TodaysPrices`, `FloorSheetOf`, `FloorSheetBySymbol`

### Fixed
- Token management robustness during network issues

## [0.1.1] - 2025-12-31

### Fixed
- Missing `/` in endpoint URLs for `GetCompanyDetails`, `GetPriceVolumeHistory`, `GetMarketDepth`, `GetFloorSheetOf`, `GetDailyScripPriceGraph`
- Response type mismatches in `GetSupplyDemand`, `PriceHistory`, `MarketDepth`
- `GetSectorScrips` now uses company list with sector info
- `GetNepseSubIndices` filtering logic
- Pointer-to-loop-variable bugs, swallowed errors, lint issues
- URL encoding via `url.Values`

### Changed
- **BREAKING**: `GetSupplyDemand()` returns `*SupplyDemandData`
- Added NEPSE index ID constants
- Refactored `graphs.go` with `IndexType` enum

### Removed
- `SupplyDemandEntry` type

## [0.0.1] - 2025-12-30

### Added
- Initial release with 50+ endpoints
- Type-safe dual API pattern (ID and Symbol methods)
- WASM-based token management
- Market data: summaries, securities, prices, floor sheets, top lists, indices
- Retry logic with exponential backoff
- Context support, structured errors, connection pooling

### Performance
- `GetSectorScrips()` ~75% faster (5-10s → ~2.3s)
- Single API call for sector data

### Security
- Configurable TLS, secure token lifecycle, input validation

### Known Issues
- Graph endpoints return empty (server-side issue)
