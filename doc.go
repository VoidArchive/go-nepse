// Package nepse provides a type-safe Go client for NEPSE market data.
//
// Covers market summaries, securities, prices, floor sheets, indices,
// company fundamentals, dividends, and corporate actions.
//
// Example:
//
//	opts := nepse.DefaultOptions()
//	opts.TLSVerification = false // NEPSE servers have certificate issues
//
//	client, err := nepse.NewClient(opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	summary, err := client.MarketSummary(context.Background())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Turnover: Rs. %.2f\n", summary.TotalTurnover)
package nepse

const (
	Version   = "0.2.0"
	UserAgent = "go-nepse/" + Version
)

// Date formats used by NEPSE API.
const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// Sector names.
const (
	SectorBanking          = "Banking"
	SectorDevelopmentBank  = "Development Bank"
	SectorFinance          = "Finance"
	SectorHotelTourism     = "Hotel Tourism"
	SectorHydro            = "Hydro"
	SectorInvestment       = "Investment"
	SectorLifeInsurance    = "Life Insurance"
	SectorManufacturing    = "Manufacturing"
	SectorMicrofinance     = "Microfinance"
	SectorMutualFund       = "Mutual Fund"
	SectorNonLifeInsurance = "Non Life Insurance"
	SectorOthers           = "Others"
	SectorTrading          = "Trading"
	SectorPromoterShare    = "Promoter Share"
)
