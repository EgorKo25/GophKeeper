package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/storage"

	"github.com/dgrijalva/jwt-go"
)

// Auth is a struct for authentication and session control
type Auth struct {
	secret string
}

// NewAuth is a contractor
func NewAuth(secret string) *Auth {
	return &Auth{
		secret: secret,
	}
}

var (
	ErrTokenInvalid    = errors.New("token invalid")
	ErrClaimsNotOfType = errors.New("token claims are not of type *tokenClaims")
	ErrSigningMethod   = errors.New("invalid singing method")
)

type claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func (a *Auth) RefreshTokens(access, refresh, login string) ([]*http.Cookie, error) {
	accessParsed, err := jwt.ParseWithClaims(access, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	refreshParsed, err := jwt.ParseWithClaims(refresh, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if accessParsed.Valid && refreshParsed.Valid {
		return a.GenerateTokensAndCreateCookie(&storage.User{Login: login})
	}

	return nil, ErrTokenInvalid

}

func (a *Auth) ParseWithClaims(token string) (string, error) {
	tokenParsed, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return []byte(a.secret), nil
	})
	if err != nil {
		return "", err
	}

	claimsParsed, ok := tokenParsed.Claims.(*claims)
	if !ok {
		return "", ErrClaimsNotOfType
	}

	if !tokenParsed.Valid {
		return "", ErrTokenInvalid
	}

	return claimsParsed.Name, nil

}

func (a *Auth) GenerateTokensAndCreateCookie(user *storage.User) ([]*http.Cookie, error) {

	var accessToken, refreshToken string
	var err error

	accessToken, err = a.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err = a.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	exp := time.Now()

	cookies := make([]*http.Cookie, 3)
	cookies[0] = a.getUserCookie(user, exp)
	cookies[1] = a.getAccessCookie(accessToken, exp)
	cookies[2] = a.getRefreshCookie(refreshToken, exp)

	return cookies, nil

}

// getRefreshCookie generate access cookie
func (a *Auth) getRefreshCookie(token string, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "Refresh-token",
		Value:   token,
		Expires: exp,
	}
}

// getAccessCookie generate access cookie
func (a *Auth) getAccessCookie(token string, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "Accesses-token",
		Value:   token,
		Expires: exp,
	}
}

// getUserCookie generate user cookie
func (a *Auth) getUserCookie(user *storage.User, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "User",
		Value:   user.Login,
		Expires: exp,
	}
}

// generateToken generate user's jwt token
func (a *Auth) generateToken(user *storage.User, exp time.Time) (string, error) {
	cl := &claims{
		Name: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateAccessToken generate user's jwt token
func (a *Auth) generateAccessToken(user *storage.User) (string, error) {
	exp := time.Now().Add(1 * time.Hour)

	return a.generateToken(user, exp)
}

// generateRefreshToken generate user's jwt token
func (a *Auth) generateRefreshToken(user *storage.User) (string, error) {

	exp := time.Now().Add(24 * time.Hour)

	return a.generateToken(user, exp)
}
