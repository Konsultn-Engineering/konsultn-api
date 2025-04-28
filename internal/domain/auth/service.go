package auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"log"
)

type AuthService interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	GenerateToken(ctx context.Context, UID string, userId string) (string, error)
}

type firebaseAuthService struct {
	client *auth.Client
}

func NewFirebaseAuthService(client *auth.Client) AuthService {
	return &firebaseAuthService{client: client}
}

func (s *firebaseAuthService) CreateUser(ctx context.Context, email, password string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)

	user, err := s.client.CreateUser(ctx, params)
	if err != nil {
		log.Println("Error creating user in Firebase:", err)
		return "", err
	}
	return user.UID, nil
}

func (s *firebaseAuthService) GenerateToken(ctx context.Context, UID string, userId string) (string, error) {
	claims := map[string]interface{}{
		"role":   "freelancer", // Assume "freelancer", "client", etc.
		"userID": userId,
	}

	err := s.client.SetCustomUserClaims(ctx, UID, claims)
	if err != nil {
		return "", err
	}

	token, err := s.client.CustomToken(ctx, UID)
	if err != nil {
		log.Println("Error generating token:", err)
		return "", err
	}
	return token, nil
}
