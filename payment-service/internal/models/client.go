package models

type Client struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Ime      string `json:"ime"`
	Prezime  string `json:"prezime"`
	Email    string `json:"email"`
}
