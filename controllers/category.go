package controllers

import (
	"github.com/sing3demons/app/v2/store"
)

func NewCategoryHandler(store *store.GormStore) *Category {
	return &Category{store: store}
}

//Category return database
type Category struct {
	store *store.GormStore
}

type categoryForm struct {
	Name string `json:"name" binding:"required"`
}

type categoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type categoryPaging struct {
	Items  []categoryResponse `json:"items"`
	Paging *pagingResult      `json:"paging"`
}
