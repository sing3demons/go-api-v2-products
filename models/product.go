package models

type Product struct {
	ID    uint   `gorm:"unique;not null"`
	Name  string `gorm:"not null"`
	Desc  string `gorm:"not null"`
	Price int    `gorm:"not null"`
	Image string `gorm:"not null"`
}
