package controllers

import (
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/cache"
	"github.com/sing3demons/app/v2/models"
	"gorm.io/gorm"
)

//Product - struct
type Product struct {
	DB     *gorm.DB
	Cacher *cache.Cacher
}

type createProductForm struct {
	Name  string                `form:"name" binding:"required"`
	Desc  string                `form:"desc" binding:"required"`
	Price int                   `form:"price" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
}

type updateProductForm struct {
	Name  string                `form:"name"`
	Desc  string                `form:"desc"`
	Price int                   `form:"price"`
	Image *multipart.FileHeader `form:"image"`
	CategoryID uint                  `form:"categoryId"`
}

type productRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type producsPaging struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

//FindAll - query-proucts
func (p *Product) FindAll(ctx *gin.Context) {
	products := []models.Product{}

	query := p.DB.Preload("Category").Order("id desc")
		if category := ctx.Query("category"); category != "" {
			c, _ := strconv.Atoi(category)
			query = query.Where("category_id = ?", c)
		}

	pagination := pagination{
		ctx:     ctx,
		query:   p.DB,
		records: &products,
	}
	paging := pagination.pagingResource()

	serializedProducts := []productRespons{}
	copier.Copy(&serializedProducts, &products)

	ctx.JSON(http.StatusOK, gin.H{"products": producsPaging{Items: serializedProducts, Paging: paging}})
}

// FindOne - /:id
func (p *Product) FindOne(ctx *gin.Context) {
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
func (p *Product) Create(ctx *gin.Context) {
	var form createProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	copier.Copy(&product, &form)

	if err := p.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	p.setProductImage(ctx, &product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

// UpdateAll - update all
func (p *Product) UpdateAll(ctx *gin.Context) {
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

	if err := p.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	p.setProductImage(ctx, product)

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"products": serializedProduct})

}

// Delete - Delete product
func (p *Product) Delete(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := p.DB.Unscoped().Delete(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted...."})
}

func (p *Product) findProductByID(ctx *gin.Context) (*models.Product, error) {
	var product models.Product
	id := ctx.Param("id")

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) setProductImage(ctx *gin.Context, products *models.Product) error {
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

	p.DB.Save(products)

	return nil
}
