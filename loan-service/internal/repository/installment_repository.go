package repository

import (
	"github.com/RAF-SI-2025/EXBanka-3-Backend/loan-service/internal/models"
	"gorm.io/gorm"
)

type InstallmentRepository struct {
	db *gorm.DB
}

func NewInstallmentRepository(db *gorm.DB) *InstallmentRepository {
	return &InstallmentRepository{db: db}
}

func (r *InstallmentRepository) CreateBatch(items []models.LoanInstallment) error {
	return r.db.Create(&items).Error
}

// ListByLoanID returns all installments for a loan ordered by RedniBroj.
func (r *InstallmentRepository) ListByLoanID(loanID uint) ([]models.LoanInstallment, error) {
	var items []models.LoanInstallment
	if err := r.db.Where("loan_id = ?", loanID).
		Order("redni_broj ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
