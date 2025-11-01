package id

import (
	"strings"

	"github.com/twmb/murmur3"
)

const (
	base32Chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	codeLength  = 6
	typeChar    = "2"
)

func GenerateID(input string, offset int) string {
	// 衝突回避のための試行
	//for offset := range 32 { // 32bit以内で試行
	buf := new(strings.Builder)
	code := generateCode(input, offset)
	codeElements := []string{typeChar, code}
	fullCode := strings.Join(codeElements, "")
	buf.WriteString(fullCode)
	checkDigit := CalculateLuhnCheckDigit(fullCode)
	buf.WriteByte(checkDigit)
	//fmt.Fprintf(os.Stderr, "%s\n", buf.String())
	// if offset <= 5 {
	// 	continue
	// }
	return buf.String()
	//}
}

func generateCode(hashInput string, offset int) string {
	// CRC32でハッシュ値を計算
	origHash := murmur3.StringSum32(hashInput)

	// offsetに応じてbitをシフト
	hash := (origHash << offset) | (origHash >> (32 - offset))
	//fmt.Fprintf(os.Stderr, "hash(off=%d): %b | %b -> %b, ", offset, (origHash << offset), (origHash >> (32 - offset)), hash)

	// base32に変換
	code := new(strings.Builder)
	for i := range codeLength {
		char := base32Chars[hash%32]
		code.WriteByte(char)
		if i == 2 {
			code.WriteString("-")
		}
		hash /= 32
	}
	reversedCode := new(strings.Builder)
	strCode := code.String()
	for i := len(strCode) - 1; 0 <= i; i-- {
		reversedCode.WriteByte(strCode[i])
	}

	return reversedCode.String()
}

func CalculateLuhnCheckDigit(input string) byte {
	// Luhn mod 32アルゴリズムを使用してチェックデジットを計算
	// base29文字を数値にマッピング
	charToNum := make(map[rune]int, len(base32Chars))
	for i, c := range base32Chars {
		charToNum[c] = i
	}

	sum := 0
	isEven := false

	// 右から左へ処理（Luhn mod 32アルゴリズム）
	for i := len(input) - 1; 0 <= i; i-- {
		if input[i] == '-' {
			continue
		}
		digit := charToNum[rune(input[i])]

		if isEven {
			digit *= 2
			if digit >= 32 {
				digit = (digit / 32) + (digit % 32)
			}
		}

		sum += digit
		isEven = !isEven
	}

	// チェックデジット計算（mod 32）
	checkSum := (32 - (sum % 32)) % 32

	return base32Chars[checkSum]
}
