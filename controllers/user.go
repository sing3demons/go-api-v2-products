package controllers

import (
	"app/config"
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

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

//FindAll - find user
func (u *Users) FindAll(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	if sub.(models.User).Role != "Admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbindeb"})
		return
	}

	var users []models.User
	if err := u.DB.Order("id desc").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var serializedUsers []userResponse
	copier.Copy(&serializedUsers, &users)
	ctx.JSON(http.StatusOK, gin.H{"users": serializedUsers})
}

//Create - INSERT INTO `users`
func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

//FindOne - SELECT * FROM users ORDER BY id LIMIT 1;
func (u *Users) FindOne(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//Update - UPDATE users SET
func (u *Users) Update(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&user, &form)

	if err := u.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	setUserImage(ctx, user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}

//Delete - delete
func (u *Users) Delete(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := u.DB.Unscoped().Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

}

//Promote ->
func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//Demote ->
func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
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

	db := config.GetDB()
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
