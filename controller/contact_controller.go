package controller

import (
	"contact-book-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

type PostRequestBody struct {
	FirstName string
	LastName  string
	Email     string
	PhoneNum  string
}

// GetContacts returns a list of contacts by page number specified.
// Note that the result-set per page is set to 10 by default
func GetContacts(c *gin.Context) {
	accID := c.MustGet("AccID").(int)
	db := c.MustGet("DB").(*sqlx.DB)
	var contacts []model.Contact
	//page, err := strconv.ParseInt(c.Param("page"), 10, 16)	// result can fit inside int32 which will accomodate 3,276,700 records
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invalid page number! Page number should be greater than 0."})
		return
	}
	if page < 1 {
		page = 1
	}

	db.Select(&contacts, "SELECT * FROM contacts WHERE `account_id` = ? ORDER BY `first_name` ASC LIMIT ?, 10", accID, (page*10)-10)

	if contacts == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func DeleteContact(c *gin.Context) {

}

func CreateContact(c *gin.Context) {

}

func EditContact(c *gin.Context) {

}
