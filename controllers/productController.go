package controllers

import (
	"app/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//Product - struct
type Product struct {
	DB *gorm.DB
}

type createProductForm struct {
	Name  string                `form:"name" binding:"required"`
	Desc  string                `form:"desc" binding:"required"`
	Price int                   `form:"price" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

type createProductRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

func (p *Product) FindAll(ctx *gin.Context) {

}

// FindOne - /:id
func (p *Product) FindOne(ctx *gin.Context) {

}

// Create - insert data
func (p *Product) Create(ctx *gin.Context) {
	var form createProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// form => product
	// product := models.Product{
	// 	Name:  form.Name,
	// 	Desc:  form.Desc,
	// 	Price: form.Price,
	// }

	var product models.Product
	copier.Copy(&product, &form)

	if err := p.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	p.setProductImage(ctx, &product)

	serializedProduct := createProductRespons{}
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

func (p *Product) setProductImage(ctx *gin.Context, products *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return nil
	}

	if products.Image == "" {
		//1. ตัด http://localhost:80800/uploads/prosucts/<ID>/image.png ให้เหลือ /uploads/prosucts/<ID>/image.png
		products.Image = strings.Replace(products.Image, os.Getenv("HOST"), "", 1)
		//2. แทนค่าพาธปัจจุบัน<WD>/uploads/prosucts/<ID>/image.png
		pwd, _ := os.Getwd()
		//3.remove <WD>/uploads/prosucts/<ID>/image.png
		os.Remove(pwd + products.Image)
	}

	path := "uploads/products/" + strconv.Itoa(int(products.ID))
	os.MkdirAll(path, 0755)

	filename := path + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	products.Image = os.Getenv("HOST") + "/" + filename

	p.DB.Save(products)

	return nil
}
