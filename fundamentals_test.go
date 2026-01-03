package nepse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_CompanyProfile(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
		case "/api/nots/security/profile/123":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(CompanyProfile{
				CompanyName:          "Test Company Ltd",
				CompanyEmail:         "test@example.com",
				CompanyContactPerson: "John Doe",
				AddressField:         "Kathmandu",
				Town:                 "Kathmandu",
			})
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	profile, err := client.CompanyProfile(context.Background(), 123)
	if err != nil {
		t.Fatalf("CompanyProfile failed: %v", err)
	}

	if profile.CompanyName != "Test Company Ltd" {
		t.Errorf("expected company name 'Test Company Ltd', got '%s'", profile.CompanyName)
	}
	if profile.CompanyEmail != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", profile.CompanyEmail)
	}
}

func TestClient_BoardOfDirectors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
		case "/api/nots/security/boardOfDirectors/123":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]BoardMember{
				{
					FirstName:   "John",
					MiddleName:  "",
					LastName:    "Doe",
					Designation: "Chairman",
				},
				{
					FirstName:   "Jane",
					MiddleName:  "Mary",
					LastName:    "Smith",
					Designation: "Director",
				},
			})
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	members, err := client.BoardOfDirectors(context.Background(), 123)
	if err != nil {
		t.Fatalf("BoardOfDirectors failed: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("expected 2 members, got %d", len(members))
	}
	if members[0].FullName() != "John Doe" {
		t.Errorf("expected 'John Doe', got '%s'", members[0].FullName())
	}
	if members[1].FullName() != "Jane Mary Smith" {
		t.Errorf("expected 'Jane Mary Smith', got '%s'", members[1].FullName())
	}
}

func TestClient_CorporateActions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
		case "/api/nots/security/corporate-actions/123":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]CorporateAction{
				{
					ActiveStatus:    "BONUS_APPROVED",
					FiscalYear:      "2080-2081",
					BonusPercentage: 10.5,
				},
			})
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	actions, err := client.CorporateActions(context.Background(), 123)
	if err != nil {
		t.Fatalf("CorporateActions failed: %v", err)
	}

	if len(actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(actions))
	}
	if !actions[0].IsBonus() {
		t.Error("expected IsBonus() to return true")
	}
	if actions[0].BonusPercentage != 10.5 {
		t.Errorf("expected bonus 10.5%%, got %.2f%%", actions[0].BonusPercentage)
	}
}

func TestClient_Reports(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
		case "/api/nots/application/reports/123":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]Report{
				{
					ID:           1,
					ActiveStatus: "A",
					FiscalReport: &FiscalReport{
						ID:               1,
						PEValue:          15.5,
						EPSValue:         25.0,
						NetWorthPerShare: 120.0,
						ReportTypeMaster: &ReportTypeMaster{
							ID:         1,
							ReportName: "Annual Report",
						},
						FinancialYear: &FinancialYear{
							ID:           1,
							FYName:       "2023-2024",
							FYNameNepali: "2080-2081",
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	reports, err := client.Reports(context.Background(), 123)
	if err != nil {
		t.Fatalf("Reports failed: %v", err)
	}

	if len(reports) != 1 {
		t.Errorf("expected 1 report, got %d", len(reports))
	}
	if !reports[0].IsAnnual() {
		t.Error("expected IsAnnual() to return true")
	}
	if reports[0].FiscalReport.EPSValue != 25.0 {
		t.Errorf("expected EPS 25.0, got %.2f", reports[0].FiscalReport.EPSValue)
	}
}

func TestClient_Dividends(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/authenticate/prove":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResponse())
		case "/api/nots/application/dividend/123":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]Dividend{
				{
					ID:           1,
					ActiveStatus: "A",
					CompanyNews: &CompanyNews{
						ID:           1,
						NewsHeadline: "Dividend Declaration",
						NewsBody:     "Cash Dividend: 10%, Bonus: 5%",
						DividendsNotice: &DividendNotice{
							ID:           1,
							CashDividend: 10.0,
							BonusShare:   5.0,
							RightShare:   0,
							FinancialYear: &FinancialYear{
								ID:           1,
								FYNameNepali: "2080-2081",
							},
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(&Options{
		BaseURL:     server.URL,
		HTTPTimeout: 5 * time.Second,
		MaxRetries:  0,
		Config: &Config{
			BaseURL:   server.URL,
			Endpoints: DefaultEndpoints(),
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	dividends, err := client.Dividends(context.Background(), 123)
	if err != nil {
		t.Fatalf("Dividends failed: %v", err)
	}

	if len(dividends) != 1 {
		t.Errorf("expected 1 dividend, got %d", len(dividends))
	}
	if !dividends[0].HasCashDividend() {
		t.Error("expected HasCashDividend() to return true")
	}
	if !dividends[0].HasBonusDividend() {
		t.Error("expected HasBonusDividend() to return true")
	}
	if dividends[0].CashPercentage() != 10.0 {
		t.Errorf("expected cash dividend 10.0%%, got %.2f%%", dividends[0].CashPercentage())
	}
	if dividends[0].FiscalYear() != "2080-2081" {
		t.Errorf("expected fiscal year '2080-2081', got '%s'", dividends[0].FiscalYear())
	}
}

func TestBoardMember_FullName(t *testing.T) {
	tests := []struct {
		name     string
		member   BoardMember
		expected string
	}{
		{
			name: "with middle name",
			member: BoardMember{
				FirstName:  "John",
				MiddleName: "William",
				LastName:   "Doe",
			},
			expected: "John William Doe",
		},
		{
			name: "without middle name",
			member: BoardMember{
				FirstName:  "Jane",
				MiddleName: "",
				LastName:   "Smith",
			},
			expected: "Jane Smith",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.member.FullName(); got != tt.expected {
				t.Errorf("FullName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCorporateAction_TypeChecks(t *testing.T) {
	bonus := CorporateAction{BonusPercentage: 10.0}
	if !bonus.IsBonus() {
		t.Error("expected IsBonus() true for bonus action")
	}

	rightPct := 15.0
	right := CorporateAction{RightPercentage: &rightPct}
	if !right.IsRight() {
		t.Error("expected IsRight() true for rights action")
	}

	cashDiv := 5.0
	cash := CorporateAction{CashDividend: &cashDiv}
	if !cash.IsCashDividend() {
		t.Error("expected IsCashDividend() true for cash dividend")
	}
}

func TestReport_TypeChecks(t *testing.T) {
	annual := Report{
		FiscalReport: &FiscalReport{
			ReportTypeMaster: &ReportTypeMaster{ReportName: "Annual Report"},
		},
	}
	if !annual.IsAnnual() {
		t.Error("expected IsAnnual() true")
	}
	if annual.IsQuarterly() {
		t.Error("expected IsQuarterly() false for annual report")
	}

	quarterly := Report{
		FiscalReport: &FiscalReport{
			ReportTypeMaster: &ReportTypeMaster{ReportName: "Quarterly Report"},
			QuarterMaster:    &QuarterMaster{QuarterName: "First Quarter"},
		},
	}
	if !quarterly.IsQuarterly() {
		t.Error("expected IsQuarterly() true")
	}
	if quarterly.QuarterName() != "First Quarter" {
		t.Errorf("expected 'First Quarter', got '%s'", quarterly.QuarterName())
	}
}
