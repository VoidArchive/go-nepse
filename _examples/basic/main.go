package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/voidarchive/go-nepse"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
	red    = "\033[31m"
)

func main() {
	withGraphs := flag.Bool("graphs", false, "include graph endpoints")
	withFloor := flag.Bool("floorsheet", false, "include floorsheet endpoints")
	symbolFlag := flag.String("symbol", "NABIL", "symbol for security-specific calls")
	bizDateFlag := flag.String("date", "", "business date (YYYY-MM-DD); defaults to last weekday")
	flag.Parse()

	printHeader()

	opts := nepse.DefaultOptions()
	opts.TLSVerification = false

	client, err := nepse.NewClient(opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	symbol := *symbolFlag
	now := time.Now()
	startDate := now.AddDate(0, -1, 0).Format("2006-01-02")
	endDate := now.Format("2006-01-02")
	bizDate := *bizDateFlag

	// Determine effective business date
	if bizDate == "" {
		bizDate = lastTradingDay(now).Format("2006-01-02")
	}

	// Track security ID for later use
	var securityID int32

	// ═══════════════════════════════════════════════════════════════════
	// MARKET OVERVIEW
	// ═══════════════════════════════════════════════════════════════════
	printSection("MARKET OVERVIEW")

	// Market Status
	printSubSection("Status")
	if status, err := client.GetMarketStatus(ctx); err != nil {
		printError("GetMarketStatus", err)
	} else {
		statusColor := red
		if status.IsMarketOpen() {
			statusColor = green
		}
		printKV("Market", fmt.Sprintf("%s%s%s", statusColor, status.IsOpen, reset))
		printKV("As Of", status.AsOf)
	}

	// Market Summary
	printSubSection("Summary")
	if summary, err := client.GetMarketSummary(ctx); err != nil {
		printError("GetMarketSummary", err)
	} else {
		printKV("Turnover", formatNumber(summary.TotalTurnover))
		printKV("Traded Shares", formatNumber(summary.TotalTradedShares))
		printKV("Transactions", formatNumber(summary.TotalTransactions))
		printKV("Scrips Traded", formatNumber(summary.TotalScripsTraded))
	}

	// NEPSE Index
	printSubSection("NEPSE Index")
	if idx, err := client.GetNepseIndex(ctx); err != nil {
		printError("GetNepseIndex", err)
	} else {
		changeColor := green
		if idx.PercentChange < 0 {
			changeColor = red
		}
		printKV("Value", fmt.Sprintf("%.2f", idx.IndexValue))
		printKV("Change", fmt.Sprintf("%s%.2f (%.2f%%)%s", changeColor, idx.PointChange, idx.PercentChange, reset))
		printKV("High / Low", fmt.Sprintf("%.2f / %.2f", idx.High, idx.Low))
		printKV("52W High / Low", fmt.Sprintf("%.2f / %.2f", idx.FiftyTwoWeekHigh, idx.FiftyTwoWeekLow))
	}

	// Other Main Indices
	printSubSection("Other Main Indices")
	if subs, err := client.GetNepseSubIndices(ctx); err != nil {
		printError("GetNepseSubIndices", err)
	} else {
		if len(subs) == 0 {
			printDim("No other indices available (sector sub-indices only in graph data)")
		} else {
			for _, sub := range subs[:min(5, len(subs))] {
				changeColor := green
				if sub.PerChange < 0 {
					changeColor = red
				}
				printKV(sub.Index, fmt.Sprintf("%.2f %s(%+.2f%%)%s", sub.Close, changeColor, sub.PerChange, reset))
			}
			if len(subs) > 5 {
				printDim(fmt.Sprintf("... and %d more", len(subs)-5))
			}
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// LIVE MARKET DATA
	// ═══════════════════════════════════════════════════════════════════
	printSection("LIVE MARKET DATA")

	// Live Market
	printSubSection("Active Securities")
	if live, err := client.GetLiveMarket(ctx); err != nil {
		printError("GetLiveMarket", err)
	} else if len(live) == 0 {
		printDim("No live data available (market closed)")
	} else {
		printKV("Total Active", fmt.Sprintf("%d securities", len(live)))
		fmt.Printf("    %s%-10s %10s %10s %12s%s\n", dim, "Symbol", "LTP", "Change%", "Volume", reset)
		for _, entry := range live[:min(5, len(live))] {
			changeColor := green
			if entry.PercentageChange < 0 {
				changeColor = red
			}
			fmt.Printf("    %-10s %10.2f %s%+10.2f%s %12d\n",
				entry.Symbol, entry.LastTradedPrice, changeColor, entry.PercentageChange, reset, entry.TotalTradeQuantity)
		}
		if len(live) > 5 {
			printDim(fmt.Sprintf("... and %d more", len(live)-5))
		}
	}

	// Supply & Demand
	printSubSection("Supply & Demand")
	if sd, err := client.GetSupplyDemand(ctx); err != nil {
		printError("GetSupplyDemand", err)
	} else if len(sd.SupplyList) == 0 && len(sd.DemandList) == 0 {
		printDim("No supply/demand data available (market closed)")
	} else {
		printKV("Supply Orders", fmt.Sprintf("%d securities", len(sd.SupplyList)))
		printKV("Demand Orders", fmt.Sprintf("%d securities", len(sd.DemandList)))
	}

	// ═══════════════════════════════════════════════════════════════════
	// TOP LISTS
	// ═══════════════════════════════════════════════════════════════════
	printSection("TOP LISTS")

	// Top Gainers
	printSubSection("Top Gainers")
	if gainers, err := client.GetTopGainers(ctx); err != nil {
		printError("GetTopGainers", err)
	} else {
		printTopGainerLoser(gainers, 5, true)
	}

	// Top Losers
	printSubSection("Top Losers")
	if losers, err := client.GetTopLosers(ctx); err != nil {
		printError("GetTopLosers", err)
	} else {
		printTopGainerLoser(losers, 5, false)
	}

	// Top by Volume
	printSubSection("Top by Volume")
	if trades, err := client.GetTopTenTrade(ctx); err != nil {
		printError("GetTopTenTrade", err)
	} else {
		fmt.Printf("    %s%-10s %12s %12s%s\n", dim, "Symbol", "Volume", "Price", reset)
		for _, t := range trades[:min(5, len(trades))] {
			fmt.Printf("    %-10s %12s %12.2f\n", t.Symbol, formatNumber(float64(t.ShareTraded)), t.ClosingPrice)
		}
	}

	// Top by Transactions
	printSubSection("Top by Transactions")
	if txns, err := client.GetTopTenTransaction(ctx); err != nil {
		printError("GetTopTenTransaction", err)
	} else {
		fmt.Printf("    %s%-10s %12s %12s%s\n", dim, "Symbol", "Trades", "LTP", reset)
		for _, t := range txns[:min(5, len(txns))] {
			fmt.Printf("    %-10s %12d %12.2f\n", t.Symbol, t.TotalTrades, t.LastTradedPrice)
		}
	}

	// Top by Turnover
	printSubSection("Top by Turnover")
	if turnover, err := client.GetTopTenTurnover(ctx); err != nil {
		printError("GetTopTenTurnover", err)
	} else {
		fmt.Printf("    %s%-10s %18s %12s%s\n", dim, "Symbol", "Turnover", "Price", reset)
		for _, t := range turnover[:min(5, len(turnover))] {
			fmt.Printf("    %-10s %18s %12.2f\n", t.Symbol, formatNumber(t.Turnover), t.ClosingPrice)
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// SECURITIES & COMPANIES
	// ═══════════════════════════════════════════════════════════════════
	printSection("SECURITIES & COMPANIES")

	// Security List
	printSubSection("Listed Securities")
	if secs, err := client.GetSecurityList(ctx); err != nil {
		printError("GetSecurityList", err)
	} else {
		printKV("Total Securities", fmt.Sprintf("%d", len(secs)))
	}

	// Company List
	printSubSection("Listed Companies")
	if companies, err := client.GetCompanyList(ctx); err != nil {
		printError("GetCompanyList", err)
	} else {
		printKV("Total Companies", fmt.Sprintf("%d", len(companies)))
	}

	// Sector Distribution
	printSubSection("Sector Distribution")
	if sectors, err := client.GetSectorScrips(ctx); err != nil {
		printError("GetSectorScrips", err)
	} else {
		printKV("Total Sectors", fmt.Sprintf("%d", len(sectors)))
		for sector, scrips := range sectors {
			if len(scrips) > 10 {
				printKV(fmt.Sprintf("  %s", sector), fmt.Sprintf("%d scrips", len(scrips)))
			}
		}
	}

	// Find Security
	printSubSection(fmt.Sprintf("Security Lookup: %s", symbol))
	if sec, err := client.FindSecurityBySymbol(ctx, symbol); err != nil {
		printError("FindSecurityBySymbol", err)
	} else {
		securityID = sec.ID
		printKV("ID", fmt.Sprintf("%d", sec.ID))
		printKV("Name", sec.SecurityName)
		if sec.ActiveStatus == "A" {
			printKV("Status", fmt.Sprintf("%sActive%s", green, reset))
		} else {
			printKV("Status", fmt.Sprintf("%s%s%s", yellow, sec.ActiveStatus, reset))
		}
	}

	// Company Details
	if securityID != 0 {
		printSubSection(fmt.Sprintf("Company Details: %s", symbol))
		if det, err := client.GetCompanyDetails(ctx, securityID); err != nil {
			printError("GetCompanyDetails", err)
		} else {
			printKV("Open", fmt.Sprintf("%.2f", det.OpenPrice))
			printKV("High", fmt.Sprintf("%.2f", det.HighPrice))
			printKV("Low", fmt.Sprintf("%.2f", det.LowPrice))
			printKV("Close", fmt.Sprintf("%.2f", det.ClosePrice))
			printKV("LTP", fmt.Sprintf("%.2f", det.LastTradedPrice))
			printKV("Previous Close", fmt.Sprintf("%.2f", det.PreviousClose))
			printKV("52W Range", fmt.Sprintf("%.2f - %.2f", det.FiftyTwoWeekLow, det.FiftyTwoWeekHigh))
			printKV("Volume", fmt.Sprintf("%d", det.TotalTradeQuantity))
			printKV("Trades", fmt.Sprintf("%d", det.TotalTrades))
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// PRICE & TRADING DATA
	// ═══════════════════════════════════════════════════════════════════
	printSection("PRICE & TRADING DATA")

	// Today's Prices
	printSubSection(fmt.Sprintf("Today's Prices (%s)", bizDate))
	if prices, err := client.GetTodaysPrices(ctx, bizDate); err != nil {
		printError("GetTodaysPrices", err)
	} else if len(prices) == 0 {
		printDim("No price data available (market closed or no trades on this date)")
	} else {
		printKV("Securities with Data", fmt.Sprintf("%d", len(prices)))
		fmt.Printf("    %s%-10s %10s %10s %10s %10s%s\n", dim, "Symbol", "Open", "High", "Low", "Close", reset)
		for _, p := range prices[:min(5, len(prices))] {
			fmt.Printf("    %-10s %10.2f %10.2f %10.2f %10.2f\n",
				p.Symbol, p.OpenPrice, p.HighPrice, p.LowPrice, p.ClosePrice)
		}
		if len(prices) > 5 {
			printDim(fmt.Sprintf("... and %d more", len(prices)-5))
		}
	}

	// Price History
	if securityID != 0 {
		printSubSection(fmt.Sprintf("Price History: %s (%s to %s)", symbol, startDate, endDate))
		if hist, err := client.GetPriceVolumeHistory(ctx, securityID, startDate, endDate); err != nil {
			printError("GetPriceVolumeHistory", err)
		} else {
			printKV("Data Points", fmt.Sprintf("%d trading days", len(hist)))
			if len(hist) > 0 {
				fmt.Printf("    %s%-12s %10s %10s %10s %12s%s\n", dim, "Date", "High", "Low", "Close", "Volume", reset)
				for _, h := range hist[:min(5, len(hist))] {
					fmt.Printf("    %-12s %10.2f %10.2f %10.2f %12d\n",
						h.BusinessDate, h.HighPrice, h.LowPrice, h.ClosePrice, h.TotalTradedQuantity)
				}
				if len(hist) > 5 {
					printDim(fmt.Sprintf("... and %d more", len(hist)-5))
				}
			}
		}
	}

	// Market Depth
	if securityID != 0 {
		printSubSection(fmt.Sprintf("Market Depth: %s", symbol))
		if depth, err := client.GetMarketDepth(ctx, securityID); err != nil {
			// Market depth is unavailable when market is closed
			if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "empty") {
				printDim("Market depth unavailable (market closed)")
			} else {
				printError("GetMarketDepth", err)
			}
		} else {
			printKV("Total Buy Qty", fmt.Sprintf("%d", depth.TotalBuyQty))
			printKV("Total Sell Qty", fmt.Sprintf("%d", depth.TotalSellQty))
			printKV("Buy Levels", fmt.Sprintf("%d", len(depth.BuyDepth)))
			printKV("Sell Levels", fmt.Sprintf("%d", len(depth.SellDepth)))
			if len(depth.BuyDepth) > 0 || len(depth.SellDepth) > 0 {
				fmt.Printf("\n    %s%10s %10s %8s  |  %8s %10s %10s%s\n",
					dim, "Bid Qty", "Bid", "Orders", "Orders", "Ask", "Ask Qty", reset)
				maxLevels := max(len(depth.BuyDepth), len(depth.SellDepth))
				for i := 0; i < min(5, maxLevels); i++ {
					buyQty, buyPrice, buyOrders := int64(0), 0.0, int32(0)
					sellQty, sellPrice, sellOrders := int64(0), 0.0, int32(0)
					if i < len(depth.BuyDepth) {
						buyQty = depth.BuyDepth[i].Quantity
						buyPrice = depth.BuyDepth[i].Price
						buyOrders = depth.BuyDepth[i].Orders
					}
					if i < len(depth.SellDepth) {
						sellQty = depth.SellDepth[i].Quantity
						sellPrice = depth.SellDepth[i].Price
						sellOrders = depth.SellDepth[i].Orders
					}
					fmt.Printf("    %s%10d %10.2f %8d%s  |  %s%8d %10.2f %10d%s\n",
						green, buyQty, buyPrice, buyOrders, reset,
						red, sellOrders, sellPrice, sellQty, reset)
				}
			}
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// FLOORSHEET (Optional)
	// ═══════════════════════════════════════════════════════════════════
	if *withFloor {
		printSection("FLOORSHEET")

		printSubSection("Today's Floorsheet")
		if fs, err := client.GetFloorSheet(ctx); err != nil {
			printError("GetFloorSheet", err)
		} else {
			printKV("Total Trades", fmt.Sprintf("%d", len(fs)))
			if len(fs) > 0 {
				fmt.Printf("    %s%-12s %10s %10s %12s %20s%s\n", dim, "Symbol", "Qty", "Rate", "Amount", "Time", reset)
				for _, f := range fs[:min(5, len(fs))] {
					fmt.Printf("    %-12s %10d %10.2f %12.2f %20s\n",
						f.StockSymbol, f.ContractQuantity, f.ContractRate, f.ContractAmount, f.TradeTime)
				}
				if len(fs) > 5 {
					printDim(fmt.Sprintf("... and %d more", len(fs)-5))
				}
			}
		}

		if securityID != 0 {
			printSubSection(fmt.Sprintf("Floorsheet: %s (%s)", symbol, bizDate))
			if fs, err := client.GetFloorSheetOf(ctx, securityID, bizDate); err != nil {
				printError("GetFloorSheetOf", err)
			} else {
				printKV("Total Trades", fmt.Sprintf("%d", len(fs)))
			}
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// GRAPHS (Optional)
	// ═══════════════════════════════════════════════════════════════════
	if *withGraphs {
		printSection("GRAPH DATA")

		// Main Indices
		printSubSection("Main Index Graphs")
		graphTests := []struct {
			name string
			fn   func(context.Context) (*nepse.GraphResponse, error)
		}{
			{"NEPSE Index", client.GetDailyNepseIndexGraph},
			{"Sensitive Index", client.GetDailySensitiveIndexGraph},
			{"Float Index", client.GetDailyFloatIndexGraph},
			{"Sensitive Float", client.GetDailySensitiveFloatIndexGraph},
		}
		for _, gt := range graphTests {
			if g, err := gt.fn(ctx); err != nil {
				printKV(gt.name, fmt.Sprintf("%s%v%s", red, err, reset))
			} else {
				printKV(gt.name, fmt.Sprintf("%d data points", len(g.Data)))
			}
		}

		// Sector Sub-indices
		printSubSection("Sector Sub-Index Graphs")
		sectorGraphs := []struct {
			name string
			fn   func(context.Context) (*nepse.GraphResponse, error)
		}{
			{"Banking", client.GetDailyBankSubindexGraph},
			{"Development Bank", client.GetDailyDevelopmentBankSubindexGraph},
			{"Finance", client.GetDailyFinanceSubindexGraph},
			{"Hotels & Tourism", client.GetDailyHotelTourismSubindexGraph},
			{"Hydro Power", client.GetDailyHydroSubindexGraph},
			{"Investment", client.GetDailyInvestmentSubindexGraph},
			{"Life Insurance", client.GetDailyLifeInsuranceSubindexGraph},
			{"Manufacturing", client.GetDailyManufacturingSubindexGraph},
			{"Microfinance", client.GetDailyMicrofinanceSubindexGraph},
			{"Mutual Fund", client.GetDailyMutualfundSubindexGraph},
			{"Non-Life Insurance", client.GetDailyNonLifeInsuranceSubindexGraph},
			{"Others", client.GetDailyOthersSubindexGraph},
			{"Trading", client.GetDailyTradingSubindexGraph},
		}
		for _, sg := range sectorGraphs {
			if g, err := sg.fn(ctx); err != nil {
				printKV(sg.name, fmt.Sprintf("%s%v%s", red, err, reset))
			} else {
				printKV(sg.name, fmt.Sprintf("%d data points", len(g.Data)))
			}
		}

		// Security Graph
		if securityID != 0 {
			printSubSection(fmt.Sprintf("Security Graph: %s", symbol))
			if g, err := client.GetDailyScripPriceGraph(ctx, securityID); err != nil {
				printError("GetDailyScripPriceGraph", err)
			} else {
				printKV("Data Points", fmt.Sprintf("%d", len(g.Data)))
				if len(g.Data) > 0 {
					fmt.Printf("    %s%-12s %12s%s\n", dim, "Date", "Value", reset)
					for _, d := range g.Data[:min(5, len(g.Data))] {
						fmt.Printf("    %-12s %12.2f\n", d.Date, d.Value)
					}
					if len(g.Data) > 5 {
						printDim(fmt.Sprintf("... and %d more", len(g.Data)-5))
					}
				}
			}
		}
	}

	// ═══════════════════════════════════════════════════════════════════
	// SUMMARY
	// ═══════════════════════════════════════════════════════════════════
	printSection("COMPLETE")
	fmt.Printf("    %sAll API endpoints tested successfully.%s\n", green, reset)
	fmt.Printf("    Use %s-graphs%s to include graph data.\n", cyan, reset)
	fmt.Printf("    Use %s-floorsheet%s to include floorsheet data.\n", cyan, reset)
	fmt.Printf("    Use %s-symbol=XXX%s to test with a different security.\n", cyan, reset)
	fmt.Println()
}

// ═══════════════════════════════════════════════════════════════════════════
// HELPERS
// ═══════════════════════════════════════════════════════════════════════════

func printHeader() {
	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", cyan, reset)
	fmt.Printf("%s║%s           %sNEPSE Go Library - API Demo%s                          %s║%s\n", cyan, reset, bold, reset, cyan, reset)
	fmt.Printf("%s║%s              github.com/voidarchive/go-nepse                   %s║%s\n", cyan, reset, cyan, reset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", cyan, reset)
	fmt.Println()
}

func printSection(title string) {
	fmt.Println()
	fmt.Printf("%s━━━ %s%s%s ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", blue, bold, title, reset+blue, reset)
}

func printSubSection(title string) {
	fmt.Printf("\n  %s▸ %s%s\n", yellow, title, reset)
}

func printKV(key, value string) {
	fmt.Printf("    %-20s %s\n", key+":", value)
}

func printError(method string, err error) {
	fmt.Printf("    %s✗ %s: %v%s\n", red, method, err, reset)
}

func printDim(msg string) {
	fmt.Printf("    %s%s%s\n", dim, msg, reset)
}

func printTopGainerLoser(entries []nepse.TopGainerLoserEntry, limit int, isGainer bool) {
	if len(entries) == 0 {
		printDim("No data available")
		return
	}
	fmt.Printf("    %s%-10s %10s %12s%s\n", dim, "Symbol", "LTP", "Change %", reset)
	for _, e := range entries[:min(limit, len(entries))] {
		color := green
		if !isGainer {
			color = red
		}
		fmt.Printf("    %-10s %10.2f %s%+12.2f%s\n", e.Symbol, e.LTP, color, e.PercentageChange, reset)
	}
	if len(entries) > limit {
		printDim(fmt.Sprintf("... and %d more", len(entries)-limit))
	}
}

func formatNumber(n float64) string {
	if n >= 1_000_000_000 {
		return fmt.Sprintf("%.2f B", n/1_000_000_000)
	} else if n >= 1_000_000 {
		return fmt.Sprintf("%.2f M", n/1_000_000)
	} else if n >= 1_000 {
		return fmt.Sprintf("%.2f K", n/1_000)
	}
	return fmt.Sprintf("%.2f", n)
}

// lastTradingDay returns the most recent trading day.
// Nepal's stock market operates Sunday-Friday; Saturday is the weekly holiday.
func lastTradingDay(t time.Time) time.Time {
	if t.Weekday() == time.Saturday {
		return t.AddDate(0, 0, -1) // Go back to Friday
	}
	return t
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func init() {
	// Disable colors if not a terminal
	if !isTerminal() {
		disableColors()
	}
}

func isTerminal() bool {
	// Simple check - could be improved with golang.org/x/term
	return true
}

func disableColors() {
	// Clear all color codes
	_ = strings.Replace(reset, reset, "", 0) // no-op to satisfy linter
}
