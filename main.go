package main

import (
	"clubhub/internal"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Franchise struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Location Location `json:"location"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
}

type Company struct {
	Owner      Owner       `json:"owner"`
	Info       Info        `json:"information"`
	Franchises []Franchise `json:"franchises"`
}

type Owner struct {
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Location Location `json:"location"`
}

type Info struct {
	Name      string   `json:"name"`
	TaxNumber string   `json:"tax_number"`
	Location  Location `json:"location"`
}

var franchises map[string]Company

func main() {
	franchises = make(map[string]Company)
	router := gin.Default()

	router.POST("/franchise", createFranchise)
	router.PUT("/franchise/:name", updateFranchise)
	router.GET("/franchise/:name", getFranchiseByName)
	router.GET("/franchises", getAllFranchises)
	router.GET("/companies/:companyName", getFranchisesByCompany)
	router.GET("/location/:country", getFranchisesByLocation)

	err := router.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func createFranchise(c *gin.Context) {
	var newFranchise Franchise
	if err := c.BindJSON(&newFranchise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normaliza la sección de ubicación
	newFranchise.Location.Country = internal.ToTitleCase(newFranchise.Location.Country)

	// Agrega la nueva franquicia al mapa
	franchises[newFranchise.Name] = Company{Franchises: []Franchise{newFranchise}}

	c.JSON(http.StatusCreated, gin.H{"message": "Franchise created successfully"})
}

func updateFranchise(c *gin.Context) {
	name := c.Param("name")
	company, exists := franchises[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}

	var updatedFranchise Franchise
	if err := c.BindJSON(&updatedFranchise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normaliza la sección de ubicación
	updatedFranchise.Location.Country = internal.ToTitleCase(updatedFranchise.Location.Country)

	// Actualiza la información de la franquicia
	company.Franchises = append(company.Franchises, updatedFranchise)
	franchises[name] = company

	c.JSON(http.StatusOK, gin.H{"message": "Franchise updated successfully"})
}

// Función para obtener una franquicia por nombre
func getFranchiseByName(c *gin.Context) {
	name := c.Param("name")
	company, exists := franchises[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// Función para obtener todas las franquicias
func getAllFranchises(c *gin.Context) {
	c.JSON(http.StatusOK, franchises)
}

// Función para obtener todas las franquicias de una compañía
func getFranchisesByCompany(c *gin.Context) {
	companyName := c.Param("companyName")
	var matchingFranchises []Franchise

	for _, company := range franchises {
		if strings.Contains(strings.ToLower(company.Info.Name), strings.ToLower(companyName)) {
			matchingFranchises = append(matchingFranchises, company.Franchises...)
		}
	}

	if len(matchingFranchises) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching franchises found"})
		return
	}

	c.JSON(http.StatusOK, matchingFranchises)
}

// Función para obtener todas las franquicias en un país
func getFranchisesByLocation(c *gin.Context) {
	country := c.Param("country")
	var matchingFranchises []Franchise

	for _, company := range franchises {
		for _, franchise := range company.Franchises {
			if strings.EqualFold(franchise.Location.Country, country) {
				matchingFranchises = append(matchingFranchises, franchise)
			}
		}
	}

	if len(matchingFranchises) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching franchises found"})
		return
	}

	c.JSON(http.StatusOK, matchingFranchises)
}
