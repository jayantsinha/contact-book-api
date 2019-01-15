package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostRequestBody struct {
	FirstName string
	LastName  string
	Email     string
	PhoneNum  string
}

// GetContacts returns a list of contacts by page number specified
func GetContacts(c *gin.Context) {
	accID := c.MustGet("AccID").(int)


	c.JSON(http.StatusOK, gin.H{"message": accID})
}

func DeleteContact(c *gin.Context) {

}

func CreateContact(c *gin.Context) {

}

func EditContact(c *gin.Context) {

}
