package controllers

import (
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//Category return database
type Category struct {
	DB *gorm.DB
}

type categoryForm struct {
	Name string `json:"name" binding:"required"`
}

type categoryRespons struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type categoryPaging struct {
	Items  []categoryRespons `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

//FindAll find categories
func (c *Category) FindAll(ctx *gin.Context) {
	var categories []models.Category

	pagination := pagination{ctx: ctx, query: c.DB, records: &categories}
	p := pagination.pagingResource()

	serializedCategories := []categoryRespons{}
	copier.Copy(&serializedCategories, &categories)

	ctx.JSON(http.StatusCreated, gin.H{"categories": categoryPaging{Items: serializedCategories, Paging: p}})
}

//FindOne find category
func (c *Category) FindOne(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedCategory := categoryRespons{}
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}

//Create insert categoty
func (c *Category) Create(ctx *gin.Context) {
	var form categoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &form)

	if err := c.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryRespons
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})
}

//Update update category
func (c *Category) Update(ctx *gin.Context) {
	var form categoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&category).Update("name", form.Name).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}

//Delete delete category
func (c *Category) Delete(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Unscoped().Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *Category) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")
	if err := c.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
