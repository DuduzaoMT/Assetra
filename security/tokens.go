package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func NewToken(userId string, roles []string) (string, error) {
       claims := jwt.MapClaims{
	       "exp":    time.Now().Add(15 * time.Minute).Unix(),
	       "iat":    time.Now().Unix(),
	       "nbf":    time.Now().Unix(),
	       "iss":    tokenIssuer,
	       "sub":    userId,
	       "roles":  roles,
	       "jti":    userId, // JWT ID for token revocation
       }
       token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
       return token.SignedString(securityJWTkey)
}

func NewAccessToken(UserId string, roles []string) (string,error){
	return NewToken(UserId, roles);
}

func parseJwtCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return securityJWTkey, nil
}

func ExtractToken(r *http.Request) (string, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidToken
	}
	
	// Expected format: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", ErrInvalidToken
	}
	
	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", ErrInvalidToken
	}
	
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", ErrInvalidToken
	}
	
	return token, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, parseJwtCallback)
}

type TokenPayload struct {
	UserId    string
	Roles     []string
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

	// Extract roles (JWT unmarshals arrays as []interface{})
	var roles []string
	if rolesRaw, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rolesRaw {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}
	}

	createdAt, _ := claims["iat"].(float64)
	expiresAt, _ := claims["exp"].(float64)

	return &TokenPayload{
		UserId:    userId,
		Roles:     roles,
		CreatedAt: time.Unix(int64(createdAt), 0),
		ExpiresAt: time.Unix(int64(expiresAt), 0),
	}, nil
}

// NewRefreshToken generates a cryptographically secure random refresh token
func NewRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken creates a SHA-256 hash of the refresh token
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SetRefreshTokenCookie sets the refresh token in an httpOnly cookie
func SetRefreshTokenCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})
}

// ExtractRefreshToken reads the refresh token from the cookie
func ExtractRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Println("error extracting refresh token from cookie:", err)
		return "", errors.New("refresh token not found")
	}
	if cookie.Value == "" {
		return "", errors.New("refresh token is empty")
	}
	return cookie.Value, nil
}

// ClearAuthCookies removes the refresh token cookie
func ClearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
