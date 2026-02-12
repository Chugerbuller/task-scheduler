package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func signIn(w http.ResponseWriter, r *http.Request) {
	var req passwordReq
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	password := os.Getenv("TODO_PASSWORD")
	if len(password) == 0 {
		writeError(w, fmt.Errorf("Can`t load env"), http.StatusInternalServerError)
		return
	}
	if password == req.Password {
		token, err := createToken(password)
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
	jwtSecret := []byte(os.Getenv("TODO_PASSWORD"))
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

	currentPasswordHash := getCurrentPasswordHash()
	if claims.PasswordHash != currentPasswordHash {
		return false, fmt.Errorf("The password has been changed ")
	}
	return true, nil
}
func createToken(password string) (string, error) {
	jwtSecret := []byte(os.Getenv("TODO_PASSWORD"))
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("Can`t load env")
	}
	passHash := getCurrentPasswordHash()
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
func getCurrentPasswordHash() string {
	password := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
