package nepse

import (
	"context"
	"fmt"
	"time"
)

// dummyData is a static array used by NEPSE's obfuscation algorithm
// to compute POST payload IDs for graph endpoints.
var dummyData = [100]int{
	147, 117, 239, 143, 157, 312, 161, 612, 512, 804,
	411, 527, 170, 511, 421, 667, 764, 621, 301, 106,
	133, 793, 411, 511, 312, 423, 344, 346, 653, 758,
	342, 222, 236, 811, 711, 611, 122, 447, 128, 199,
	183, 135, 489, 703, 800, 745, 152, 863, 134, 211,
	142, 564, 375, 793, 212, 153, 138, 153, 648, 611,
	151, 649, 318, 143, 117, 756, 119, 141, 717, 113,
	112, 146, 162, 660, 693, 261, 362, 354, 251, 641,
	157, 178, 631, 192, 734, 445, 192, 883, 187, 122,
	591, 731, 852, 384, 565, 596, 451, 772, 624, 691,
}

// computeBasePayloadID computes the base payload value used by graph endpoints.
// Returns: dummyData[dummyID] + dummyID + 2 * day
func (c *Client) computeBasePayloadID(ctx context.Context) (int, int, error) {
	// Get the dummy ID from market status
	status, err := c.MarketStatus(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get market status: %w", err)
	}
	dummyID := int(status.ID)

	// Ensure dummyID is within bounds
	if dummyID < 0 || dummyID >= len(dummyData) {
		dummyID = dummyID % len(dummyData)
		if dummyID < 0 {
			dummyID += len(dummyData)
		}
	}

	// Get current day of month in Nepal timezone (NEPSE expects NPT)
	loc, _ := time.LoadLocation("Asia/Kathmandu")
	day := time.Now().In(loc).Day()

	// Compute base value: dummyData[dummyID] + dummyID + 2 * day
	e := dummyData[dummyID] + dummyID + 2*day

	return e, day, nil
}

// computeIndexGraphPayloadID computes the POST payload ID for index graph endpoints.
// Uses salt values in addition to the base calculation.
func (c *Client) computeIndexGraphPayloadID(ctx context.Context) (int, error) {
	e, day, err := c.computeBasePayloadID(ctx)
	if err != nil {
		return 0, err
	}

	// Get salt values
	salts, err := c.authManager.GetSalts(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get salts: %w", err)
	}

	// Compute final payload ID using salt values
	// Logic: if (e % 10 < 5) use salts[3] * day - salts[2], else use salts[1] * day - salts[0]
	// Python uses 1-indexed array, so: salts[3] = Salt4, salts[1] = Salt2, salts[2] = Salt3, salts[0] = Salt1
	var payloadID int
	if e%10 < 5 {
		payloadID = e + salts.Salt4*day - salts.Salt3
	} else {
		payloadID = e + salts.Salt2*day - salts.Salt1
	}

	return payloadID, nil
}

// computeScripGraphPayloadID computes the POST payload ID for security/scrip graph endpoints.
// Uses only the base calculation without salt adjustment.
func (c *Client) computeScripGraphPayloadID(ctx context.Context) (int, error) {
	e, _, err := c.computeBasePayloadID(ctx)
	if err != nil {
		return 0, err
	}
	return e, nil
}

// graphPostPayload is the request body for graph POST endpoints.
type graphPostPayload struct {
	ID int `json:"id"`
}

// IndexType represents the type of market index for graph data retrieval.
type IndexType int

const (
	IndexNepse IndexType = iota
	IndexSensitive
	IndexFloat
	IndexSensitiveFloat
	IndexBanking
	IndexDevBank
	IndexFinance
	IndexHotelTourism
	IndexHydro
	IndexInvestment
	IndexLifeInsurance
	IndexManufacturing
	IndexMicrofinance
	IndexMutualFund
	IndexNonLifeInsurance
	IndexOthers
	IndexTrading
)

// indexEndpoint returns the endpoint for a given index type.
func (c *Client) indexEndpoint(indexType IndexType) string {
	endpoints := c.config.Endpoints
	switch indexType {
	case IndexNepse:
		return endpoints.GraphNepseIndex
	case IndexSensitive:
		return endpoints.GraphSensitiveIndex
	case IndexFloat:
		return endpoints.GraphFloatIndex
	case IndexSensitiveFloat:
		return endpoints.GraphSensitiveFloatIndex
	case IndexBanking:
		return endpoints.GraphBankingSubindex
	case IndexDevBank:
		return endpoints.GraphDevBankSubindex
	case IndexFinance:
		return endpoints.GraphFinanceSubindex
	case IndexHotelTourism:
		return endpoints.GraphHotelSubindex
	case IndexHydro:
		return endpoints.GraphHydroSubindex
	case IndexInvestment:
		return endpoints.GraphInvestmentSubindex
	case IndexLifeInsurance:
		return endpoints.GraphLifeInsSubindex
	case IndexManufacturing:
		return endpoints.GraphManufacturingSubindex
	case IndexMicrofinance:
		return endpoints.GraphMicrofinanceSubindex
	case IndexMutualFund:
		return endpoints.GraphMutualFundSubindex
	case IndexNonLifeInsurance:
		return endpoints.GraphNonLifeInsSubindex
	case IndexOthers:
		return endpoints.GraphOthersSubindex
	case IndexTrading:
		return endpoints.GraphTradingSubindex
	default:
		return endpoints.GraphNepseIndex
	}
}

// DailyIndexGraph returns intraday graph data points for any market index.
func (c *Client) DailyIndexGraph(ctx context.Context, indexType IndexType) (*GraphResponse, error) {
	payloadID, err := c.computeIndexGraphPayloadID(ctx)
	if err != nil {
		return nil, err
	}

	var arr []GraphDataPoint
	if err := c.apiPostRequest(ctx, c.indexEndpoint(indexType), graphPostPayload{ID: payloadID}, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// DailyNepseIndexGraph returns intraday graph data for the main NEPSE index.
func (c *Client) DailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexNepse)
}

// DailySensitiveIndexGraph returns intraday graph data for the sensitive index.
func (c *Client) DailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexSensitive)
}

// DailyFloatIndexGraph returns intraday graph data for the float index.
func (c *Client) DailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexFloat)
}

// DailySensitiveFloatIndexGraph returns intraday graph data for the sensitive float index.
func (c *Client) DailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexSensitiveFloat)
}

// DailyBankSubindexGraph returns intraday graph data for the banking sector sub-index.
func (c *Client) DailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexBanking)
}

// DailyDevelopmentBankSubindexGraph returns intraday graph data for the development bank sector.
func (c *Client) DailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexDevBank)
}

// DailyFinanceSubindexGraph returns intraday graph data for the finance sector.
func (c *Client) DailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexFinance)
}

// DailyHotelTourismSubindexGraph returns intraday graph data for the hotel & tourism sector.
func (c *Client) DailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexHotelTourism)
}

// DailyHydroSubindexGraph returns intraday graph data for the hydropower sector.
func (c *Client) DailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexHydro)
}

// DailyInvestmentSubindexGraph returns intraday graph data for the investment sector.
func (c *Client) DailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexInvestment)
}

// DailyLifeInsuranceSubindexGraph returns intraday graph data for the life insurance sector.
func (c *Client) DailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexLifeInsurance)
}

// DailyManufacturingSubindexGraph returns intraday graph data for the manufacturing sector.
func (c *Client) DailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexManufacturing)
}

// DailyMicrofinanceSubindexGraph returns intraday graph data for the microfinance sector.
func (c *Client) DailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexMicrofinance)
}

// DailyMutualfundSubindexGraph returns intraday graph data for the mutual fund sector.
func (c *Client) DailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexMutualFund)
}

// DailyNonLifeInsuranceSubindexGraph returns intraday graph data for the non-life insurance sector.
func (c *Client) DailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexNonLifeInsurance)
}

// DailyOthersSubindexGraph returns intraday graph data for the others sector.
func (c *Client) DailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexOthers)
}

// DailyTradingSubindexGraph returns intraday graph data for the trading sector.
func (c *Client) DailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	return c.DailyIndexGraph(ctx, IndexTrading)
}

// DailyScripGraph returns intraday price graph data for a specific security.
func (c *Client) DailyScripGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	payloadID, err := c.computeScripGraphPayloadID(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyDailyGraph, securityID)
	var arr []GraphDataPoint
	if err := c.apiPostRequest(ctx, endpoint, graphPostPayload{ID: payloadID}, &arr); err != nil {
		return nil, err
	}
	return &GraphResponse{Data: arr}, nil
}

// DailyScripGraphBySymbol returns intraday price graph data for a security by ticker symbol.
func (c *Client) DailyScripGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.DailyScripGraph(ctx, security.ID)
}
