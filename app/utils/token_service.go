package utils

import (
	"github.com/zidni722/login-crud-user/app/models"
	"github.com/zidni722/login-crud-user/app/repositories/interface"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

var (
	key = []byte(viper.GetString("token.key"))
	issuer = viper.GetString("token.issuer")
)

type CustomClaims struct {
	UserId int `json:"user_id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Authable interface {
	Encode(userId int, username string) (string, error)
	Decode(token string) (*models.User, error)
}

type JWTService struct {
	Db *gorm.DB
	UserRepository _interface.IUserRepository
}

type TokenService struct {}

func NewJWTService(db *gorm.DB, userRepository _interface.IUserRepository) *JWTService {
	return &JWTService{
		Db: db,
		UserRepository: userRepository,
	}
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *JWTService) Encode(userId int, email string) (string, error) {
	expiredAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	claims := CustomClaims{
		UserId: userId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer: issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func (s *JWTService) Decode(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err == nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			var user models.User
			s.UserRepository.FindById(s.Db, &user, claims.UserId)
			if user != (models.User{}) {
				return &user, nil
			}
		}
	}

	return nil, &UnAuthenticatedError{}
}

func (s *TokenService) Encode(userId int, email string) (string, error) {
	expiredAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	claims := CustomClaims{
		UserId: userId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer: issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func (s *TokenService) Decode(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err == nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			if claims.UserId == 0 && claims.Email == viper.GetString("app.owner") {
				return &models.User{}, nil
			}
		}
	}

	return nil, &UnAuthenticatedError{}
}
