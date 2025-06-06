package tests

import (
	"grpc-service-ref/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/nonam00/protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
    secret = "test-secret"
    passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
    ctx, st := suite.New(t)

    email := gofakeit.Email()
    pass := randomFakePassword()

    respRegister, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
        Email: email,
        Password: pass,
    })

    require.NoError(t, err)
    assert.NotEmpty(t, respRegister.GetUserId())

    respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
  	    Email:    email,
  	    Password: pass,
    })

    loginTime := time.Now()

    require.NoError(t, err)

    token := respLogin.GetToken()
    require.NotEmpty(t, token)
    
    // TODO: secret
    tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
        return []byte(secret), nil
    })
    require.NoError(t, err)

    claims, ok := tokenParsed.Claims.(jwt.MapClaims)
    assert.True(t, ok)
    
    assert.Equal(t, respRegister.GetUserId(), int64(claims["uid"].(float64)))

    const deltaSeconds = 1
    assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
    ctx, st := suite.New(t)
    
    email := gofakeit.Email()
    pass := randomFakePassword()

    respRegister1, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
        Email: email,
        Password: pass,
    })

    require.NoError(t, err)
    require.NotEmpty(t, respRegister1.GetUserId())
    
    respRegister2, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
        Email: email,
        Password: pass,
    })
    require.Error(t, err)
    assert.Empty(t, respRegister2.GetUserId())
    assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
    ctx, st := suite.New(t)

    tests := []struct {
        name        string
        email       string
        password    string
        expectedErr string
    } {
        {
			      name:        "Register with Empty Password",
			      email:       gofakeit.Email(),
			      password:    "",
			      expectedErr: "password is required",
		    },
		    {
			      name:        "Register with Empty Email",
			      email:       "",
			      password:    randomFakePassword(),
			      expectedErr: "email is required",
		    },
		    {
			      name:        "Register with Both Empty",
			      email:       "",
			      password:    "",
			      expectedErr: "email is required",
	      },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            _, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
                Email: tt.email,
                Password: tt.password,
            })
            require.Error(t, err)
            require.Contains(t, err.Error(), tt.expectedErr)
        })
    }
}

func TestLogin_FailCases(t *testing.T) {
	  ctx, st := suite.New(t)

	  tests := []struct {
		    name        string
		    email       string
		    password    string
		    appID       int32
		    expectedErr string
	  } {
		    {
			      name:        "Login with Empty Password",
			      email:       gofakeit.Email(),
			      password:    "",
			      expectedErr: "password is required",
		    },
		    {
			      name:        "Login with Empty Email",
			      email:       "",
			      password:    randomFakePassword(),
			      expectedErr: "email is required",
		    },
		    {
			      name:        "Login with Both Empty Email and Password",
			      email:       "",
			      password:    "",
			      expectedErr: "email is required",
		    },
		    {
			      name:        "Login with Non-Matching Password",
			      email:       gofakeit.Email(),
			      password:    randomFakePassword(),
			      expectedErr: "invalid email or password",
		    },
    }

	  for _, tt := range tests {
        tt := tt
		    t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
			      _, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				        Email:    gofakeit.Email(),
				        Password: randomFakePassword(),
			      })
			      require.NoError(t, err)

			      _, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				        Email:    tt.email,
				        Password: tt.password,
			      })
		      	require.Error(t, err)
			      require.Contains(t, err.Error(), tt.expectedErr)
		    })
	  }
}

func randomFakePassword() string {
    return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
