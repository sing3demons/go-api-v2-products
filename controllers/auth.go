package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/models"
	"github.com/sing3demons/app/v2/store"
)

//Auth database
type Auth struct {
	store *store.GormStore
}

func NewAuthHandler(store *store.GormStore) *Auth {
	return &Auth{store: store}
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
	Email string `json:"email"`
	Name  string `json:"name" `
}

//GetProfile godoc
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} authResponse
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/profile [get]
func (a *Auth) GetProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	var user models.User = sub.(models.User)
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

// SignUp godoc
// @Tags auth
// @Accept  json
// @Produce  json
// @Param register body authForm true "register"
// @Success 201 {object} authResponse
// @Failure 422 {object} map[string]any
// @Router /api/v1/auth/register [post]
func (a *Auth) SignUp(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}
	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.store.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser, "message": "ลงทะเบียนสำเร็จ"})
}

//UpdateImageProfile - upload image
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path uint true "id"
// @Param avatar formData file true "avatar"
// @Success 200 {object} userResponse
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/profile/{id} [patch]
func (a *Auth) UpdateImageProfile(ctx *gin.Context) {

	sub, _ := ctx.Get("sub")
	user := sub.(models.User)

	setUserImage(ctx, &user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}

// UpdateProfile godoc
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path uint true "id"
// @Param name formData string true "name"
// @Param email formData string true "email"
// @Param avatar formData file true "avatar"
// @Success 200 {object} userResponse
// @Failure 422 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/profile/{id} [put]
func (a *Auth) UpdateProfile(ctx *gin.Context) {
	form := updateProfileForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(models.User)
	copier.Copy(&user, &form)

	if err := a.store.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	setUserImage(ctx, &user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}
