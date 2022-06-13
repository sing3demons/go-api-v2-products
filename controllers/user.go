package controllers

import (
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/app/v2/database"
	"github.com/sing3demons/app/v2/models"
	"gorm.io/gorm"
)

//Users - receiver adater
type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email  string                `form:"email" binding:"omitempty,email"`
	Avatar *multipart.FileHeader `form:"avatar"`
	Name   string                `form:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email" `
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return nil
	}

	db := database.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + filename
	db.Save(user)

	return nil
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
