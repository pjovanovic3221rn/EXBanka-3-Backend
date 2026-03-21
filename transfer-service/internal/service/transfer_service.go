package service

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/RAF-SI-2025/EXBanka-3-Backend/transfer-service/internal/models"
)

// AccountRepositoryInterface defined here to avoid circular imports.
type AccountRepositoryInterface interface {
	FindByID(id uint) (*models.Account, error)
	UpdateFields(id uint, fields map[string]interface{}) error
}

// TransferRepositoryInterface defined here to avoid circular imports.
type TransferRepositoryInterface interface {
	Create(transfer *models.Transfer) error
	FindByID(id uint) (*models.Transfer, error)
	ListByAccountID(accountID uint, filter models.TransferFilter) ([]models.Transfer, int64, error)
	ListByClientID(clientID uint, filter models.TransferFilter) ([]models.Transfer, int64, error)
}

// ExchangeRateServiceInterface allows mocking in tests.
type ExchangeRateServiceInterface interface {
	GetRate(fromCurrencyKod, toCurrencyKod string) (float64, error)
}

type CreateTransferInput struct {
	RacunPosiljaocaID uint
	RacunPrimaocaID   uint
	Iznos             float64
	Svrha             string
}

type TransferService struct {
	accountRepo  AccountRepositoryInterface
	transferRepo TransferRepositoryInterface
	exchangeSvc  ExchangeRateServiceInterface
}

func NewTransferServiceWithRepos(
	accountRepo AccountRepositoryInterface,
	transferRepo TransferRepositoryInterface,
	exchangeSvc ExchangeRateServiceInterface,
) *TransferService {
	return &TransferService{
		accountRepo:  accountRepo,
		transferRepo: transferRepo,
		exchangeSvc:  exchangeSvc,
	}
}

func (s *TransferService) ListTransfersByAccount(accountID uint, filter models.TransferFilter) ([]models.Transfer, int64, error) {
	return s.transferRepo.ListByAccountID(accountID, filter)
}

func (s *TransferService) ListTransfersByClient(clientID uint, filter models.TransferFilter) ([]models.Transfer, int64, error) {
	return s.transferRepo.ListByClientID(clientID, filter)
}

func (s *TransferService) CreateTransfer(input CreateTransferInput) (*models.Transfer, error) {
	if input.Iznos <= 0 {
		return nil, fmt.Errorf("transfer amount must be positive")
	}
	if input.RacunPosiljaocaID == input.RacunPrimaocaID {
		return nil, fmt.Errorf("sender and receiver accounts must be different")
	}

	sender, err := s.accountRepo.FindByID(input.RacunPosiljaocaID)
	if err != nil {
		return nil, fmt.Errorf("sender account not found: %w", err)
	}
	receiver, err := s.accountRepo.FindByID(input.RacunPrimaocaID)
	if err != nil {
		return nil, fmt.Errorf("receiver account not found: %w", err)
	}

	if sender.RaspolozivoStanje < input.Iznos {
		return nil, fmt.Errorf("insufficient balance: available %.2f, requested %.2f",
			sender.RaspolozivoStanje, input.Iznos)
	}
	if input.Iznos > sender.DnevniLimit {
		return nil, fmt.Errorf("amount %.2f exceeds daily limit %.2f",
			input.Iznos, sender.DnevniLimit)
	}

	// Determine converted amount, exchange rate, and commission.
	kurs := 1.0
	konvertovaniIznos := input.Iznos
	valutaIznosa := sender.Currency.Kod
	provizijaProcent := 0.0
	provizija := 0.0

	if sender.CurrencyID != receiver.CurrencyID {
		// All cross-currency transfers go through RSD (bazna valuta).
		// Step 1: from sender currency to RSD
		// Step 2: from RSD to receiver currency
		var rsdAmount float64

		if sender.Currency.Kod == "RSD" {
			// Already in RSD, no first conversion needed
			rsdAmount = input.Iznos
			kurs = 1.0
		} else {
			// Convert sender currency → RSD (prodajni kurs)
			kursToRSD, err2 := s.exchangeSvc.GetRate(sender.Currency.Kod, "RSD")
			if err2 != nil {
				return nil, fmt.Errorf("failed to get exchange rate %s→RSD: %w", sender.Currency.Kod, err2)
			}
			rsdAmount = input.Iznos * kursToRSD
			kurs = kursToRSD
		}

		if receiver.Currency.Kod == "RSD" {
			// Receiver is RSD, no second conversion needed
			konvertovaniIznos = math.Round(rsdAmount*100) / 100
		} else {
			// Convert RSD → receiver currency (prodajni kurs)
			kursFromRSD, err2 := s.exchangeSvc.GetRate("RSD", receiver.Currency.Kod)
			if err2 != nil {
				return nil, fmt.Errorf("failed to get exchange rate RSD→%s: %w", receiver.Currency.Kod, err2)
			}
			konvertovaniIznos = math.Round(rsdAmount*kursFromRSD*100) / 100
			// Store effective rate (from→to via RSD)
			if input.Iznos > 0 {
				kurs = konvertovaniIznos / input.Iznos
			}
		}

		// Commission 0-1%
		provizijaProcent = math.Round(rand.Float64()*100) / 100 // 0.00 - 1.00%
		provizija = math.Round(input.Iznos*provizijaProcent) / 100

		// Sender pays: original amount + commission
		ukupnoZaSkidanje := input.Iznos + provizija

		if sender.RaspolozivoStanje < ukupnoZaSkidanje {
			return nil, fmt.Errorf("insufficient balance: available %.2f, required %.2f (amount %.2f + commission %.2f)",
				sender.RaspolozivoStanje, ukupnoZaSkidanje, input.Iznos, provizija)
		}

		// Update sender: deduct amount + commission
		if err := s.accountRepo.UpdateFields(sender.ID, map[string]interface{}{
			"stanje":             sender.Stanje - ukupnoZaSkidanje,
			"raspolozivo_stanje": sender.RaspolozivoStanje - ukupnoZaSkidanje,
		}); err != nil {
			return nil, fmt.Errorf("failed to update sender balance: %w", err)
		}
	} else {
		// Same currency: no commission, no conversion
		if err := s.accountRepo.UpdateFields(sender.ID, map[string]interface{}{
			"stanje":             sender.Stanje - input.Iznos,
			"raspolozivo_stanje": sender.RaspolozivoStanje - input.Iznos,
		}); err != nil {
			return nil, fmt.Errorf("failed to update sender balance: %w", err)
		}
	}

	// Update receiver balance.
	if err := s.accountRepo.UpdateFields(receiver.ID, map[string]interface{}{
		"stanje":             receiver.Stanje + konvertovaniIznos,
		"raspolozivo_stanje": receiver.RaspolozivoStanje + konvertovaniIznos,
	}); err != nil {
		return nil, fmt.Errorf("failed to update receiver balance: %w", err)
	}

	transfer := &models.Transfer{
		RacunPosiljaocaID: input.RacunPosiljaocaID,
		RacunPrimaocaID:   input.RacunPrimaocaID,
		Iznos:             input.Iznos,
		ValutaIznosa:      valutaIznosa,
		KonvertovaniIznos: konvertovaniIznos,
		Kurs:              kurs,
		Provizija:         provizija,
		ProvizijaProcent:  provizijaProcent,
		Svrha:             input.Svrha,
		Status:            "uspesno",
		VremeTransakcije:  time.Now(),
	}

	if err := s.transferRepo.Create(transfer); err != nil {
		return nil, fmt.Errorf("failed to save transfer: %w", err)
	}

	return transfer, nil
}
