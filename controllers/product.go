package controllers

import (
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

type productResponse struct {
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
	Items  []productResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{}) error
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
