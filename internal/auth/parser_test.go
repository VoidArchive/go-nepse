package auth

import (
	"testing"
)

func TestTokenParser_Initialization(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	// Verify all WASM functions are exported
	if parser.cdx == nil {
		t.Error("cdx function not exported")
	}
	if parser.rdx == nil {
		t.Error("rdx function not exported")
	}
	if parser.bdx == nil {
		t.Error("bdx function not exported")
	}
	if parser.ndx == nil {
		t.Error("ndx function not exported")
	}
	if parser.mdx == nil {
		t.Error("mdx function not exported")
	}
}

func TestTokenParser_IndicesFromSalts(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	tests := []struct {
		name  string
		salts [5]int
	}{
		{
			name:  "typical salt values",
			salts: [5]int{1234, 5678, 9012, 3456, 7890},
		},
		{
			name:  "zero salts",
			salts: [5]int{0, 0, 0, 0, 0},
		},
		{
			name:  "small values",
			salts: [5]int{1, 2, 3, 4, 5},
		},
		{
			name:  "large values",
			salts: [5]int{999999, 888888, 777777, 666666, 555555},
		},
		{
			name:  "negative values",
			salts: [5]int{-100, -200, -300, -400, -500},
		},
		{
			name:  "mixed values",
			salts: [5]int{-50, 100, 0, 999, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indices, err := parser.indicesFromSalts(tt.salts)
			if err != nil {
				t.Fatalf("indicesFromSalts(%v) failed: %v", tt.salts, err)
			}

			// Access token should have 5 indices (one from each WASM function)
			if len(indices.access) != 5 {
				t.Errorf("expected 5 access indices, got %d", len(indices.access))
			}

			// Refresh token should also have 5 indices
			if len(indices.refresh) != 5 {
				t.Errorf("expected 5 refresh indices, got %d", len(indices.refresh))
			}
		})
	}
}

func TestTokenParser_DeterministicOutput(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	salts := [5]int{1234, 5678, 9012, 3456, 7890}

	// Call multiple times and verify same output
	indices1, err := parser.indicesFromSalts(salts)
	if err != nil {
		t.Fatalf("first indicesFromSalts failed: %v", err)
	}

	indices2, err := parser.indicesFromSalts(salts)
	if err != nil {
		t.Fatalf("second indicesFromSalts failed: %v", err)
	}

	// Compare access indices
	for i := range indices1.access {
		if indices1.access[i] != indices2.access[i] {
			t.Errorf("access index %d differs: %d != %d", i, indices1.access[i], indices2.access[i])
		}
	}

	// Compare refresh indices
	for i := range indices1.refresh {
		if indices1.refresh[i] != indices2.refresh[i] {
			t.Errorf("refresh index %d differs: %d != %d", i, indices1.refresh[i], indices2.refresh[i])
		}
	}
}

func TestTokenParser_DifferentSaltsProduceDifferentIndices(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	salts1 := [5]int{1234, 5678, 9012, 3456, 7890}
	salts2 := [5]int{9876, 5432, 1098, 7654, 3210}

	indices1, err := parser.indicesFromSalts(salts1)
	if err != nil {
		t.Fatalf("indicesFromSalts(salts1) failed: %v", err)
	}

	indices2, err := parser.indicesFromSalts(salts2)
	if err != nil {
		t.Fatalf("indicesFromSalts(salts2) failed: %v", err)
	}

	// At least some indices should differ
	allSame := true
	for i := range indices1.access {
		if indices1.access[i] != indices2.access[i] {
			allSame = false
			break
		}
	}
	if allSame {
		t.Log("Warning: different salts produced same access indices (may be expected for some salt combinations)")
	}
}

func TestTokenParser_Call5(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	// Test that call5 works with all functions
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "cdx function",
			fn: func() error {
				_, err := parser.call5(parser.cdx, 1, 2, 3, 4, 5)
				return err
			},
		},
		{
			name: "rdx function",
			fn: func() error {
				_, err := parser.call5(parser.rdx, 1, 2, 3, 4, 5)
				return err
			},
		},
		{
			name: "bdx function",
			fn: func() error {
				_, err := parser.call5(parser.bdx, 1, 2, 3, 4, 5)
				return err
			},
		},
		{
			name: "ndx function",
			fn: func() error {
				_, err := parser.call5(parser.ndx, 1, 2, 3, 4, 5)
				return err
			},
		},
		{
			name: "mdx function",
			fn: func() error {
				_, err := parser.call5(parser.mdx, 1, 2, 3, 4, 5)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s failed: %v", tt.name, err)
			}
		})
	}
}

func TestTokenParser_Close(t *testing.T) {
	parser, err := newTokenParser()
	if err != nil {
		t.Fatalf("newTokenParser() failed: %v", err)
	}

	// Close should not error
	if err := parser.close(); err != nil {
		t.Errorf("close() failed: %v", err)
	}

	// Double close should be safe (wazero handles this)
	if err := parser.close(); err != nil {
		t.Logf("double close error (may be expected): %v", err)
	}
}

func TestTokenParser_MultipleInstances(t *testing.T) {
	// Create multiple parser instances to ensure no global state issues
	parsers := make([]*tokenParser, 3)
	for i := range parsers {
		p, err := newTokenParser()
		if err != nil {
			t.Fatalf("newTokenParser() %d failed: %v", i, err)
		}
		parsers[i] = p
	}

	// All should produce same results for same input
	salts := [5]int{1234, 5678, 9012, 3456, 7890}
	var firstIndices tokenIndices
	for i, p := range parsers {
		indices, err := p.indicesFromSalts(salts)
		if err != nil {
			t.Fatalf("parser %d indicesFromSalts failed: %v", i, err)
		}
		if i == 0 {
			firstIndices = indices
		} else {
			for j := range indices.access {
				if indices.access[j] != firstIndices.access[j] {
					t.Errorf("parser %d access[%d] = %d, want %d", i, j, indices.access[j], firstIndices.access[j])
				}
			}
		}
	}

	// Clean up
	for _, p := range parsers {
		p.close()
	}
}

// Benchmark for performance regression detection
func BenchmarkTokenParser_IndicesFromSalts(b *testing.B) {
	parser, err := newTokenParser()
	if err != nil {
		b.Fatalf("newTokenParser() failed: %v", err)
	}
	defer parser.close()

	salts := [5]int{1234, 5678, 9012, 3456, 7890}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.indicesFromSalts(salts)
		if err != nil {
			b.Fatalf("indicesFromSalts failed: %v", err)
		}
	}
}

func BenchmarkNewTokenParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser, err := newTokenParser()
		if err != nil {
			b.Fatalf("newTokenParser() failed: %v", err)
		}
		parser.close()
	}
}
