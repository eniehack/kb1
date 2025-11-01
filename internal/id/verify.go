package id

import (
	"strings"
)

func VerifyID(id string) bool {
	// IDの長さをチェック
	if len(id) < 8 {
		return false
	}

	// ハイフンを除去して正規化
	hyphenTrimmedID := strings.ReplaceAll(id, "-", "")

	if len(hyphenTrimmedID) != 8 {
		return false
	}

	normalizedID := strings.ToUpper(hyphenTrimmedID)

	// 全文字がbase32文字セットに含まれるかチェック
	for _, c := range normalizedID {
		if !strings.ContainsRune(base32Chars, c) {
			return false
		}
	}

	// チェックデジットを計算して比較
	codeWithoutCheck := normalizedID[:7] // Type + Code部分
	expectedCheck := CalculateLuhnCheckDigit(codeWithoutCheck)
	actualCheck := normalizedID[7]

	return expectedCheck == actualCheck
}
