package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Company struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description,omitempty"`
	Employees   int    `json:"employees"`
	Registered  bool   `json:"registered"`
	Type        string `json:"type"`
}

var (
	Db *gorm.DB
)

func Connection() {
	dsn := "root:@tcp(localhost:3306)/assessment?parseTime=true"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Auto-migrate the company model
	err = Db.AutoMigrate(&Company{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
}

func CreateCompany(c *gin.Context) {
	var company Company

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company.ID = uuid.New().String()

	// Insert the company into the database
	result := Db.Create(&company)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func PatchCompany(c *gin.Context) {
	id := c.Param("id")
	var updatedCompany Company

	if err := c.ShouldBindJSON(&updatedCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the company by ID
	var foundCompany Company
	result := Db.First(&foundCompany, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	// Update the company
	foundCompany.Name = updatedCompany.Name
	foundCompany.Description = updatedCompany.Description
	foundCompany.Employees = updatedCompany.Employees
	foundCompany.Registered = updatedCompany.Registered
	foundCompany.Type = updatedCompany.Type

	// Save the changes to the database
	result = Db.Save(&foundCompany)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, foundCompany)
}

func DeleteCompany(c *gin.Context) {
	id := c.Param("id")

	// Delete the company from the database
	result := Db.Delete(&Company{}, "id = ?", id)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

func GetCompany(c *gin.Context) {
	id := c.Param("id")

	// Find the company by ID
	var company Company
	result := Db.First(&company, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}
