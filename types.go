package nepse

import "encoding/json"

// MarketSummaryItem represents a single item in the market summary response.
type MarketSummaryItem struct {
	Detail string  `json:"detail"`
	Value  float64 `json:"value"`
}

// MarketSummary represents the processed market summary data.
type MarketSummary struct {
	TotalTurnover             float64
	TotalTradedShares         float64
	TotalTransactions         float64
	TotalScripsTraded         float64
	TotalMarketCapitalization float64
	TotalFloatMarketCap       float64
}

// MarketStatus represents the current market status.
type MarketStatus struct {
	IsOpen string `json:"isOpen"`
	AsOf   string `json:"asOf"`
	ID     int32  `json:"id"`
}

// IsMarketOpen returns true if the market is currently open.
func (m *MarketStatus) IsMarketOpen() bool {
	return m.IsOpen == "OPEN"
}

// NepseIndexRaw represents the raw NEPSE index response item.
type NepseIndexRaw struct {
	ID               int32   `json:"id"`
	Index            string  `json:"index"`
	Close            float64 `json:"close"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	PerChange        float64 `json:"perChange"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// NepseIndex represents the NEPSE main index (ID 58).
type NepseIndex struct {
	IndexValue       float64 `json:"close"`
	PercentChange    float64 `json:"perChange"`
	PointChange      float64 `json:"change"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// SubIndex represents a sector sub-index.
type SubIndex struct {
	ID               int32   `json:"id"`
	Index            string  `json:"index"`
	Close            float64 `json:"close"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	PerChange        float64 `json:"perChange"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// Security represents a listed security/company.
// Note: The security list API only returns id, symbol, securityName, and activeStatus.
// For sector info, use GetCompanyList() instead.
type Security struct {
	ID           int32  `json:"id"`
	Symbol       string `json:"symbol"`
	SecurityName string `json:"securityName"`
	ActiveStatus string `json:"activeStatus"`
}

// Company represents company information.
type Company struct {
	ID             int32  `json:"id"`
	CompanyName    string `json:"companyName"`
	Symbol         string `json:"symbol"`
	SecurityName   string `json:"securityName"`
	Status         string `json:"status"`
	CompanyEmail   string `json:"companyEmail"`
	Website        string `json:"website"`
	SectorName     string `json:"sectorName"`
	RegulatoryBody string `json:"regulatoryBody"`
	InstrumentType string `json:"instrumentType"`
}

// InstrumentType represents the type of financial instrument.
type InstrumentType struct {
	ID           int32  `json:"id"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	ActiveStatus string `json:"activeStatus"`
}

// ParsedValue represents a value with its source string and parsed numeric value.
type ParsedValue struct {
	Source      string  `json:"source"`
	ParsedValue float64 `json:"parsedValue"`
}

// ShareGroup represents the share group classification.
type ShareGroup struct {
	ID              int32   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	CapitalRangeMin int64   `json:"capitalRangeMin"`
	ModifiedBy      *string `json:"modifiedBy"`
	ModifiedDate    *string `json:"modifiedDate"`
	ActiveStatus    string  `json:"activeStatus"`
	IsDefault       string  `json:"isDefault"`
}

// SectorMaster represents sector information.
type SectorMaster struct {
	ID                int32  `json:"id"`
	SectorDescription string `json:"sectorDescription"`
	ActiveStatus      string `json:"activeStatus"`
	RegulatoryBody    string `json:"regulatoryBody"`
}

// CompanyInfo represents company information within classification.
type CompanyInfo struct {
	ID                        int32        `json:"id"`
	CompanyShortName          string       `json:"companyShortName"`
	CompanyName               string       `json:"companyName"`
	Email                     string       `json:"email"`
	CompanyWebsite            string       `json:"companyWebsite"`
	CompanyContactPerson      string       `json:"companyContactPerson"`
	SectorMaster              SectorMaster `json:"sectorMaster"`
	CompanyRegistrationNumber string       `json:"companyRegistrationNumber"`
	ActiveStatus              string       `json:"activeStatus"`
}

// TodayPrice represents today's price data for a security.
type TodayPrice struct {
	ID                  int32   `json:"id"`
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closePrice"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	PreviousClose       float64 `json:"previousClose"`
	DifferenceRs        float64 `json:"differenceRs"`
	PercentageChange    float64 `json:"percentageChange"`
	TotalTrades         int32   `json:"totalTrades"`
	BusinessDate        string  `json:"businessDate"`
	SecurityID          int32   `json:"securityId"`
	LastTradedPrice     float64 `json:"lastTradedPrice"`
	MaxPrice            float64 `json:"maxPrice"`
	MinPrice            float64 `json:"minPrice"`
}

// PriceHistory represents historical OHLCV data for a security.
// Note: NEPSE API does not provide open price in historical data.
type PriceHistory struct {
	BusinessDate        string  `json:"businessDate"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closePrice"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	TotalTrades         int32   `json:"totalTrades"`
}

// FloorSheetEntry represents a single floor sheet entry.
type FloorSheetEntry struct {
	ContractID       int64   `json:"contractId"`
	StockSymbol      string  `json:"stockSymbol"`
	SecurityName     string  `json:"securityName"`
	BuyerMemberID    int32   `json:"buyerMemberId"`
	SellerMemberID   int32   `json:"sellerMemberId"`
	ContractQuantity int64   `json:"contractQuantity"`
	ContractRate     float64 `json:"contractRate"`
	BusinessDate     string  `json:"businessDate"`
	TradeTime        string  `json:"tradeTime"`
	SecurityID       int32   `json:"securityId"`
	ContractAmount   float64 `json:"contractAmount"`
	BuyerBrokerName  string  `json:"buyerBrokerName"`
	SellerBrokerName string  `json:"sellerBrokerName"`
	TradeBookID      int64   `json:"tradeBookId"`
}

// FloorSheetResponse represents the paginated floor sheet response.
type FloorSheetResponse struct {
	FloorSheets struct {
		Content          []FloorSheetEntry `json:"content"`
		PageNumber       int32             `json:"number"`
		Size             int32             `json:"size"`
		TotalElements    int64             `json:"totalElements"`
		TotalPages       int32             `json:"totalPages"`
		First            bool              `json:"first"`
		Last             bool              `json:"last"`
		NumberOfElements int32             `json:"numberOfElements"`
	} `json:"floorsheets"`
}

// DepthEntry represents a single entry in market depth.
type DepthEntry struct {
	StockID  int32   `json:"stockId"`
	Price    float64 `json:"orderBookOrderPrice"`
	Quantity int64   `json:"quantity"`
	Orders   int32   `json:"orderCount"`
	IsBuy    int     `json:"isBuy"`
}

// MarketDepthRaw represents the raw API response for market depth.
type MarketDepthRaw struct {
	TotalBuyQty  int64 `json:"totalBuyQty"`
	TotalSellQty int64 `json:"totalSellQty"`
	MarketDepth  struct {
		BuyList  []DepthEntry `json:"buyMarketDepthList"`
		SellList []DepthEntry `json:"sellMarketDepthList"`
	} `json:"marketDepth"`
}

// MarketDepth represents processed market depth information.
type MarketDepth struct {
	TotalBuyQty  int64
	TotalSellQty int64
	BuyDepth     []DepthEntry
	SellDepth    []DepthEntry
}

// TopGainerLoserEntry represents entries in top gainers/losers lists.
type TopGainerLoserEntry struct {
	Symbol           string  `json:"symbol"`
	SecurityName     string  `json:"securityName"`
	SecurityID       int32   `json:"securityId"`
	LTP              float64 `json:"ltp"`
	PointChange      float64 `json:"pointChange"`
	PercentageChange float64 `json:"percentageChange"`
}

// TopTradeEntry represents entries in top trade (volume) list.
type TopTradeEntry struct {
	Symbol       string  `json:"symbol"`
	SecurityName string  `json:"securityName"`
	SecurityID   int32   `json:"securityId"`
	ShareTraded  int64   `json:"shareTraded"`
	ClosingPrice float64 `json:"closingPrice"`
}

// TopTurnoverEntry represents entries in top turnover list.
type TopTurnoverEntry struct {
	Symbol       string  `json:"symbol"`
	SecurityName string  `json:"securityName"`
	SecurityID   int32   `json:"securityId"`
	Turnover     float64 `json:"turnover"`
	ClosingPrice float64 `json:"closingPrice"`
}

// TopTransactionEntry represents entries in top transaction list.
type TopTransactionEntry struct {
	Symbol          string  `json:"symbol"`
	SecurityName    string  `json:"securityName"`
	SecurityID      int32   `json:"securityId"`
	TotalTrades     int32   `json:"totalTrades"`
	LastTradedPrice float64 `json:"lastTradedPrice"`
}

// TopListEntry is deprecated, use specific types instead.
// Kept for backward compatibility.
type TopListEntry struct {
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	ClosePrice          float64 `json:"closePrice"`
	PercentageChange    float64 `json:"percentageChange"`
	DifferenceRs        float64 `json:"differenceRs"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	TotalTrades         int32   `json:"totalTrades"`
	HighPrice           float64 `json:"highPrice,omitempty"`
	LowPrice            float64 `json:"lowPrice,omitempty"`
	OpenPrice           float64 `json:"openPrice,omitempty"`
	PreviousClose       float64 `json:"previousClose,omitempty"`
}

// GraphDataPoint represents a single data point in graph data.
// The NEPSE API returns different formats for different endpoints:
// - Index graphs: [timestamp, value] arrays
// - Scrip graphs: {"time": timestamp, "value": value} objects
type GraphDataPoint struct {
	Timestamp int64
	Value     float64
}

// UnmarshalJSON implements custom unmarshaling for GraphDataPoint.
// Handles both array format [timestamp, value] and object format {"time": ..., "value": ...}.
func (g *GraphDataPoint) UnmarshalJSON(data []byte) error {
	// Try array format first (index graphs): [timestamp, value]
	var arr [2]float64
	if err := json.Unmarshal(data, &arr); err == nil {
		g.Timestamp = int64(arr[0])
		g.Value = arr[1]
		return nil
	}

	// Try object format (scrip graphs): {"time": ..., "value": ...}
	var obj struct {
		Time  int64   `json:"time"`
		Value float64 `json:"value"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	g.Timestamp = obj.Time
	g.Value = obj.Value
	return nil
}

// GraphResponse represents graph data response.
type GraphResponse struct {
	Data []GraphDataPoint
}

// CompanyDetailsRaw represents the raw nested company details response.
type CompanyDetailsRaw struct {
	SecurityMcsData struct {
		SecurityID          string  `json:"securityId"`
		OpenPrice           float64 `json:"openPrice"`
		HighPrice           float64 `json:"highPrice"`
		LowPrice            float64 `json:"lowPrice"`
		TotalTradeQuantity  int64   `json:"totalTradeQuantity"`
		TotalTrades         int32   `json:"totalTrades"`
		LastTradedPrice     float64 `json:"lastTradedPrice"`
		PreviousClose       float64 `json:"previousClose"`
		BusinessDate        string  `json:"businessDate"`
		ClosePrice          float64 `json:"closePrice"`
		FiftyTwoWeekHigh    float64 `json:"fiftyTwoWeekHigh"`
		FiftyTwoWeekLow     float64 `json:"fiftyTwoWeekLow"`
		LastUpdatedDateTime string  `json:"lastUpdatedDateTime"`
	} `json:"securityMcsData"`
	SecurityData struct {
		ID               int32  `json:"id"`
		Symbol           string `json:"symbol"`
		SecurityName     string `json:"securityName"`
		ActiveStatus     string `json:"activeStatus"`
		PermittedToTrade string `json:"permittedToTrade"`
		Email            string `json:"email"`
		Sector           string `json:"sector"`
	} `json:"securityData"`
}

// CompanyDetails represents processed company information.
type CompanyDetails struct {
	ID               int32  `json:"id"`
	Symbol           string `json:"symbol"`
	SecurityName     string `json:"securityName"`
	SectorName       string `json:"sectorName"`
	Email            string `json:"email"`
	ActiveStatus     string `json:"activeStatus"`
	PermittedToTrade string `json:"permittedToTrade"`

	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closePrice"`
	LastTradedPrice     float64 `json:"lastTradedPrice"`
	PreviousClose       float64 `json:"previousClose"`
	TotalTradeQuantity  int64   `json:"totalTradeQuantity"`
	TotalTrades         int32   `json:"totalTrades"`
	FiftyTwoWeekHigh    float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow     float64 `json:"fiftyTwoWeekLow"`
	BusinessDate        string  `json:"businessDate"`
	LastUpdatedDateTime string  `json:"lastUpdatedDateTime"`
}

// LiveMarketEntry represents live market data entry.
type LiveMarketEntry struct {
	SecurityID          string  `json:"securityId"`
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	LastTradedPrice     float64 `json:"lastTradedPrice"`
	TotalTradeQuantity  int64   `json:"totalTradeQuantity"`
	TotalTradeValue     float64 `json:"totalTradeValue"`
	PreviousClose       float64 `json:"previousClose"`
	PercentageChange    float64 `json:"percentageChange"`
	LastTradedVolume    int64   `json:"lastTradedVolume"`
	LastUpdatedDateTime string  `json:"lastUpdatedDateTime"`
	AverageTradedPrice  float64 `json:"averageTradedPrice"`
}

// SectorScrips represents scrips grouped by sector.
type SectorScrips map[string][]string

// PaginatedResponse represents a generic paginated response.
type PaginatedResponse[T any] struct {
	Content          []T   `json:"content"`
	PageNumber       int32 `json:"number"`
	Size             int32 `json:"size"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int32 `json:"totalPages"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int32 `json:"numberOfElements"`
}

// CompanyProfile represents detailed company profile information.
type CompanyProfile struct {
	CompanyName          string `json:"companyName"`
	CompanyEmail         string `json:"companyEmail"`
	CompanyProfile       string `json:"companyProfile"`
	CompanyContactPerson string `json:"companyContactPerson"`
	LogoFilePath         string `json:"logoFilePath"`
	AddressType          string `json:"addressType"`
	AddressField         string `json:"addressField"`
	PhoneNumber          string `json:"phoneNumber"`
	Fax                  string `json:"fax"`
	Town                 string `json:"town"`
}

// BoardMember represents a board of directors member.
type BoardMember struct {
	FirstName       string  `json:"firstName"`
	MiddleName      string  `json:"middleName"`
	LastName        string  `json:"lastName"`
	Designation     string  `json:"designation"`
	MemberPhotoPath *string `json:"memberPhotoPath"`
	Description     string  `json:"description"`
}

// FullName returns the complete name of the board member.
func (b *BoardMember) FullName() string {
	if b.MiddleName != "" {
		return b.FirstName + " " + b.MiddleName + " " + b.LastName
	}
	return b.FirstName + " " + b.LastName
}

// CorporateAction represents a corporate action (bonus, rights, cash dividend).
type CorporateAction struct {
	ActiveStatus          string   `json:"activeStatus"`
	AuthorizationComments *string  `json:"authorizationComments"`
	SubmittedDate         string   `json:"submittedDate"`
	FilePath              string   `json:"filePath"`
	DocumentID            int32    `json:"documentId"`
	RatioNum              float64  `json:"ratioNum"`
	RatioDen              float64  `json:"ratioDen"`
	CashDividend          *float64 `json:"cashDividend"`
	FiscalYear            string   `json:"fiscalYear"`
	RightAmountPerShare   *float64 `json:"rightAmountPerShare"`
	BonusPercentage       float64  `json:"bonusPercentage"`
	RightPercentage       *float64 `json:"rightPercentage"`
	SdID                  int32    `json:"sdId"`
}

// IsBonus returns true if this corporate action is a bonus share.
func (c *CorporateAction) IsBonus() bool {
	return c.BonusPercentage > 0
}

// IsRight returns true if this corporate action is a rights issue.
func (c *CorporateAction) IsRight() bool {
	return c.RightPercentage != nil && *c.RightPercentage > 0
}

// IsCashDividend returns true if this corporate action is a cash dividend.
func (c *CorporateAction) IsCashDividend() bool {
	return c.CashDividend != nil && *c.CashDividend > 0
}

// FinancialYear represents a fiscal year.
type FinancialYear struct {
	ID           int32  `json:"id"`
	FYName       string `json:"fyName"`
	FYNameNepali string `json:"fyNameNepali"`
	FromYear     string `json:"fromYear"`
	ToYear       string `json:"toYear"`
}

// QuarterMaster represents a fiscal quarter.
type QuarterMaster struct {
	ID          int32  `json:"id"`
	QuarterName string `json:"quarterName"`
}

// ReportTypeMaster represents a report type (Annual/Quarterly).
type ReportTypeMaster struct {
	ID         int32  `json:"id"`
	ReportName string `json:"reportName"`
}

// FiscalReport contains financial metrics from a report.
type FiscalReport struct {
	ID               int32             `json:"id"`
	QuarterMaster    *QuarterMaster    `json:"quarterMaster"`
	ReportTypeMaster *ReportTypeMaster `json:"reportTypeMaster"`
	FinancialYear    *FinancialYear    `json:"financialYear"`
	PEValue          float64           `json:"peValue"`
	EPSValue         float64           `json:"epsValue"`
	PaidUpCapital    float64           `json:"paidUpCapital"`
	ProfitAmount     float64           `json:"profitAmount"`
	NetWorthPerShare float64           `json:"netWorthPerShare"`
	Remarks          *string           `json:"remarks"`
}

// ReportDocument represents a document attached to a report.
type ReportDocument struct {
	ID            int32  `json:"id"`
	SubmittedDate string `json:"submittedDate"`
	FilePath      string `json:"filePath"`
	EncryptedID   string `json:"encryptedId"`
}

// Report represents a quarterly or annual financial report.
type Report struct {
	ID                             int32            `json:"id"`
	ActiveStatus                   string           `json:"activeStatus"`
	ModifiedDate                   string           `json:"modifiedDate"`
	ApplicationType                int32            `json:"applicationType"`
	ApplicationStatus              int32            `json:"applicationStatus"`
	FiscalReport                   *FiscalReport    `json:"fiscalReport"`
	ApplicationDocumentDetailsList []ReportDocument `json:"applicationDocumentDetailsList"`
}

// IsAnnual returns true if this is an annual report.
func (r *Report) IsAnnual() bool {
	if r.FiscalReport != nil && r.FiscalReport.ReportTypeMaster != nil {
		return r.FiscalReport.ReportTypeMaster.ReportName == "Annual Report"
	}
	return false
}

// IsQuarterly returns true if this is a quarterly report.
func (r *Report) IsQuarterly() bool {
	if r.FiscalReport != nil && r.FiscalReport.ReportTypeMaster != nil {
		return r.FiscalReport.ReportTypeMaster.ReportName == "Quarterly Report"
	}
	return false
}

// QuarterName returns the quarter name (e.g., "First Quarter") or empty if annual.
func (r *Report) QuarterName() string {
	if r.FiscalReport != nil && r.FiscalReport.QuarterMaster != nil {
		return r.FiscalReport.QuarterMaster.QuarterName
	}
	return ""
}

// DividendNotice contains dividend declaration details.
type DividendNotice struct {
	ID            int32          `json:"id"`
	FinancialYear *FinancialYear `json:"financialYear"`
	CashDividend  float64        `json:"cashDividend"`
	BonusShare    float64        `json:"bonusShare"`
	RightShare    float64        `json:"rightShare"`
	Remarks       *string        `json:"remarks"`
}

// CompanyNews contains news/announcement details.
type CompanyNews struct {
	ID              int32           `json:"id"`
	NewsSource      string          `json:"newsSource"`
	NewsHeadline    string          `json:"newsHeadline"`
	NewsBody        string          `json:"newsBody"`
	NewsType        string          `json:"newsType"`
	ExpiryDate      string          `json:"expiryDate"`
	DividendsNotice *DividendNotice `json:"dividendsNotice"`
}

// Dividend represents a dividend declaration.
type Dividend struct {
	ID                int32        `json:"id"`
	ActiveStatus      string       `json:"activeStatus"`
	ModifiedDate      string       `json:"modifiedDate"`
	ApplicationType   int32        `json:"applicationType"`
	ApplicationStatus int32        `json:"applicationStatus"`
	CompanyNews       *CompanyNews `json:"companyNews"`
}

// HasCashDividend returns true if this dividend includes cash.
func (d *Dividend) HasCashDividend() bool {
	if d.CompanyNews != nil && d.CompanyNews.DividendsNotice != nil {
		return d.CompanyNews.DividendsNotice.CashDividend > 0
	}
	return false
}

// HasBonusDividend returns true if this dividend includes bonus shares.
func (d *Dividend) HasBonusDividend() bool {
	if d.CompanyNews != nil && d.CompanyNews.DividendsNotice != nil {
		return d.CompanyNews.DividendsNotice.BonusShare > 0
	}
	return false
}

// CashPercentage returns the cash dividend percentage.
func (d *Dividend) CashPercentage() float64 {
	if d.CompanyNews != nil && d.CompanyNews.DividendsNotice != nil {
		return d.CompanyNews.DividendsNotice.CashDividend
	}
	return 0
}

// BonusPercentage returns the bonus dividend percentage.
func (d *Dividend) BonusPercentage() float64 {
	if d.CompanyNews != nil && d.CompanyNews.DividendsNotice != nil {
		return d.CompanyNews.DividendsNotice.BonusShare
	}
	return 0
}

// FiscalYear returns the fiscal year of the dividend.
func (d *Dividend) FiscalYear() string {
	if d.CompanyNews != nil && d.CompanyNews.DividendsNotice != nil && d.CompanyNews.DividendsNotice.FinancialYear != nil {
		return d.CompanyNews.DividendsNotice.FinancialYear.FYNameNepali
	}
	return ""
}
