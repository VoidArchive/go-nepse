package nepse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Index IDs used by NEPSE API.
const (
	nepseIndexID          = 58
	sensitiveIndexID      = 57
	floatIndexID          = 62
	sensitiveFloatIndexID = 63
)

// MarketSummary returns aggregate market statistics including turnover, volume, and capitalization.
func (c *Client) MarketSummary(ctx context.Context) (*MarketSummary, error) {
	var rawItems []MarketSummaryItem
	if err := c.apiRequest(ctx, c.config.Endpoints.MarketSummary, &rawItems); err != nil {
		return nil, err
	}

	summary := &MarketSummary{}
	for _, item := range rawItems {
		switch item.Detail {
		case "Total Turnover Rs:":
			summary.TotalTurnover = item.Value
		case "Total Traded Shares":
			summary.TotalTradedShares = item.Value
		case "Total Transactions":
			summary.TotalTransactions = item.Value
		case "Total Scrips Traded":
			summary.TotalScripsTraded = item.Value
		case "Total Market Capitalization Rs:":
			summary.TotalMarketCapitalization = item.Value
		case "Total Float Market Capitalization Rs:":
			summary.TotalFloatMarketCap = item.Value
		}
	}

	return summary, nil
}

// MarketStatus returns whether the market is currently open or closed.
func (c *Client) MarketStatus(ctx context.Context) (*MarketStatus, error) {
	var status MarketStatus
	if err := c.apiRequest(ctx, c.config.Endpoints.MarketOpen, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// NepseIndex returns the main NEPSE index with current value, change, and 52-week range.
func (c *Client) NepseIndex(ctx context.Context) (*NepseIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := c.apiRequest(ctx, c.config.Endpoints.NepseIndex, &rawIndices); err != nil {
		return nil, err
	}

	for i := range rawIndices {
		if rawIndices[i].ID == nepseIndexID {
			return &NepseIndex{
				IndexValue:       rawIndices[i].Close,
				PercentChange:    rawIndices[i].PerChange,
				PointChange:      rawIndices[i].Change,
				High:             rawIndices[i].High,
				Low:              rawIndices[i].Low,
				PreviousClose:    rawIndices[i].PreviousClose,
				FiftyTwoWeekHigh: rawIndices[i].FiftyTwoWeekHigh,
				FiftyTwoWeekLow:  rawIndices[i].FiftyTwoWeekLow,
				CurrentValue:     rawIndices[i].CurrentValue,
				GeneratedTime:    rawIndices[i].GeneratedTime,
			}, nil
		}
	}

	return nil, NewNotFoundError("NEPSE Index")
}

// SubIndices returns other main indices (Sensitive, Float, Sensitive Float)
// excluding the main NEPSE index.
// Note: Sector sub-indices are only available through graph endpoints.
func (c *Client) SubIndices(ctx context.Context) ([]SubIndex, error) {
	var rawIndices []NepseIndexRaw
	if err := c.apiRequest(ctx, c.config.Endpoints.NepseIndex, &rawIndices); err != nil {
		return nil, err
	}

	// Only exclude the main NEPSE index, include the other 3 main indices
	subIndices := make([]SubIndex, 0, len(rawIndices))
	for i := range rawIndices {
		if rawIndices[i].ID != nepseIndexID {
			subIndices = append(subIndices, SubIndex(rawIndices[i]))
		}
	}

	return subIndices, nil
}

// LiveMarket returns real-time price and volume data for all actively traded securities.
func (c *Client) LiveMarket(ctx context.Context) ([]LiveMarketEntry, error) {
	var liveMarket []LiveMarketEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.LiveMarket, &liveMarket); err != nil {
		return nil, err
	}
	return liveMarket, nil
}

// SupplyDemandData represents the combined supply and demand response.
type SupplyDemandData struct {
	SupplyList []SupplyDemandItem `json:"supplyList"`
	DemandList []SupplyDemandItem `json:"demandList"`
}

// SupplyDemandItem represents a single item in supply or demand list.
type SupplyDemandItem struct {
	SecurityID    int32  `json:"securityId"`
	Symbol        string `json:"symbol"`
	SecurityName  string `json:"securityName"`
	TotalQuantity int64  `json:"totalQuantity"`
	TotalOrder    int32  `json:"totalOrder"`
}

// SupplyDemand returns aggregate supply and demand data.
func (c *Client) SupplyDemand(ctx context.Context) (*SupplyDemandData, error) {
	var data SupplyDemandData
	if err := c.apiRequest(ctx, c.config.Endpoints.SupplyDemand, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// TopGainers returns securities with the highest percentage gains for the trading day.
func (c *Client) TopGainers(ctx context.Context) ([]TopGainerLoserEntry, error) {
	var topGainers []TopGainerLoserEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopGainers, &topGainers); err != nil {
		return nil, err
	}
	return topGainers, nil
}

// TopLosers returns securities with the highest percentage losses for the trading day.
func (c *Client) TopLosers(ctx context.Context) ([]TopGainerLoserEntry, error) {
	var topLosers []TopGainerLoserEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopLosers, &topLosers); err != nil {
		return nil, err
	}
	return topLosers, nil
}

// TopTenTrade returns the ten securities with the highest traded share volume.
func (c *Client) TopTenTrade(ctx context.Context) ([]TopTradeEntry, error) {
	var topTrade []TopTradeEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTrade, &topTrade); err != nil {
		return nil, err
	}
	return topTrade, nil
}

// TopTenTransaction returns the ten securities with the most transactions.
func (c *Client) TopTenTransaction(ctx context.Context) ([]TopTransactionEntry, error) {
	var topTransaction []TopTransactionEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTransaction, &topTransaction); err != nil {
		return nil, err
	}
	return topTransaction, nil
}

// TopTenTurnover returns the ten securities with the highest trading turnover (value).
func (c *Client) TopTenTurnover(ctx context.Context) ([]TopTurnoverEntry, error) {
	var topTurnover []TopTurnoverEntry
	if err := c.apiRequest(ctx, c.config.Endpoints.TopTurnover, &topTurnover); err != nil {
		return nil, err
	}
	return topTurnover, nil
}

// TodaysPrices returns price data for all securities on a given business date.
// If businessDate is empty, returns data for the current trading day.
//
// Note: This endpoint may return empty results. NEPSE's web interface uses a POST request
// that requires additional authentication not currently supported by this library.
// For current prices, consider using [Client.TopGainers], [Client.TopLosers], or
// [Client.Company] which return LTP (last traded price) data.
func (c *Client) TodaysPrices(ctx context.Context, businessDate string) ([]TodayPrice, error) {
	endpoint := c.config.Endpoints.TodaysPrice
	if businessDate != "" {
		params := url.Values{}
		params.Set("businessDate", businessDate)
		params.Set("size", "500")
		endpoint += "?" + params.Encode()
	}

	var todayPrices []TodayPrice
	if err := c.apiRequest(ctx, endpoint, &todayPrices); err != nil {
		return nil, err
	}
	return todayPrices, nil
}

// PriceHistory returns historical OHLCV data for a security within a date range.
func (c *Client) PriceHistory(ctx context.Context, securityID int32, startDate, endDate string) ([]PriceHistory, error) {
	params := url.Values{}
	params.Set("size", "500")
	params.Set("startDate", startDate)
	params.Set("endDate", endDate)
	endpoint := fmt.Sprintf("%s/%d?%s", c.config.Endpoints.CompanyPriceHistory, securityID, params.Encode())

	var response struct {
		Content []PriceHistory `json:"content"`
	}

	if err := c.apiRequest(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	return response.Content, nil
}

// PriceHistoryBySymbol returns historical OHLCV data for a security by symbol.
func (c *Client) PriceHistoryBySymbol(ctx context.Context, symbol string, startDate, endDate string) ([]PriceHistory, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.PriceHistory(ctx, security.ID, startDate, endDate)
}

// MarketDepth returns the order book (bid/ask levels) for a security.
func (c *Client) MarketDepth(ctx context.Context, securityID int32) (*MarketDepth, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.MarketDepth, securityID)

	var raw MarketDepthRaw
	if err := c.apiRequest(ctx, endpoint, &raw); err != nil {
		return nil, err
	}

	return &MarketDepth{
		TotalBuyQty:  raw.TotalBuyQty,
		TotalSellQty: raw.TotalSellQty,
		BuyDepth:     raw.MarketDepth.BuyList,
		SellDepth:    raw.MarketDepth.SellList,
	}, nil
}

// MarketDepthBySymbol returns the order book for a security by ticker symbol.
func (c *Client) MarketDepthBySymbol(ctx context.Context, symbol string) (*MarketDepth, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.MarketDepth(ctx, security.ID)
}

// Securities returns all tradable securities on the exchange.
func (c *Client) Securities(ctx context.Context) ([]Security, error) {
	var securities []Security
	if err := c.apiRequest(ctx, c.config.Endpoints.SecurityList, &securities); err != nil {
		return nil, err
	}
	return securities, nil
}

// Companies returns all listed companies on the exchange.
func (c *Client) Companies(ctx context.Context) ([]Company, error) {
	var companies []Company
	if err := c.apiRequest(ctx, c.config.Endpoints.CompanyList, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

// Company returns comprehensive information including price data for a security.
func (c *Client) Company(ctx context.Context, securityID int32) (*CompanyDetails, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDetails, securityID)

	var rawDetails CompanyDetailsRaw
	if err := c.apiRequest(ctx, endpoint, &rawDetails); err != nil {
		return nil, err
	}

	details := &CompanyDetails{
		ID:               rawDetails.SecurityData.ID,
		Symbol:           rawDetails.SecurityData.Symbol,
		SecurityName:     rawDetails.SecurityData.SecurityName,
		SectorName:       rawDetails.SecurityData.Sector,
		Email:            rawDetails.SecurityData.Email,
		ActiveStatus:     rawDetails.SecurityData.ActiveStatus,
		PermittedToTrade: rawDetails.SecurityData.PermittedToTrade,

		OpenPrice:           rawDetails.SecurityMcsData.OpenPrice,
		HighPrice:           rawDetails.SecurityMcsData.HighPrice,
		LowPrice:            rawDetails.SecurityMcsData.LowPrice,
		ClosePrice:          rawDetails.SecurityMcsData.ClosePrice,
		LastTradedPrice:     rawDetails.SecurityMcsData.LastTradedPrice,
		PreviousClose:       rawDetails.SecurityMcsData.PreviousClose,
		TotalTradeQuantity:  rawDetails.SecurityMcsData.TotalTradeQuantity,
		TotalTrades:         rawDetails.SecurityMcsData.TotalTrades,
		FiftyTwoWeekHigh:    rawDetails.SecurityMcsData.FiftyTwoWeekHigh,
		FiftyTwoWeekLow:     rawDetails.SecurityMcsData.FiftyTwoWeekLow,
		BusinessDate:        rawDetails.SecurityMcsData.BusinessDate,
		LastUpdatedDateTime: rawDetails.SecurityMcsData.LastUpdatedDateTime,
	}

	return details, nil
}

// CompanyBySymbol returns comprehensive information for a security by ticker symbol.
func (c *Client) CompanyBySymbol(ctx context.Context, symbol string) (*CompanyDetails, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.Company(ctx, security.ID)
}

// SecurityDetail returns comprehensive security information including shareholding data.
// This uses a POST request to fetch additional data not available via [Client.Company].
func (c *Client) SecurityDetail(ctx context.Context, securityID int32) (*SecurityDetail, error) {
	payloadID, err := c.computeScripGraphPayloadID(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDetails, securityID)

	var raw SecurityDetailRaw
	if err := c.apiPostRequest(ctx, endpoint, graphPostPayload{ID: payloadID}, &raw); err != nil {
		return nil, err
	}

	return &SecurityDetail{
		ID:               raw.Security.ID,
		Symbol:           raw.Security.Symbol,
		ISIN:             raw.Security.Isin,
		PermittedToTrade: raw.Security.PermittedToTrade,
		FaceValue:        raw.Security.FaceValue,

		ListedShares:    int64(raw.StockListedShares),
		PaidUpCapital:   raw.PaidUpCapital,
		IssuedCapital:   raw.IssuedCapital,
		MarketCap:       raw.MarketCapitalization,
		PublicShares:    raw.PublicShares,
		PublicPercent:   raw.PublicPercentage,
		PromoterShares:  int64(raw.PromoterShares),
		PromoterPercent: raw.PromoterPercentage,

		OpenPrice:           raw.SecurityDailyTradeDTO.OpenPrice,
		HighPrice:           raw.SecurityDailyTradeDTO.HighPrice,
		LowPrice:            raw.SecurityDailyTradeDTO.LowPrice,
		ClosePrice:          raw.SecurityDailyTradeDTO.ClosePrice,
		LastTradedPrice:     raw.SecurityDailyTradeDTO.LastTradedPrice,
		PreviousClose:       raw.SecurityDailyTradeDTO.PreviousClose,
		TotalTradedQuantity: raw.SecurityDailyTradeDTO.TotalTradeQuantity,
		TotalTrades:         raw.SecurityDailyTradeDTO.TotalTrades,
		FiftyTwoWeekHigh:    raw.SecurityDailyTradeDTO.FiftyTwoWeekHigh,
		FiftyTwoWeekLow:     raw.SecurityDailyTradeDTO.FiftyTwoWeekLow,
		BusinessDate:        raw.SecurityDailyTradeDTO.BusinessDate,
		LastUpdatedDateTime: raw.SecurityDailyTradeDTO.LastUpdatedDateTime,
	}, nil
}

// SecurityDetailBySymbol returns comprehensive security information by ticker symbol.
func (c *Client) SecurityDetailBySymbol(ctx context.Context, symbol string) (*SecurityDetail, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.SecurityDetail(ctx, security.ID)
}

// DebugSecurityDetailRaw returns the raw JSON response from the security detail endpoint.
// This is useful for debugging the API response structure.
func (c *Client) DebugSecurityDetailRaw(ctx context.Context, securityID int32) ([]byte, error) {
	payloadID, err := c.computeScripGraphPayloadID(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDetails, securityID)
	return c.apiPostRequestRaw(ctx, endpoint, graphPostPayload{ID: payloadID})
}

// SectorScrips returns a map of sector names to their constituent security symbols.
func (c *Client) SectorScrips(ctx context.Context) (SectorScrips, error) {
	// Use company list which includes sector information
	companies, err := c.Companies(ctx)
	if err != nil {
		return nil, err
	}

	sectorScrips := make(SectorScrips)

	for _, company := range companies {
		sectorName := company.SectorName
		if strings.HasSuffix(company.Symbol, "P") {
			sectorName = "Promoter Share"
		} else if sectorName == "" {
			sectorName = "Others"
		}

		sectorScrips[sectorName] = append(sectorScrips[sectorName], company.Symbol)
	}

	return sectorScrips, nil
}

// FindSecurity returns the security with the given ID.
func (c *Client) FindSecurity(ctx context.Context, securityID int32) (*Security, error) {
	return c.findSecurityByID(ctx, securityID)
}

// FindSecurityBySymbol returns the security with the given ticker symbol.
func (c *Client) FindSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	return c.findSecurityBySymbol(ctx, symbol)
}

func (c *Client) findSecurityByID(ctx context.Context, id int32) (*Security, error) {
	if id <= 0 {
		return nil, NewInvalidClientRequestError("security ID must be positive")
	}

	securities, err := c.Securities(ctx)
	if err != nil {
		return nil, err
	}

	for i := range securities {
		if securities[i].ID == id {
			return &securities[i], nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("security with ID %d", id))
}

func (c *Client) findSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, NewInvalidClientRequestError("symbol cannot be empty")
	}

	securities, err := c.Securities(ctx)
	if err != nil {
		return nil, err
	}

	for i := range securities {
		if securities[i].Symbol == symbol {
			return &securities[i], nil
		}
	}

	return nil, NewNotFoundError("security with symbol " + symbol)
}

// FloorSheet returns all trades executed on the exchange for the current trading day.
// Handles both array and paginated response formats.
// Note: Returns empty slice if no trades have occurred yet.
func (c *Client) FloorSheet(ctx context.Context) ([]FloorSheetEntry, error) {
	params := url.Values{}
	params.Set("size", "500")
	params.Set("sort", "contractId,desc")
	endpoint := c.config.Endpoints.FloorSheet + "?" + params.Encode()

	data, err := c.apiRequestRaw(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	// Try direct array format (may be empty during market hours before trades occur).
	var floorSheetArray []FloorSheetEntry
	if err := json.Unmarshal(data, &floorSheetArray); err == nil {
		return floorSheetArray, nil
	}

	// Try paginated format.
	var firstPage FloorSheetResponse
	if err := json.Unmarshal(data, &firstPage); err != nil {
		return nil, NewInvalidServerResponseError("unrecognized floor sheet response format")
	}

	all := firstPage.FloorSheets.Content
	total := firstPage.FloorSheets.TotalPages
	for p := int32(1); p < total; p++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, p)
		var page FloorSheetResponse
		if err := c.apiRequest(ctx, pageEndpoint, &page); err != nil {
			return nil, err
		}
		all = append(all, page.FloorSheets.Content...)
	}
	return all, nil
}

// FloorSheetOf returns all trades for a specific security on a given business date.
//
// IMPORTANT: As of December 2025, NEPSE has blocked this endpoint at the server level.
// All requests return 403 Forbidden. Use [Client.FloorSheet] instead for general floorsheet data.
func (c *Client) FloorSheetOf(ctx context.Context, securityID int32, businessDate string) ([]FloorSheetEntry, error) {
	params := url.Values{}
	params.Set("businessDate", businessDate)
	params.Set("size", "500")
	params.Set("sort", "contractid,desc")
	endpoint := fmt.Sprintf("%s/%d?%s", c.config.Endpoints.CompanyFloorsheet, securityID, params.Encode())

	var firstPage FloorSheetResponse
	if err := c.apiRequest(ctx, endpoint, &firstPage); err != nil {
		return nil, err
	}

	if len(firstPage.FloorSheets.Content) == 0 {
		return []FloorSheetEntry{}, nil
	}

	allEntries := firstPage.FloorSheets.Content
	totalPages := firstPage.FloorSheets.TotalPages

	for page := int32(1); page < totalPages; page++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, page)

		var pageResponse FloorSheetResponse
		if err := c.apiRequest(ctx, pageEndpoint, &pageResponse); err != nil {
			return nil, err
		}

		allEntries = append(allEntries, pageResponse.FloorSheets.Content...)
	}

	return allEntries, nil
}

// FloorSheetBySymbol returns all trades for a specific security by symbol on a given date.
//
// HACK: As of December 2025, NEPSE has blocked this endpoint at the server level.
// All requests return 403 Forbidden. Use [Client.FloorSheet] instead for general floorsheet data.
func (c *Client) FloorSheetBySymbol(ctx context.Context, symbol string, businessDate string) ([]FloorSheetEntry, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.FloorSheetOf(ctx, security.ID, businessDate)
}
