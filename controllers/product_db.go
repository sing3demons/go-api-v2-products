package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/models"
)

// FindAll godoc
// @Summary Show an products
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query uint false "page"
// @Param limit query uint false "limit"
// @Success 200 {object} productsPaging
// @Router /api/v1/products [get]
func (p *ProductHandler) FindAll(ctx *gin.Context) {
	query1CacheKey := "items::product"
	query2CacheKey := "items::page"

	serializedProduct := []productResponse{}
	var paging *pagingResult

	cacheItems, err := p.cacher.MGet([]string{query1CacheKey, query2CacheKey})
	if err != nil {
		log.Println(err.Error())
	}

	productJS := cacheItems[0]
	pageJS := cacheItems[1]

	if productJS != nil && len(productJS.(string)) > 0 {
		err := json.Unmarshal([]byte(productJS.(string)), &serializedProduct)
		if err != nil {
			p.cacher.Del(query1CacheKey)
			log.Println(err.Error())
		}

	}

	itemToCaches := map[string]interface{}{}

	var paginationItem *pagingResult
	if productJS == nil {
		var products []models.Product
		// pagination := pagination{ctx: ctx, query: p.store, records: &products}
		pagination := NewPaginationHandler(ctx, p.store, &products)
		paginationItem = pagination.pagingResource()
		copier.Copy(&serializedProduct, &products)

		itemToCaches[query1CacheKey] = serializedProduct
	}

	if pageJS != nil && len(pageJS.(string)) > 0 {
		err := json.Unmarshal([]byte(pageJS.(string)), &paging)
		if err != nil {
			p.cacher.Del(query2CacheKey)
			log.Println(err.Error())
		}
	}

	if paging == nil {
		paging = paginationItem
		itemToCaches[query2CacheKey] = paging
	}

	if len(itemToCaches) > 0 {
		timeToExpire := 10 * time.Second // m
		fmt.Println("M_SET")

		// Set cache using M SET
		err := p.cacher.MSet(itemToCaches)
		if err != nil {
			log.Println(err.Error())
		}

		// Set time to expire
		keys := []string{}
		for k := range itemToCaches {
			keys = append(keys, k)
		}
		err = p.cacher.Expires(keys, timeToExpire)
		if err != nil {
			log.Println(err.Error())
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"products": productsPaging{Items: serializedProduct, Paging: paging}})

}

// @Summary FindOne - /:id
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} productResponse
// @Router /api/v1/products/{id} [get]
func (p *ProductHandler) FindOne(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

// Create godoc
// @Summary add an product
// @Description add by form product
// @Tags products
// @Accept  mpfd
// @Produce  json
// @Security BearerAuth
// @Param name formData string true "name"
// @Param desc formData string true "desc"
// @Param price formData int true "price"
// @Param image formData file true "image"
// @Param categoryId formData uint true "categoryId"
// @Success 200 {object} productResponse
// @Router /api/v1/products [post]
func (p *ProductHandler) Create(ctx *gin.Context) {
	var form createProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	copier.Copy(&product, &form)

	if err := p.store.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	p.setProductImage(ctx, &product)

	var serializedProduct productResponse
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

// UpdateAll godoc
// @Summary update an products
// @Description update by form product
// @Tags products
// @Accept  mpfd
// @Produce  json
// @Security BearerAuth
// @Param id path string true "id"
// @Param name formData string false "name"
// @Param desc formData string false "desc"
// @Param price formData int false "price"
// @Param image formData file false "image"
// @Param categoryId formData uint false "categoryId"
// @Success 200 {object} productResponse
// @Router /api/v1/products/{id} [put]
func (p *ProductHandler) UpdateAll(ctx *gin.Context) {
	var form updateProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&product, &form)

	if err := p.store.Save(product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	p.setProductImage(ctx, product)

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"products": serializedProduct})

}

// Delete godoc
// @Summary	delete an product
// @Description	delete by json product
// @Tags	products
// @Accept	json
// @Produce	json
// @Param id path string true "id"
// @Success	200  string  string "{"message": "deleted"}"
// @Failure	422  {object} string "Bad Request"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/categories/{id} [delete]
// @Success 200 {object} productResponse
// @Router /api/v1/products/{id} [delete]
func (p *ProductHandler) Delete(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := p.store.Delete(product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted...."})
}
