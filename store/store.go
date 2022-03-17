package store

import (
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStorm(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (tx *GormStore) Save(value interface{}) *gorm.DB {
	return tx.db.Save(value)
}

func (tx *GormStore) Create(value interface{}) *gorm.DB {
	return tx.db.Create(value)
}

func (tx *GormStore) Update(model interface{}, column string, value interface{}) *gorm.DB {
	return tx.db.Model(model).Update(column, value)
}

func (tx *GormStore) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return tx.db.Delete(value, conds...)
}

func (tx *GormStore) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return tx.db.First(dest, conds...)
}

func (s *GormStore) PreloadAndOrder(preload string, order interface{}) *gorm.DB {
	return s.db.Preload(preload).Order(order)
}

func (tx *GormStore) Model(value interface{}) *gorm.DB {
	return tx.db.Model(value)
}
func (tx *GormStore) Count(value interface{}, count *int64) *gorm.DB {
	return tx.db.Model(value).Count(count)
}

func (t *GormStore) Where(preload string, order interface{}, query interface{}, args ...interface{}) *gorm.DB {
	t.PreloadAndOrder(preload, order)
	return t.db.Where(query, args...)
}

func (tx *GormStore) Offset(offset int) *GormStore {
	tx.db.Offset(offset)
	return nil
}
func (tx *GormStore) Limit(limit int) *GormStore {
	tx.db.Limit(limit)
	return nil
}

func (tx *GormStore) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return tx.db.Find(dest, conds...)
}
