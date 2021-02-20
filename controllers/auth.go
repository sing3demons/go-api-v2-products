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
	sub, _ := ctx.Get("jwt_id")
	ctx.JSON(http.StatusOK, gin.H{"user": sub})
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

//SignIn - POST /api/v1/login
func (a *Auth) SignIn(ctx *gin.Context) {
}

//UpdateProfile - PUT /api/v1/profile
func (a *Auth) UpdateProfile(ctx *gin.Context) {}
