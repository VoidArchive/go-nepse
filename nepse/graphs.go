package nepse

import (
	"context"
	"fmt"
)

// Graph Data GET API Methods (aligned with Client interface)

// GetDailyNepseIndexGraph Index Graph Methods
func (h *HTTPClient) GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["nepse_index_daily_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily NEPSE index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["sensitive_index_daily_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily sensitive index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["float_index_daily_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily float index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["sensitive_float_index_daily_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily sensitive float index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyBankSubindexGraph Sector Sub-Index Graph Methods
func (h *HTTPClient) GetDailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["banking_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily banking sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["development_bank_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily development bank sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["finance_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily finance sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["hotel_tourism_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily hotel tourism sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["hydro_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily hydro sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["investment_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily investment sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["life_insurance_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily life insurance sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["manufacturing_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily manufacturing sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["microfinance_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily microfinance sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["mutual_fund_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily mutual fund sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["non_life_insurance_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily non-life insurance sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["others_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily others sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error) {
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, h.config.APIEndpoints["trading_sub_index_graph"], &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily trading sub-index graph: %w", err)
	}
	return &GraphResponse{Data: arr}, nil
}

// GetDailyScripPriceGraph Company-Specific Graph
func (h *HTTPClient) GetDailyScripPriceGraph(ctx context.Context, securityID int32) (*GraphResponse, error) {
	endpoint := fmt.Sprintf("%s%d", h.config.APIEndpoints["company_daily_graph"], securityID)
	var arr []GraphDataPoint
	if err := h.apiRequest(ctx, endpoint, &arr); err != nil {
		return nil, fmt.Errorf("failed to get daily scrip price graph for security %d: %w", securityID, err)
	}
	return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyScripPriceGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error) {
	security, err := h.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to find security %s: %w", symbol, err)
	}
	return h.GetDailyScripPriceGraph(ctx, security.ID)
}

// GetAllSubIndexGraphs helpers
func (h *HTTPClient) GetAllSubIndexGraphs(ctx context.Context) (map[string]*GraphResponse, error) {
	subIndex := map[string]func(context.Context) (*GraphResponse, error){
		"banking":            h.GetDailyBankSubindexGraph,
		"development_bank":   h.GetDailyDevelopmentBankSubindexGraph,
		"finance":            h.GetDailyFinanceSubindexGraph,
		"hotel_tourism":      h.GetDailyHotelTourismSubindexGraph,
		"hydro":              h.GetDailyHydroSubindexGraph,
		"investment":         h.GetDailyInvestmentSubindexGraph,
		"life_insurance":     h.GetDailyLifeInsuranceSubindexGraph,
		"manufacturing":      h.GetDailyManufacturingSubindexGraph,
		"microfinance":       h.GetDailyMicrofinanceSubindexGraph,
		"mutual_fund":        h.GetDailyMutualfundSubindexGraph,
		"non_life_insurance": h.GetDailyNonLifeInsuranceSubindexGraph,
		"others":             h.GetDailyOthersSubindexGraph,
		"trading":            h.GetDailyTradingSubindexGraph,
	}
	out := make(map[string]*GraphResponse)
	for name, fn := range subIndex {
		g, err := fn(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s sub-index graph: %w", name, err)
		}
		out[name] = g
	}
	return out, nil
}

func (h *HTTPClient) GetAllMainIndexGraphs(ctx context.Context) (map[string]*GraphResponse, error) {
	main := map[string]func(context.Context) (*GraphResponse, error){
		"nepse":           h.GetDailyNepseIndexGraph,
		"sensitive":       h.GetDailySensitiveIndexGraph,
		"float":           h.GetDailyFloatIndexGraph,
		"sensitive_float": h.GetDailySensitiveFloatIndexGraph,
	}
	out := make(map[string]*GraphResponse)
	for name, fn := range main {
		g, err := fn(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s index graph: %w", name, err)
		}
		out[name] = g
	}
	return out, nil
}
