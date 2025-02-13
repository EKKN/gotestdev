package models

import (
	"time"
)

type Marketing struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique" validate:"required"`
	IsActive bool   `json:"is_active" gorm:"not null"`
}

type Penjualan struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TransactionNumber string    `json:"transaction_number" gorm:"column:transaction_number" validate:"required"`
	MarketingID       uint      `json:"marketing_id" validate:"required"`
	Marketing         Marketing `json:"marketing" gorm:"foreignKey:MarketingID"`
	Date              time.Time `json:"date" gorm:"type:date" validate:"required"`
	CargoFee          float64   `json:"cargo_fee" validate:"gte=0"`
	TotalBalance      float64   `json:"total_balance" validate:"gte=0"`
	GrandTotal        float64   `json:"grand_total" validate:"gte=0"`
}

type PenjualanInput struct {
	ID                uint    `json:"id" gorm:"primaryKey"`
	TransactionNumber string  `json:"transaction_number" gorm:"column:transaction_number" validate:"required"`
	MarketingID       uint    `json:"marketing_id" validate:"required"`
	Date              string  `json:"date" gorm:"type:date" validate:"required"`
	CargoFee          float64 `json:"cargo_fee"` // Bisa null, nanti di-handle dalam validasi
	TotalBalance      float64 `json:"total_balance" validate:"required,gte=1"`
	GrandTotal        float64 `json:"grand_total" validate:"gte=0"`
}

// Model untuk Komisi Persen
type KomisiPersen struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	MinOmzet   float64 `json:"min_omzet" validate:"gte=0"`
	MaxOmzet   float64 `json:"max_omzet" validate:"gte=0"`
	Persentase float64 `json:"persentase" validate:"gte=0,lte=100"`
}

// Model untuk tabel Pembayaran
type Pembayaran struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PenjualanID  uint      `json:"penjualan_id" validate:"required"`
	Penjualan    Penjualan `json:"penjualan" gorm:"foreignKey:PenjualanID"`
	CreditMonths uint      `json:"credit_months" validate:"required"`
	InterestRate float64   `json:"interest_rate" validate:"required,gte=0"`
	Installment  float64   `json:"installment" validate:"required,gte=0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type PembayaranDetail struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	PembayaranID uint       `json:"pembayaran_id" validate:"required"`
	Pembayaran   Pembayaran `json:"pembayaran" gorm:"foreignKey:PembayaranID"`
	Year         uint       `json:"year" validate:"required"`
	Month        uint       `json:"month" validate:"required"`
	Amount       float64    `json:"amount" validate:"required,gte=0"`
	Status       string     `json:"status" gorm:"type:varchar(255)"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// Model untuk tabel Kredit
type Kredit struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	Months       uint    `json:"months" gorm:"unique;not null"`
	InterestRate float64 `json:"interest_rate" gorm:"not null"`
}
