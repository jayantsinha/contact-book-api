package main

import (
	"contact-book-api/model"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

// AuthHeaderMiddleware checks for Authorization header and sets account ID in request scope
func AuthHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			respondWithError(http.StatusUnauthorized, "Authorization Required", c)
			return
		}
		authval := strings.SplitN(token, " ", 2)
		if len(authval) != 2 || authval[0] != "Basic" {
			respondWithError(http.StatusUnauthorized, "Authorization Required", c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(authval[1])
		userpasspair := strings.Split(string(payload), ":")
		log.Println(base64.StdEncoding.EncodeToString([]byte(authval[1])))
		accountID := validateAccount(userpasspair[0], userpasspair[1])
		if accountID == -1 {
			respondWithError(http.StatusUnauthorized, "Invalid Credentials", c)
			return
		}
		c.Set("AccID", accountID)
		c.Next()
	}
}

// validateAccount validates email and password string and returns the account ID on success or -1 on failure
func validateAccount(email, password string) int {
	var account []model.Account
	DB.Select(&account, "SELECT * FROM `accounts` WHERE `email` = ?", email)
	log.Println(email)
	// validate password
	err := bcrypt.CompareHashAndPassword([]byte(account[0].Password), []byte(password))
	if err != nil {
		return -1
	}
	return account[0].AccountID
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}
