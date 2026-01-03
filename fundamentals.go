package nepse

import (
	"context"
	"fmt"
)

// CompanyProfile returns detailed profile information for a security.
func (c *Client) CompanyProfile(ctx context.Context, securityID int32) (*CompanyProfile, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CompanyProfile, securityID)

	var profile CompanyProfile
	if err := c.apiRequest(ctx, endpoint, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// CompanyProfileBySymbol returns detailed profile information for a security by symbol.
func (c *Client) CompanyProfileBySymbol(ctx context.Context, symbol string) (*CompanyProfile, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.CompanyProfile(ctx, security.ID)
}

// BoardOfDirectors returns the board of directors for a security.
func (c *Client) BoardOfDirectors(ctx context.Context, securityID int32) ([]BoardMember, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.BoardOfDirectors, securityID)

	var members []BoardMember
	if err := c.apiRequest(ctx, endpoint, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// BoardOfDirectorsBySymbol returns the board of directors for a security by symbol.
func (c *Client) BoardOfDirectorsBySymbol(ctx context.Context, symbol string) ([]BoardMember, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.BoardOfDirectors(ctx, security.ID)
}

// CorporateActions returns corporate actions (bonus, rights, dividends) for a security.
func (c *Client) CorporateActions(ctx context.Context, securityID int32) ([]CorporateAction, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.CorporateActions, securityID)

	var actions []CorporateAction
	if err := c.apiRequest(ctx, endpoint, &actions); err != nil {
		return nil, err
	}
	return actions, nil
}

// CorporateActionsBySymbol returns corporate actions for a security by symbol.
func (c *Client) CorporateActionsBySymbol(ctx context.Context, symbol string) ([]CorporateAction, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.CorporateActions(ctx, security.ID)
}

// Reports returns quarterly and annual reports for a security.
func (c *Client) Reports(ctx context.Context, securityID int32) ([]Report, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.Reports, securityID)

	var reports []Report
	if err := c.apiRequest(ctx, endpoint, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

// ReportsBySymbol returns quarterly and annual reports for a security by symbol.
func (c *Client) ReportsBySymbol(ctx context.Context, symbol string) ([]Report, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.Reports(ctx, security.ID)
}

// Dividends returns dividend history for a security.
func (c *Client) Dividends(ctx context.Context, securityID int32) ([]Dividend, error) {
	endpoint := fmt.Sprintf("%s/%d", c.config.Endpoints.Dividend, securityID)

	var dividends []Dividend
	if err := c.apiRequest(ctx, endpoint, &dividends); err != nil {
		return nil, err
	}
	return dividends, nil
}

// DividendsBySymbol returns dividend history for a security by symbol.
func (c *Client) DividendsBySymbol(ctx context.Context, symbol string) ([]Dividend, error) {
	security, err := c.findSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return c.Dividends(ctx, security.ID)
}
