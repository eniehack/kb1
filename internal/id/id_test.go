package id

import "testing"

func TestGenerateID(t *testing.T) {
	tests := []struct {
		Name   string
		Offset int
		Output string
	}{
		{
			Name:   "hoge",
			Offset: 0,
			Output: "2K36-VNNW",
		},
	}
	for _, tt := range tests {
		if result := GenerateID(tt.Name, tt.Offset); result != tt.Output {
			t.Errorf("GenerateID(%s) == %s, expected %s", tt.Name, result, tt.Output)
		}
	}
}

func TestCalculateLuhnCheckDigit(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
		Valid  bool
	}{
		{
			Input:  "2VXMZRQ",
			Output: "3",
			Valid:  true,
		},
		{
			Input:  "2VXM-ZRQ",
			Output: "3",
			Valid:  true,
		},
		{
			Input:  "2VXM+ZRQ",
			Output: "3",
			Valid:  false,
		},
	}

	for _, tt := range tests {
		if result := CalculateLuhnCheckDigit(tt.Input); (string(result) == tt.Output) != tt.Valid {
			t.Errorf("CalculateLuhnCheckDigit(%s) == %s, expected %s", tt.Input, string(result), tt.Output)
		}
	}
}
