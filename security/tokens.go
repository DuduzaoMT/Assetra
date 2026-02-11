package security

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	securityJWTkey  []byte
	ErrInvalidToken = errors.New("invalid jwt")
	tokenIssuer     = "assetra"
)

func init() {
	err := godotenv.Load("../.env") // load .env file
	if err != nil {
		log.Panic("Failed to load .env file")
	}
	securityJWTkey = []byte(os.Getenv("SECURITY_KEY"))
	if len(securityJWTkey) < 32 {
		log.Panic("SECURITY_KEY must be at least 32 characters long")
	}
}

func NewToken(userId string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    tokenIssuer,
		Subject:   userId,
		ID:        userId, // jti - JWT ID for token revocation
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(securityJWTkey)
}

func NewAccessToken(UserId string) (string,error){
	return NewToken(UserId);
}

func parseJwtCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return securityJWTkey, nil
}

func ExtractToken(r *http.Request) (string, error) {
	// Extract token from access_token cookie
	cookie, err := r.Cookie("access_token")
	if err != nil {
		log.Println("error extracting token from cookie:", err)
		return "", ErrInvalidToken
	}
	if cookie.Value == "" {
		return "", ErrInvalidToken
	}
	return cookie.Value, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, parseJwtCallback)
}

type TokenPayload struct {
	UserId    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewTokenPayload(tokenString string) (*TokenPayload, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return nil, ErrInvalidToken
	}
	
	// Validate issuer
	issuer, ok := claims["iss"].(string)
	if !ok || issuer != tokenIssuer {
		return nil, fmt.Errorf("invalid token issuer")
	}
	
	// Extract subject (userId)
	userId, ok := claims["sub"].(string)
	if !ok || userId == "" {
		return nil, fmt.Errorf("invalid user ID in token")
	}
	
	createdAt, _ := claims["iat"].(float64)
	expiresAt, _ := claims["exp"].(float64)
	
	return &TokenPayload{
		UserId:    userId,
		CreatedAt: time.Unix(int64(createdAt), 0),
		ExpiresAt: time.Unix(int64(expiresAt), 0),
	}, nil
}

// SetAccessTokenCookie sets the access token in an httpOnly cookie
func SetAccessTokenCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   1800, // 30 minutes
		},
	)
}

// ClearAuthCookies removes the access token cookie
func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		},
	)
}
