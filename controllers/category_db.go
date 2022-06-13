package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/models"
)

// FindAll godoc
// @Summary show an categories
// @Description get all categories
// @Tags categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  categoryPaging
// @Router       /api/v1/categories [get]
func (c *Category) FindAll(ctx *gin.Context) {
	var categories []models.Category

	// pagination := pagination{ctx: ctx, query: c.store, records: &categories}
	pagination := NewPaginationHandler(ctx, c.store, &categories)
	p := pagination.pagingResource()

	serializedCategories := []categoryResponse{}
	copier.Copy(&serializedCategories, &categories)

	ctx.JSON(http.StatusCreated, gin.H{"categories": categoryPaging{Items: serializedCategories, Paging: p}})
}

// FindOne godoc
// @Summary show an category
// @Description get string by ID
// @Tags categories
// @Accept       json
// @Produce      json
// @Param        id   path      uint  true  "id"
// @Success      200  {object}  categoryPaging
// @Failure      404  {object}  map[string]any	"{"error": "record not found"}"
// @Router       /api/v1/categories/{id} [get]
func (c *Category) FindOne(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedCategory := categoryResponse{}
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}

// Create godoc
// @Summary  add an category
// @Description add by json category
// @Tags categories
// @Accept       json
// @Produce      json
// @Param	form	body	categoryForm true "form"
// @Success      201  {object}  categoryResponse
// @Failure      422  {object}  map[string]any	"{"error": "unprocessable entity"}"
// @Router       /api/v1/categories [post]
func (c *Category) Create(ctx *gin.Context) {
	var form categoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &form)

	if err := c.store.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})
}

// Update godoc
// @Summary	update an category
// @Description	update by json category
// @Tags	categories
// @Accept	json
// @Produce	json
// @Param	id	path	int	true  "id"
// @Param	form	body	categoryForm true "form"
// @Success	200  {object}  categoryResponse
// @Failure	422  {object} string "Bad Request"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/categories/{id} [put]
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

	if err := c.store.Update(&category, "name", form.Name).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}

// Delete godoc
// @Summary	delete an category
// @Description	delete by json category
// @Tags	categories
// @Accept	json
// @Produce	json
// @Param	id	path	int	true  "id"
// @Success	200  string  string "{"message": "success"}"
// @Failure	422  {object} string "Bad Request"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/categories/{id} [delete]
func (c *Category) Delete(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.store.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *Category) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")
	if err := c.store.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
