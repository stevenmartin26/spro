package service

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type JWTManager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func loadPrivateKey(path string) *rsa.PrivateKey {
	privateKeyBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	return privateKey
}

func loadPublicKey(path string) *rsa.PublicKey {
	publicKeyBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	return publicKey
}

func NewJWTManager(publicKeyPath, privateKeyPath string) *JWTManager {
	return &JWTManager{
		privateKey: loadPrivateKey(privateKeyPath),
		publicKey:  loadPublicKey(publicKeyPath),
	}
}

func (s *JWTManager) GenerateToken(userID uuid.UUID) (string, *common.CustomError) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID.String()
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		log.Fatal(err)
		return "", common.NewCustomError(common.ErrUnexpectedError, "internal server error")
	}
	return tokenString, nil
}

func (s *JWTManager) ValidateToken(accessToken string) (uuid.UUID, *common.CustomError) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, common.NewCustomError(common.ErrInvalidInput, "unexpected signing method")
		}
		return s.publicKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, common.NewCustomError(common.ErrUnauthorized, "invalid access token")
	}

	userID, err := uuid.Parse(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		return uuid.Nil, common.NewCustomError(common.ErrUnauthorized, "invalid access token")
	}
	return userID, nil
}
