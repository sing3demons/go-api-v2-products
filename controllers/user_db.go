package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/sing3demons/app/v2/models"
)

// FindAll godoc
// @Summary      List users
// @Description	get users
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Success	200  {object}  []userResponse
// @Failure	403  {object} map[string]any "{"error": "Forbidden"}"
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users [get]
func (u *Users) FindAll(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	if sub.(models.User).Role != "Admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
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

// Create godoc
// @Summary	add an users
// @Description	add user by form User
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param register body createUserForm true "register"
// @Success	200  {object}  userResponse
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users [post]
func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		fmt.Println("error ShouldBindJSON")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		fmt.Println("error create")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

// FindOne godoc
// @Summary		Show an users
// @Description	get user by id
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param	id	path	int	true  "id"
// @Success	200  {object}  userResponse
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users/{id} [get]
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

// Update godoc
// @Summary		update an users
// @Description	update user by form User
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param	id	path	int	true  "id"
// @Param name formData string true "name"
// @Param email formData string false "email"
// @Param avatar formData file true "avatar"
// @Success	200  {object}  userResponse
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users/{id} [put]
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

// Delete godoc
// @Summary		update an users
// @Description	update user by form User
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param	id	path	int	true  "id"
// @Success	200
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users/{id} [delete]
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

// Promote godoc
// @Summary		update an users
// @Description Admin mode
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param	id	path	int	true  "id"
// @Success	200
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users/{id}/promote [patch]
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

// Demote godoc
// @Summary		update an users
// @Description Admin mode
// @Tags	users
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param	id	path	int	true  "id"
// @Success	200
// @Failure	422  {object} map[string]any "{"error": "Bad Request"}"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router	/api/v1/users/{id}/demote [patch]
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
