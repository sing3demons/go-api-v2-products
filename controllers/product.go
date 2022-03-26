package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/cache"
	"github.com/sing3demons/app/v2/models"
	"github.com/sing3demons/app/v2/store"
)

//Product - struct
type ProductHandler struct {
	store  *store.GormStore
	cacher *cache.Cacher
}

func NewProductHandler(store *store.GormStore, cacher *cache.Cacher) *ProductHandler {
	return &ProductHandler{store: store, cacher: cacher}
}

type createProductForm struct {
	Name       string                `form:"name" binding:"required"`
	Desc       string                `form:"desc" binding:"required"`
	Price      int                   `form:"price" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
}

type updateProductForm struct {
	Name       string                `form:"name"`
	Desc       string                `form:"desc"`
	Price      int                   `form:"price"`
	Image      *multipart.FileHeader `form:"image"`
	CategoryID uint                  `form:"categoryId"`
}

type productRespons struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type productsPaging struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{}) error
}

//FindAll - query-products
func (p *ProductHandler) FindAll(ctx *gin.Context) {
	query1CacheKey := "items::product"
	query2CacheKey := "items::page"

	serializedProduct := []productRespons{}
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

		// Set cache using MSET
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

// FindOne - /:id
func (p *ProductHandler) FindOne(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

// Create - insert data
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

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

// UpdateAll - update all
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

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"products": serializedProduct})

}

// Delete - Delete product
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

func (p *ProductHandler) findProductByID(ctx *gin.Context) (*models.Product, error) {
	var product models.Product
	id := ctx.Param("id")

	if err := p.store.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *ProductHandler) setProductImage(ctx *gin.Context, products *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return nil
	}

	if products.Image == "" {
		products.Image = strings.Replace(products.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + products.Image)
	}

	path := "uploads/products/" + strconv.Itoa(int(products.ID))
	os.MkdirAll(path, 0755)

	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	products.Image = os.Getenv("HOST") + "/" + filename

	p.store.Save(products)

	return nil
}
