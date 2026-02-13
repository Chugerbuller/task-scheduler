package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	PasswordHash string `json:"password_hash"` // Хэш пароля в токене
}
type passwordReq struct {
	Password string `json:"password"`
}
type passwordResp struct {
	Token string `json:"token"`
}

var (
	password  string
	jwtSecret []byte
)

func signIn(w http.ResponseWriter, r *http.Request) {
	var req passwordReq
	if err := jsonDeserialize(r.Body, &req); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if len(password) == 0 {
		writeError(w, fmt.Errorf("Password wasn`t init"), http.StatusInternalServerError)
		return
	}
	if password == req.Password {
		token, err := createToken()
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		resp := passwordResp{
			Token: token,
		}
		writeJson(w, resp)
		return
	}
	writeError(w, fmt.Errorf("Invalid password"), http.StatusBadRequest)
}

func tokenIsValid(tokenStr string) (bool, error) {
	if len(jwtSecret) == 0 {
		return false, fmt.Errorf("Can`t load env")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return false, fmt.Errorf("Invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("Invalid token claims ")
	}

	passwordHash := getPasswordHash()
	if claims.PasswordHash != passwordHash {
		return false, fmt.Errorf("The password has been changed ")
	}
	return true, nil
}
func createToken() (string, error) {
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("Can`t load env")
	}
	passHash := getPasswordHash()
	claims := Claims{
		PasswordHash: passHash,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func getPasswordHash() string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
func InitPassword(pswd string) {
	password = pswd
}
func InitJwtSecret(jwtSec string) {
	jwtSecret = []byte(jwtSec)
}
