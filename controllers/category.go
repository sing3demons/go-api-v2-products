package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/app/v2/models"
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

func (c *Category) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")
	if err := c.store.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
