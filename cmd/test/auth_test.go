package main

import (
	"testing"
	"time"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	EmptyAppId = 0
	AppId      = 1
	AppToken   = "rest-token"
	timeDelta  = 1
)

func TestRegisterLogin(t *testing.T) {
	ctx, cl := NewTestClient(t)
	defer ctx.Done()

	email := gofakeit.Email()
	password := fakePassword()
	RegReq := &grpcpetv1.RegisterRequest{
		Email:    email,
		Password: password,
	}

	RegResp, err := cl.AuthClient.Register(ctx, RegReq)
	require.NoError(t, err)
	assert.NotEmpty(t, RegResp)

	LogReq := &grpcpetv1.LoginRequest{
		AppId:    AppId,
		Email:    email,
		Password: password,
	}

	LogResp, err := cl.AuthClient.Login(ctx, LogReq)
	require.NoError(t, err)
	assert.NotEmpty(t, LogResp)

	loginTime := time.Now()
	token := LogResp.GetToken()

	claims, err := parseToken(token, AppToken)
	require.NoError(t, err)
	assert.NotEmpty(t, claims)

	assert.Equal(t, RegResp.UserId, int64(claims.UserId))
	assert.Equal(t, RegReq.Email, (claims.Email))
	assert.Equal(t, LogReq.AppId, int32(claims.AppId))

	assert.InDelta(t, loginTime.Add(cl.Cfg.TokenTTL).Unix(), float64(claims.ExpiresAt.Time.Unix()), timeDelta)
}

func TestDoubleRegister(t *testing.T) {
	ctx, cl := NewTestClient(t)
	defer ctx.Done()

	email := gofakeit.Email()
	password := fakePassword()
	RegReq := &grpcpetv1.RegisterRequest{
		Email:    email,
		Password: password,
	}

	RegRes, err := cl.AuthClient.Register(ctx, RegReq)
	require.NoError(t, err)
	assert.NotEmpty(t, RegRes)

	RegRes2, err := cl.AuthClient.Register(ctx, RegReq)
	require.Error(t, err)
	assert.Empty(t, RegRes2.GetUserId())
	assert.ErrorContains(t, err, "User already exists")

}
func TestInvalidRegistration(t *testing.T) {
	ctx, cl := NewTestClient(t)
	defer ctx.Done()

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Empty EMAIL registration",
			email:       "",
			password:    fakePassword(),
			expectedErr: "Email is required",
		},
		{
			name:        "Empty PASSWORD registration",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "Password is required",
		},
		{
			name:        "Empty EMAIL and PASSWORD registration",
			email:       "",
			password:    "",
			expectedErr: "Email is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := cl.AuthClient.Register(ctx, &grpcpetv1.RegisterRequest{
				Email:    test.email,
				Password: test.password,
			})
			require.Error(t, err)
			// require.Contains(t, err.Error(), test.expectedErr)
			assert.ErrorContains(t, err, test.expectedErr)
		})
	}
}

func TestInvalidLogin(t *testing.T) {
	ctx, cl := NewTestClient(t)
	defer ctx.Done()

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Empty EMAIL login",
			email:       "",
			password:    fakePassword(),
			appID:       AppId,
			expectedErr: "Email is required",
		},
		{
			name:        "Empty PASSWORD login",
			email:       gofakeit.Email(),
			password:    "",
			appID:       AppId,
			expectedErr: "Password is required",
		},
		{
			name:        "Empty EMAIL and PASSWORD login",
			email:       "",
			password:    "",
			appID:       AppId,
			expectedErr: "Email is required",
		},
		{
			name:        "login with WRONG PASSWORD",
			email:       gofakeit.Email(),
			password:    fakePassword(),
			appID:       AppId,
			expectedErr: "Invalid email or password",
		},
		{
			name:        "Empty AppID login",
			email:       gofakeit.Email(),
			password:    fakePassword(),
			appID:       0,
			expectedErr: "AppId is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := cl.AuthClient.Register(ctx, &grpcpetv1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: fakePassword(),
			})
			require.NoError(t, err)

			_, err = cl.AuthClient.Login(ctx, &grpcpetv1.LoginRequest{
				Email:    test.email,
				Password: test.password,
				AppId:    test.appID,
			})
			require.Error(t, err)
			// require.Contains(t, err.Error(), test.expectedErr)
			assert.ErrorContains(t, err, test.expectedErr)
		})
	}

}

func TestRealInvalidLogin(t *testing.T) {
	ctx, cl := NewTestClient(t)

	tempEmail := gofakeit.Email()
	testReq := &grpcpetv1.RegisterRequest{
		Email:    tempEmail,
		Password: fakePassword(),
	}
	testLogin := &grpcpetv1.LoginRequest{
		Email:    tempEmail,
		Password: fakePassword(),
		AppId:    AppId,
	}

	_, err := cl.AuthClient.Register(ctx, testReq)
	require.NoError(t, err)

	_, err = cl.AuthClient.Login(ctx, testLogin)
	require.Error(t, err)
	assert.ErrorContains(t, err, "Invalid email or password")
}

func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 10)
}
