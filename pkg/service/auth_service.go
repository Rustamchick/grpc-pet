package Service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"grpc-pet/pkg/models"
	"grpc-pet/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log         *logrus.Logger
	repos       repository.Authentification
	appProvider repository.AppProvider
	TokenTTL    time.Duration
}

func NewAuthService(log *logrus.Logger, repos repository.Authentification, AppProvider repository.AppProvider, TokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:         log,
		repos:       repos,
		appProvider: AppProvider,
		TokenTTL:    TokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string, appid int) (string, error) {
	const loc = "AuthService.Login()"

	log := s.log.WithField("loc", loc)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error hashing password: %s", err.Error())
		return "", err
	}

	user, err := s.repos.Login(ctx, email, hashedPassword, appid)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotExists) {
			log.Errorf("error login user: %s", err.Error())
			return "", ErrInvalidCredentials
		}
		log.Errorf("error login user: %s", err.Error())
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Infof("error CompareHashAndPassword: %s", err.Error())
		return "", err
	}

	app, err := s.appProvider.GetApp(ctx, appid)
	if err != nil {
		return "", err
	}

	token, err := s.GenerateToken(user, app)
	if err != nil {
		log.Info("Error generating token: %s", err.Error())
		return "", err
	}
	return token, nil
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int64, error) {
	const loc = "AuthService.Register()"

	log := s.log.WithField("loc", loc)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error hashing password: %s", err.Error())
		return 0, err
	}

	userid, err := s.repos.RegisterNewUser(ctx, email, hashedPassword)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			log.Errorf("error registering new user: %s", err.Error())
			return 0, ErrInvalidCredentials
		}
		log.Errorf("Error registering new user: %s", err.Error()) // TODO: handle varios errors
		return 0, err
	}

	return userid, nil
}

func (s *AuthService) IsAdmin(ctx context.Context, userid int) (bool, error) {
	const loc = "AuthService.IsAdmin()"

	log := s.log.WithField("loc", loc)

	log.Infof("Is Admin")

	return false, nil
}

// GenetateToken does not ready
func (s *AuthService) GenerateToken(user models.User, app models.App) (string, error) {
	const signingKey = "fewfdxsf32t4yt22saf4231r"

	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
		UserId: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	const salt = "r342fdsg8j8t4g20"

	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
