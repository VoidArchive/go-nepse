package nepse

// DefaultBaseURL is the production NEPSE API URL.
const DefaultBaseURL = "https://nepalstock.com.np"

// Endpoints holds all NEPSE API endpoint paths.
type Endpoints struct {
	// Market data
	MarketSummary string
	MarketOpen    string
	LiveMarket    string
	SupplyDemand  string
	TodaysPrice   string
	FloorSheet    string

	// Index data
	NepseIndex string

	// Top ten lists
	TopGainers      string
	TopLosers       string
	TopTrade        string
	TopTransaction  string
	TopTurnover     string

	// Security/Company data
	SecurityList        string
	CompanyList         string
	CompanyDetails      string
	CompanyPriceHistory string
	CompanyFloorsheet   string
	MarketDepth         string

	// Company fundamentals & financial data
	CompanyProfile   string
	BoardOfDirectors string
	CorporateActions string
	Reports          string
	Dividend         string

	// Graph endpoints (index charts)
	GraphNepseIndex          string
	GraphSensitiveIndex      string
	GraphFloatIndex          string
	GraphSensitiveFloatIndex string
	GraphBankingSubindex     string
	GraphDevBankSubindex     string
	GraphFinanceSubindex     string
	GraphHotelSubindex       string
	GraphHydroSubindex       string
	GraphInvestmentSubindex  string
	GraphLifeInsSubindex     string
	GraphManufacturingSubindex  string
	GraphMicrofinanceSubindex   string
	GraphMutualFundSubindex     string
	GraphNonLifeInsSubindex     string
	GraphOthersSubindex         string
	GraphTradingSubindex        string

	// Graph endpoints (company)
	CompanyDailyGraph string
}

// Config holds configuration for the NEPSE API client.
type Config struct {
	BaseURL   string
	Endpoints Endpoints
}

// DefaultEndpoints returns the default NEPSE API endpoints.
func DefaultEndpoints() Endpoints {
	return Endpoints{
		// Market data
		MarketSummary: "/api/nots/market-summary",
		MarketOpen:    "/api/nots/nepse-data/market-open",
		LiveMarket:    "/api/nots/lives-market",
		SupplyDemand:  "/api/nots/nepse-data/supplydemand",
		TodaysPrice:   "/api/nots/nepse-data/today-price",
		FloorSheet:    "/api/nots/nepse-data/floorsheet",

		// Index data
		NepseIndex: "/api/nots/nepse-index",

		// Top ten lists
		TopGainers:     "/api/nots/top-ten/top-gainer",
		TopLosers:      "/api/nots/top-ten/top-loser",
		TopTrade:       "/api/nots/top-ten/trade",
		TopTransaction: "/api/nots/top-ten/transaction",
		TopTurnover:    "/api/nots/top-ten/turnover",

		// Security/Company data
		SecurityList:        "/api/nots/security?nonDelisted=true",
		CompanyList:         "/api/nots/company/list",
		CompanyDetails:      "/api/nots/security",
		CompanyPriceHistory: "/api/nots/market/history/security",
		CompanyFloorsheet:   "/api/nots/security/floorsheet",
		MarketDepth:         "/api/nots/nepse-data/marketdepth",

		// Company fundamentals & financial data
		CompanyProfile:   "/api/nots/security/profile",
		BoardOfDirectors: "/api/nots/security/boardOfDirectors",
		CorporateActions: "/api/nots/security/corporate-actions",
		Reports:          "/api/nots/application/reports",
		Dividend:         "/api/nots/application/dividend",

		// Graph endpoints (index charts)
		GraphNepseIndex:            "/api/nots/graph/index/58",
		GraphSensitiveIndex:        "/api/nots/graph/index/57",
		GraphFloatIndex:            "/api/nots/graph/index/62",
		GraphSensitiveFloatIndex:   "/api/nots/graph/index/63",
		GraphBankingSubindex:       "/api/nots/graph/index/51",
		GraphDevBankSubindex:       "/api/nots/graph/index/55",
		GraphFinanceSubindex:       "/api/nots/graph/index/60",
		GraphHotelSubindex:         "/api/nots/graph/index/52",
		GraphHydroSubindex:         "/api/nots/graph/index/54",
		GraphInvestmentSubindex:    "/api/nots/graph/index/67",
		GraphLifeInsSubindex:       "/api/nots/graph/index/65",
		GraphManufacturingSubindex: "/api/nots/graph/index/56",
		GraphMicrofinanceSubindex:  "/api/nots/graph/index/64",
		GraphMutualFundSubindex:    "/api/nots/graph/index/66",
		GraphNonLifeInsSubindex:    "/api/nots/graph/index/59",
		GraphOthersSubindex:        "/api/nots/graph/index/53",
		GraphTradingSubindex:       "/api/nots/graph/index/61",

		// Graph endpoints (company)
		CompanyDailyGraph: "/api/nots/market/graphdata/daily",
	}
}

// DefaultConfig returns the default NEPSE API configuration.
func DefaultConfig() *Config {
	return &Config{
		BaseURL:   DefaultBaseURL,
		Endpoints: DefaultEndpoints(),
	}
}
