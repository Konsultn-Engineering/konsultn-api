package auth

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	user2 "konsultn-api/internal/domain/user"
	"konsultn-api/pkg/firebase"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FirebaseLoginResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

// make sure to inject this from env/config
const firebaseAPIKey = "AIzaSyDN9bapCygS1xP12dK65rheq1Rzdob6WY0"

type Handler struct {
	repo        *user2.Repository[user2.User] // use the interface for flexibility
	AuthService AuthService
}

func NewHandler(repo *user2.Repository[user2.User]) *Handler {
	var authService = NewFirebaseAuthService(firebase.AuthClient)
	return &Handler{
		repo:        repo,
		AuthService: authService,
	}
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var createUserDto user2.CreateUserRequest

	if err := ctx.ShouldBindJSON(&createUserDto); err != nil {
		// If there is an error (invalid JSON or missing fields), return an error response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, existingUserError := h.repo.FindFirstBy("email", createUserDto.Email)
	if existingUserError != nil && existingUserError.Error() != "record not found" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": existingUserError.Error()})
		return
	}

	if existingUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": existingUser})
		return
	}

	uid, err := h.AuthService.CreateUser(ctx, createUserDto.Email, createUserDto.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error in creating user"})
		return
	}
	userModel := user2.ToUserModel(&createUserDto)
	userModel.UID = uid

	createdUser, err := h.repo.Save(userModel)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "error"})
		return
	}

	userToken, err := h.AuthService.GenerateToken(ctx, uid, createdUser.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"user":    createdUser,
		"token":   userToken,
	})
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Firebase sign-in with password
	payload := map[string]interface{}{
		"email":             req.Email,
		"password":          req.Password,
		"returnSecureToken": true,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		"https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+firebaseAPIKey,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil || resp.StatusCode != 200 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	defer resp.Body.Close()

	var firebaseRes FirebaseLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&firebaseRes); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Firebase response"})
		return
	}
	user, err := h.repo.FindFirstBy("uid", firebaseRes.LocalID)

	// Here you can update claims dynamically (example: role from your DB)
	newClaims := map[string]interface{}{
		"role":   "freelancer", // or fetch dynamically from your DB
		"userId": user.ID,      // fetch this from your DB based on email or uid
	}

	err = firebase.AuthClient.SetCustomUserClaims(ctx, firebaseRes.LocalID, newClaims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update custom claims"})
		return
	}

	// Verify token again (optional)
	token, err := firebase.AuthClient.VerifyIDToken(ctx, firebaseRes.IDToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	role, _ := token.Claims["role"]
	userId, _ := token.Claims["userId"]

	ctx.JSON(http.StatusOK, gin.H{
		"id_token":      firebaseRes.IDToken,
		"refresh_token": firebaseRes.RefreshToken,
		"expires_in":    firebaseRes.ExpiresIn,
		"firebase_user": firebaseRes.LocalID,
		"role":          role,
		"userId":        userId,
		"message":       "Login successful, claims updated. Please refresh your token to get updated claims.",
	})
}
