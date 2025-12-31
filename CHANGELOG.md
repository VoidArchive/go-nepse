# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Unit tests for core functionality
- Integration tests
- Rate limiting improvements

## [0.1.1] - 2025-12-31

### Fixed
- **GetCompanyDetails**: Fixed missing `/` in endpoint URL causing 404 errors
- **GetPriceVolumeHistory**: Fixed missing `/` in endpoint URL causing 404 errors
- **GetMarketDepth**: Fixed missing `/` in endpoint URL and updated response type to match API
- **GetFloorSheetOf**: Fixed missing `/` in endpoint URL
- **GetDailyScripPriceGraph**: Fixed missing `/` in endpoint URL
- **GetSupplyDemand**: Changed return type to match actual API response format (`*SupplyDemandData`)
- **GetSectorScrips**: Now uses company list which includes sector information
- **GetNepseSubIndices**: Fixed filtering logic, uses map-based exclusion
- **PriceHistory**: Updated struct to match actual API response fields
- **MarketDepth**: Updated struct to match actual API response format
- Pointer-to-loop-variable bugs in `findSecurityByID` and `findSecurityBySymbol`
- Silently swallowed errors in `GetSupplyDemand` and `GetFloorSheet`
- All lint issues (unchecked error returns in `transport.go`)
- URL encoding now uses `url.Values` for proper escaping

### Changed
- **BREAKING**: `GetSupplyDemand()` now returns `*SupplyDemandData` instead of `[]SupplyDemandEntry`
- Added constants for NEPSE index IDs (58, 57, 62, 63)
- Improved all doc comments for clarity and accuracy
- Refactored `graphs.go` to use `IndexType` enum and reduce duplication

### Removed
- Unused `SupplyDemandEntry` type (replaced by `SupplyDemandData`)

## [0.0.1] - 2025-12-30

### Added
- Initial release of NEPSE Go Library
- Complete NEPSE API coverage with 50+ endpoints
- Type-safe dual API pattern (ID and Symbol methods)
- Automatic token management with WASM-based authentication
- Comprehensive market data access:
  - Market summaries and status
  - Security and company information
  - Price and trading data
  - Floor sheet data
  - Top lists (gainers, losers, etc.)
  - Market indices and sub-indices
- Fast sector grouping without API calls
- Built-in retry logic with exponential backoff
- Context support for timeouts and cancellation
- Structured error handling with typed errors
- Connection pooling and HTTP optimization
- Production-ready examples and documentation

### Performance Improvements
- Optimized `GetSectorScrips()` method — ~75% faster execution (5–10s → ~2.3s)
- Reduced API calls from 50+ to 1 for sector data
- HTTP connection pooling for efficiency
- Smart retry mechanisms

### Security
- Comprehensive security audit completed
- Configurable TLS verification (may be required due to NEPSE server issues)
- Secure token lifecycle management with automatic refresh
- Input validation across all public APIs
- Safe error handling with no sensitive data leakage
- No sensitive data exposed in logs

### Documentation
- Comprehensive README with full API coverage
- GoDoc comments for all public APIs
- Working examples in `cmd/examples/`
- Security best practices documented

### Known Issues
- Graph endpoints currently return empty data due to NEPSE backend issues
- Issue is server-side, not a client library defect
