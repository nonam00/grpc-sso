package auth

import (
	"context"
	"errors"
	"fmt"
	"grpc-service-ref/internal/domain/models"
	"grpc-service-ref/internal/lib/jwt"
	"grpc-service-ref/internal/storage"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	usrSaver    UserSaver
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrProvider: userProvider,
		usrSaver:    userSaver,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login chechs if user with given credentials exists in the system.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error
func (a *Auth) Login(
    ctx context.Context,
	  email string,
	  password string,
	  appId int,
) (string, error) {
    const op = "auth.Login"

	  log := a.log.With(
		    slog.String("op", op),
		    slog.String("username", email),
	  )

	  log.Info("attempting to login user")

	  user, err := a.usrProvider.User(ctx, email)
	  if err != nil {

		    if errors.Is(err, storage.ErrUserNotFound) {
			      a.log.Warn("user not found", slog.String("err", err.Error()))

			      return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		    }
		    a.log.Error("failed to get user", slog.String("err", err.Error()))

		    return "", fmt.Errorf("%s: %w", op, err)
    }

	  if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		    a.log.Info("invalid credentials", slog.String("err", err.Error()))

		    return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	  }

	  app, err := a.appProvider.App(ctx, appId)
	  if err != nil {
		    return "", fmt.Errorf("%s: %w", op, err)
	  }

	  log.Info("user logged in successfully")

	  token, err := jwt.NewToken(user, app, a.tokenTTL)
	  if err != nil {
		    a.log.Error("failed to generate token", slog.String("err", err.Error()))

		  return "", fmt.Errorf("%s: %w", op, err)
	  }

	  return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(
	  ctx context.Context,
	  email string,
	  password string,
) (int64, error) {
	  const op = "Auth.RegisterNewUser"

	  log := a.log.With(
		    slog.String("op", op),
		    slog.String("email", email),
	  )

	  log.Info("registering user")

	  passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	  if err != nil {
		    log.Error("failed to generate password hash", slog.String("err", err.Error()))

		    return 0, fmt.Errorf("%s: %w", op, err)
	  }

	  id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	  if err != nil {
		    if errors.Is(err, storage.ErrUserExists) {
			      log.Warn("user already exists", slog.String("err", err.Error()))

		       	return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		    }

		    log.Error("failed to save user", slog.String("err", err.Error()))

		    return 0, fmt.Errorf("%s: %w", op, err)
	  }

	  log.Info("user registered")

	  return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	  const op = "Auth.IsAdmin"

	  log := a.log.With(
		    slog.String("op", op),
		    slog.Int64("user_id", userID),
	  )

	  log.Info("checking if user is admin")

	  isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	  if err != nil {
		    if errors.Is(err, storage.ErrUserNotFound) {
			      log.Warn("user not found", slog.String("err", err.Error()))

			      return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		    }
		    return false, fmt.Errorf("%s: %w", op, err)
	  }

	  log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	  return isAdmin, nil
}
