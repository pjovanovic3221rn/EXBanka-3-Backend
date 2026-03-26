package models

import "time"

// CardRequest stores a pending card creation request with verification code.
type CardRequest struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID        uint      `gorm:"not null" json:"account_id"`
	ClientID         uint      `gorm:"not null" json:"client_id"`
	VrstaKartice     string    `gorm:"not null" json:"vrsta_kartice"`
	NazivKartice     string    `json:"naziv_kartice"`
	ClientEmail      string    `gorm:"not null" json:"client_email"`
	ClientName       string    `json:"client_name"`
	VerifikacioniKod string    `gorm:"not null" json:"-"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	BrojPokusaja     int       `gorm:"default:0" json:"broj_pokusaja"`
	Status           string    `gorm:"default:'pending';not null" json:"status"` // pending, verified, expired, failed
	// OvlascenoLice fields (for poslovni accounts)
	OvlascenoIme          string `json:"ovlasceno_ime,omitempty"`
	OvlascenoPrezime      string `json:"ovlasceno_prezime,omitempty"`
	OvlascenoEmail        string `json:"ovlasceno_email,omitempty"`
	OvlascenoBrojTelefona string `json:"ovlasceno_broj_telefona,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}
