package controllers

import (
	"app/models"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//Auth database
type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type updateProfileForm struct {
	Name   string                `form:"name"`
	Email  string                `form:"email" `
	Avatar *multipart.FileHeader `form:"avatar"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name" `
}

//GetProfile - GET /api/v1/profile
func (a *Auth) GetProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	var user models.User = sub.(models.User)
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//SignUp - POST /api/v1/register
func (a *Auth) SignUp(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}
	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser, "message": "ลงทะเบียนสำเร็จ"})
}

//UpdateImageProfile - upload image
func (a *Auth) UpdateImageProfile(ctx *gin.Context) {

	sub, _ := ctx.Get("sub")
	user := sub.(models.User)

	setUserImage(ctx, &user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}

//UpdateProfile - PUT /api/v1/profile
func (a *Auth) UpdateProfile(ctx *gin.Context) {
	form := updateProfileForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(models.User)
	copier.Copy(&user, &form)

	if err := a.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	setUserImage(ctx, &user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}
