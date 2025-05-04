package user

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"konsultn-api/pkg/firebase"
	"net/http"
)

type Handler struct {
	repo *Repository[User] // use the interface for flexibility
}

func NewHandler(repo *Repository[User]) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListAllUsers(ctx *gin.Context) {
	users, err := h.repo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.repo.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with id " + id + " not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var user CreateUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		// If there is an error (invalid JSON or missing fields), return an error response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password)

	fbUser, err := firebase.AuthClient.CreateUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customToken, err := firebase.AuthClient.CustomToken(ctx, fbUser.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth token"})
		return
	}

	// Process the user data (e.g., save it to the database)
	var userModel = ToUserModel(user)
	userModel.UID = fbUser.UID
	createdUser, creationError := h.repo.Save(userModel)

	if creationError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": creationError.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user":    createdUser,
		"token":   customToken,
	})

}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	var id = ctx.Param("id")
	err := h.repo.DeleteById(id, false)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "resource not found"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

//func (h *Handler) SyncUserHandler(c *gin.Context) {
//	uid := c.GetString("uid")
//	email := c.GetString("email")
//
//	var user User
//	user, err := h.repo.FindFirstBy("fuid", uid)
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		// User doesn't exist yet – create
//		user = User{
//			FUID:  uid, // Use Firebase UID as your User ID
//			Email: email,
//		}
//		if _, err1 := h.repo.Save(&user); err1 != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
//			return
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "User synced",
//		"user":    user,
//	})
//}
