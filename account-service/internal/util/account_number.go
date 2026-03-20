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

// GenerateAccountNumber generates an 18-digit account number.
// Format: BBB (3) + FFFF (4) + RRRRRRRRR (9) + TT (2) = 18
func GenerateAccountNumber(tip, vrsta string, podvrsta ...string) string {
	pv := ""
	if len(podvrsta) > 0 {
		pv = podvrsta[0]
	}
	randomPart := fmt.Sprintf("%09d", rand.Int63n(1_000_000_000))
	typeCode := AccountTypeCode(tip, vrsta, pv)
	return BankCode + BranchCode + randomPart + typeCode
}

// ValidateAccountNumber checks that number is 18 digits and starts with a valid bank code.
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
	return bankCode == "111" || bankCode == "222" || bankCode == "333" || bankCode == "444"
}
