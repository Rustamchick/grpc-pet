package AuthService

import (
	"context"
	"errors"
	"grpc-pet/pkg/repository"
	AppPostgres "grpc-pet/pkg/repository/postgres/app"
	AuthPostgres "grpc-pet/pkg/repository/postgres/auth"

	"time"

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

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserExists         = errors.New("User already exists")
	ErrInvalidAppID       = errors.New("Invalid app id")
	ErrUserNotFound       = errors.New("User not found")
)

func (s *AuthService) Login(ctx context.Context, email, password string, appid int) (string, error) {
	const loc = "AuthService.Login()"

	log := s.log.WithField("loc", loc)

	user, err := s.repos.Login(ctx, email)
	if err != nil {
		if errors.Is(err, AuthPostgres.ErrUserNotExists) {
			log.Errorf("Error. User not exist: %s", err.Error())
			return "", ErrInvalidCredentials // email is not correct
		}
		log.Errorf("error login user: %s", err.Error())
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Infof("error CompareHashAndPassword: %s", err.Error())
		return "", ErrInvalidCredentials // email correct, incorrect password
	}

	app, err := s.appProvider.GetApp(ctx, appid)
	if err != nil {
		if errors.Is(err, AppPostgres.ErrAppNotFound) {
			log.Errorf("Error. App not found: %s", err.Error())
			return "", ErrInvalidAppID // appid is not correct
		}
		return "", err
	}

	token, err := GenerateToken(user, app, s.TokenTTL)
	if err != nil {
		log.Infof("Error generating token: %s", err.Error())
		return "", err
	}

	log.Info("User logged in successfully")

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
		if errors.Is(err, AuthPostgres.ErrUserExists) {
			log.Errorf("error. User already exist: %s", err.Error())
			return 0, ErrUserExists
		}
		log.Errorf("Error registering new user: %s", err.Error())
		return 0, err
	}

	log.Info("User registered successfully")

	return userid, nil
}

func (s *AuthService) IsAdmin(ctx context.Context, userid int) (bool, error) {
	const loc = "AuthService.IsAdmin()"

	log := s.log.WithField("loc", loc)
	log.Infof("Check if user is admin. userId: %d", userid)

	isAdmin, err := s.repos.IsAdmin(ctx, userid)
	if err != nil {
		if errors.Is(err, AuthPostgres.ErrUserNotExists) {
			log.Infof("Error. App not found.")
			return isAdmin, ErrUserNotFound
		}
		log.Info("Error. IsAdmin. Undefined error.")
		return false, err
	}

	return isAdmin, nil
}
