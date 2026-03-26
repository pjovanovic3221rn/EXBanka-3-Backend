package models

import "time"

// ValidLoanTypes returns the allowed values for Loan.Vrsta.
// Exported so tests can verify the business data without a DB.
func ValidLoanTypes() []string {
	return []string{"gotovinski", "stambeni", "auto", "refinansirajuci", "studentski"}
}

// ValidLoanStatuses returns the allowed values for Loan.Status.
func ValidLoanStatuses() []string {
	return []string{"zahtev", "odobren", "odbijen", "aktivan", "zatvoren"}
}

// ValidInterestTypes returns the allowed values for Loan.TipKamate.
func ValidInterestTypes() []string {
	return []string{"fiksna", "varijabilna"}
}

// ValidPeriods returns the allowed repayment periods (months) for each loan type.
func ValidPeriods() map[string][]int {
	return map[string][]int{
		"gotovinski":      {12, 24, 36, 48, 60, 72, 84},
		"auto":            {12, 24, 36, 48, 60, 72, 84},
		"studentski":      {12, 24, 36, 48, 60, 72, 84},
		"refinansirajuci": {12, 24, 36, 48, 60, 72, 84},
		"stambeni":        {60, 120, 180, 240, 300, 360},
	}
}

// ValidEmploymentStatuses returns the allowed values for StatusZaposlenja.
func ValidEmploymentStatuses() []string {
	return []string{"stalno", "privremeno", "nezaposlen"}
}

// Loan represents a bank loan issued to a client.
type Loan struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Vrsta          string    `gorm:"not null" json:"vrsta"`                       // gotovinski | stambeni | auto | refinansirajuci | studentski
	BrojRacuna     string    `gorm:"not null" json:"broj_racuna"`                 // account for payout/collection
	BrojKredita    string    `gorm:"uniqueIndex;not null" json:"broj_kredita"`
	Iznos          float64   `gorm:"not null" json:"iznos"`
	Period         int       `gorm:"not null" json:"period"`                      // months
	KamatnaStopa   float64   `gorm:"not null" json:"kamatna_stopa"`               // annual %
	TipKamate      string    `gorm:"not null" json:"tip_kamate"`                  // fiksna | varijabilna
	DatumKreiranja time.Time `json:"datum_kreiranja"`
	DatumDospeca   time.Time `json:"datum_dospeca"`
	IznosRate      float64   `gorm:"not null" json:"iznos_rate"`                  // monthly installment (annuity)
	Status         string    `gorm:"not null;default:'zahtev'" json:"status"`     // zahtev | odobren | odbijen | aktivan | zatvoren
	ClientID       uint      `gorm:"not null" json:"client_id"`
	ZaposleniID    *uint     `json:"zaposleni_id"`                                // who approved/rejected
	CurrencyID     uint      `gorm:"not null" json:"currency_id"`

	// Additional fields from specification
	SvrhaKredita      string `json:"svrha_kredita"`                               // purpose of loan
	IznosMesecnePlate float64 `json:"iznos_mesecne_plate"`                        // monthly salary
	StatusZaposlenja  string `json:"status_zaposlenja"`                            // stalno | privremeno | nezaposlen
	PeriodZaposlenja  string `json:"period_zaposlenja"`                            // employment duration at current employer
	KontaktTelefon    string `json:"kontakt_telefon"`                              // contact phone

	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Installments []LoanInstallment `gorm:"foreignKey:LoanID" json:"installments,omitempty"`
}
