package controller

import (
	"contact-book-api/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strconv"
)

type PostRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// GetContacts returns a list of contacts by page number specified.
// Handler for [GET] /contacts/page/:page
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

	selErr := db.Select(&contacts, "SELECT * FROM `contacts` WHERE `account_id` = ? ORDER BY `first_name` ASC LIMIT ?, 10", accID, (page*10)-10)

	if selErr != nil {
		log.Println("Error retrieving contacts | ", selErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving contacts."})
	}

	if contacts == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// DeleteContact deletes the contact by contact ID
// Handler for [DELETE] /contact/:id
func DeleteContactByID(c *gin.Context) {
	accID := c.MustGet("AccID").(int)
	db := c.MustGet("DB").(*sqlx.DB)
	contactID := c.Param("id")
	res, err := db.Exec("DELETE FROM `contacts` WHERE `account_id` = ? AND `contact_id` = ?", accID, contactID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to delete! Might be due to bad ID"})
		return
	}
	count, err := res.RowsAffected()
	if err != nil || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to delete!"})
		return
	}
	// not possible by DB design (MySQL)
	if count > 1 {
		log.Println("Unexpected number of records removed. AccID = ")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected number of records removed!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully."})
}

// CreateContact creates a new contact entity under the given account ID
// Handler for [POST] /contact
func CreateContact(c *gin.Context) {
	accID := c.MustGet("AccID").(int)
	db := c.MustGet("DB").(*sqlx.DB)
	reqBody := &PostRequestBody{}
	bindErr := c.BindJSON(reqBody)
	if bindErr != nil {
		log.Println("Invalid request JSON ", reqBody)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	res, err := db.Exec("INSERT INTO `contacts` (`account_id`, `first_name`, `last_name`, `email`) VALUES (?, ?, ?, ?)", accID, reqBody.FirstName, reqBody.LastName, reqBody.Email)
	if err != nil {
		log.Println("Duplicate email address in request: ", reqBody, " | ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate email address! Unable to add contact."})
		return
	}
	id, err := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Contact added successfully with ID: " + strconv.FormatInt(id, 10)})
}

// EditContactByID edits the contact by given contact ID
// Handler for [PUT] /contact/:id
// This endpoint can be used to update the email address of the contact
func EditContactByID(c *gin.Context) {
	db := c.MustGet("DB").(*sqlx.DB)
	reqBody := &PostRequestBody{}
	bindErr := c.BindJSON(reqBody)
	if bindErr != nil {
		log.Println("Invalid request JSON ", reqBody)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	contactID := c.Param("id")
	res, err := db.Exec("UPDATE `contacts` SET `email` = ?, `first_name` = ?, `last_name` = ? WHERE `contact_id` = ?", reqBody.Email, reqBody.FirstName, reqBody.LastName, contactID)
	if err != nil {
		log.Println("Duplicate email address in request or invalid contact ID | ", reqBody, " | ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate email address or invalid contact ID."})
		return
	}
	// This condition should never pass by design
	if rowsAffected, _ := res.RowsAffected(); rowsAffected > 1 {
		log.Println("Multiple rows updated by request: ", reqBody)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Multiple records updated!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contact updated successfully."})
}

// EditContactByEmail edits the contact by email address matching
// Handler for [PUT] /contact
func EditContactByEmail(c *gin.Context) {
	db := c.MustGet("DB").(*sqlx.DB)
	reqBody := &PostRequestBody{}
	bindErr := c.BindJSON(reqBody)
	if bindErr != nil {
		log.Println("Invalid request JSON ", reqBody)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	res, err := db.Exec("UPDATE `contacts` SET `first_name` = ?, `last_name` = ? WHERE `email` = ?, ", reqBody.FirstName, reqBody.LastName, reqBody.Email)
	if err != nil {
		log.Println("Invalid email address in request | ", reqBody, " | ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email address in request!"})
		return
	}
	// This condition should never pass by design
	if rowsAffected, _ := res.RowsAffected(); rowsAffected > 1 {
		log.Println("Multiple rows updated by request: ", reqBody)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Multiple records updated!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contact updated successfully."})
}

// SearchContact returns the list of contact which matches the search criteria
func SearchContact(c *gin.Context) {
	
}