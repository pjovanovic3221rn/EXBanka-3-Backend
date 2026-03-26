package util

import (
	"fmt"
	"math/rand"
)

const (
	BankCode   = "333"  // Naša banka - EXBanka-3
	BranchCode = "0001" // Default filijala
)

// AccountTypeCode returns the 2-digit type code for the account number.
//
// Tekući račun: 10 (generički)
//   Lični: 11, Poslovni: 12, Štedni: 13, Penzionerski: 14,
//   Za mlade: 15, Za studente: 16, Za nezaposlene: 17
//
// Devizni račun: 20 (generički)
//   Lični: 21, Poslovni: 22
func AccountTypeCode(tip, vrsta, podvrsta string) string {
	if tip == "devizni" {
		if vrsta == "poslovni" {
			return "22"
		}
		return "21"
	}

	// Tekući
	if vrsta == "poslovni" {
		return "12"
	}

	// Tekući lični — podvrste
	switch podvrsta {
	case "stedni":
		return "13"
	case "penzionerski":
		return "14"
	case "za_mlade":
		return "15"
	case "za_studente":
		return "16"
	case "za_nezaposlene":
		return "17"
	default:
		return "11" // standardni
	}
}

// GenerateAccountNumber generates an 18-digit account number with checksum.
// Format: BBB (3) + FFFF (4) + RRRRRRRR (8 random) + C (1 check) + TT (2) = 18
// The check digit C is chosen so that (sum of all 18 digits) % 11 == 0.
func GenerateAccountNumber(tip, vrsta string, podvrsta ...string) string {
	pv := ""
	if len(podvrsta) > 0 {
		pv = podvrsta[0]
	}
	typeCode := AccountTypeCode(tip, vrsta, pv)
	randomPart := fmt.Sprintf("%08d", rand.Int63n(100_000_000)) // 8 random digits

	// 17 digits without check digit: BBB + FFFF + RRRRRRRR + TT
	prefix := BankCode + BranchCode + randomPart // 15 digits
	suffix := typeCode                           // 2 digits

	// Sum of all known 17 digits
	sum := digitSum(prefix) + digitSum(suffix)
	check := (11 - (sum % 11)) % 11
	// If check == 10, regenerate (rare edge case)
	if check == 10 {
		return GenerateAccountNumber(tip, vrsta, podvrsta...)
	}

	return prefix + fmt.Sprintf("%d", check) + suffix
}

// digitSum returns the sum of all digits in a numeric string.
func digitSum(number string) int {
	sum := 0
	for _, c := range number {
		if c >= '0' && c <= '9' {
			sum += int(c - '0')
		}
	}
	return sum
}

// ValidateAccountNumber checks that a number is 18 digits, starts with a valid
// bank code, and passes the (sum of digits) % 11 == 0 checksum.
func ValidateAccountNumber(number string) bool {
	if len(number) != 18 {
		return false
	}
	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}
	bankCode := number[:3]
	validBank := bankCode == "111" || bankCode == "222" || bankCode == "333" || bankCode == "444"
	if !validBank {
		return false
	}
	return digitSum(number)%11 == 0
}
