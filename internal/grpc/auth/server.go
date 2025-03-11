package auth

import (
	"context"
	"errors"
	"fmt"
	"grpc-service-ref/internal/services/auth"

	ssov1 "github.com/nonam00/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Auth interface {
    Login (ctx context.Context,
        email string,
        password string,
    ) (token string, err error)
    RegisterNewUser(ctx context.Context,
        email string,
        password string,
    ) (userID int64, err error)
}

type serverAPI struct {
    ssov1.UnimplementedAuthServer
    auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
    ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
    emptyValue = 0
)

func (s *serverAPI) Login(
    ctx context.Context,
    req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
    if err := validateLogin(req); err != nil {
        return nil, err
    }

    token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())

    if err != nil {
        if errors.Is(err, auth.ErrInvalidCredentials) {
            return nil, status.Error(codes.InvalidArgument, "invalid credentials")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    // TODO: token cookie
    header := metadata.Pairs("set-cookie", fmt.Sprintf("token=%s; httponly; secure; samesite=none; maxage=3600", token))
    grpc.SendHeader(ctx, header)

    return &ssov1.LoginResponse{
        Token: token,
    }, nil
}

func (s *serverAPI) Register(
    ctx context.Context,
    req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
    if err := validateRegiter(req); err != nil {
        return nil, err
    }

    userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
    if err != nil {
        if errors.Is(err, auth.ErrUserExists) {
            return nil, status.Error(codes.AlreadyExists, "user already exists")
        }

        return nil, status.Error(codes.Internal, "internal error")
    }

    return &ssov1.RegisterResponse{
        UserId: userID,
    }, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }

    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }

    return nil
}

func validateRegiter(req *ssov1.RegisterRequest) error {
    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }

    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }
  
    return nil
}
