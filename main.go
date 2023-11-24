package main

import (
	"clubhub/config"
	"clubhub/internal"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/likexian/whois"
	"gorm.io/gorm"
)

type Franchise struct {
	gorm.Model
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Location   Location `json:"location" gorm:"foreignKey:LocationID"`
	LocationID uint
	Info       Info `json:"info" gorm:"foreignKey:InfoID"`
	InfoID     uint
	Company    Company `json:"company" gorm:"foreignKey:CompanyID"`
	CompanyID  uint
}

type Company struct {
	gorm.Model
	OwnerID    uint `json:"owner_id"`
	Franchises []Franchise
	ImageURL   string
	DomainInfo DomainInfo `gorm:"foreignKey:CompanyID;references:ID"`
}

type DomainInfo struct {
	gorm.Model
	CompanyID  uint      `json:"company_id"`
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
	OwnerName  string    `json:"owner_name"`
	OwnerEmail string    `json:"owner_email"`
}

type Location struct {
	ID      uint   `gorm:"primaryKey"`
	City    string `json:"city"`
	Country string `json:"country"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
}

type Owner struct {
	gorm.Model
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	LocationID uint   `json:"location_id"`
}
type Info struct {
	gorm.Model
	Name       string `json:"name"`
	TaxNumber  string `json:"tax_number"`
	LocationID uint   `json:"location_id"`
}

type SSLLabsResponse struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Endpoints []struct {
		IPAddress  string `json:"ipAddress"`
		ServerName string `json:"serverName"`
		Status     string `json:"status"`
		Grade      string `json:"grade"`
		Delegation int    `json:"delegation"`
	} `json:"endpoints"`
}

var db *gorm.DB

func main() {
	var err error
	port := ":3000"
	db, err = config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Migrar modelos a la base de datos
	db.AutoMigrate(&Franchise{}, &Location{}, &Company{}, &Owner{}, &Info{}, &DomainInfo{})

	// Iniciar el router
	router := gin.Default()

	// Configurar rutas
	router.POST("/franchise", createFranchise)
	router.PUT("/franchise/:id", updateFranchise)
	router.GET("/franchise/:id", getFranchiseByID)
	router.GET("/franchises", getAllFranchises)
	router.GET("/companies/:id", getFranchisesByCompany)
	router.GET("/location/:country", getFranchisesByLocation)
	router.GET("/ssl-info", getSSLInfo)
	router.GET("/domain-info", getDomainInfo)

	fmt.Println("Server running on port " + port)

	// Iniciar el servidor
	err = router.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

// Función para crear una nueva franquicia
func createFranchise(c *gin.Context) {
	var newFranchise Franchise
	if err := c.BindJSON(&newFranchise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normaliza la sección de ubicación
	newFranchise.Location.Country = internal.ToTitleCase(newFranchise.Location.Country)

	// Crear nueva franquicia en la base de datos
	db.Create(&newFranchise)

	c.JSON(http.StatusCreated, gin.H{"message": "Franchise created successfully"})
}

// Función para actualizar una franquicia existente
func updateFranchise(c *gin.Context) {
	id := c.Param("id")
	var updatedFranchise Franchise

	// Buscar la franquicia en la base de datos
	if err := db.First(&updatedFranchise, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}

	// Bind JSON a la franquicia existente
	if err := c.BindJSON(&updatedFranchise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normaliza la sección de ubicación
	updatedFranchise.Location.Country = internal.ToTitleCase(updatedFranchise.Location.Country)

	// Actualizar la franquicia en la base de datos
	db.Save(&updatedFranchise)

	c.JSON(http.StatusOK, gin.H{"message": "Franchise updated successfully"})
}

// Función para obtener una franquicia por ID
func getFranchiseByID(c *gin.Context) {
	id := c.Param("id")
	var franchise Franchise

	// Buscar la franquicia en la base de datos
	if err := db.Preload("Location").First(&franchise, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Franchise not found"})
		return
	}

	c.JSON(http.StatusOK, franchise)
}

// Función para obtener todas las franquicias
func getAllFranchises(c *gin.Context) {
	var franchises []Franchise

	// Obtener todas las franquicias de la base de datos
	db.Find(&franchises)

	c.JSON(http.StatusOK, franchises)
}

// Función para obtener todas las franquicias de una compañía
func getFranchisesByCompany(c *gin.Context) {
	id := c.Param("id")
	var company Company

	// Buscar la compañía en la base de datos
	if err := db.Preload("Franchises.Location").First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company.Franchises)
}

// Función para obtener todas las franquicias en un país
func getFranchisesByLocation(c *gin.Context) {
	country := c.Param("country")
	var matchingFranchises []Franchise

	// Obtener todas las franquicias de la base de datos que coinciden con el país
	db.Joins("JOIN locations ON franchises.location_id = locations.id").
		Where("country = ?", country).
		Find(&matchingFranchises)

	if len(matchingFranchises) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching franchises found"})
		return
	}

	c.JSON(http.StatusOK, matchingFranchises)
}

func getSSLInfo(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host parameter is required"})
		return
	}

	// Construye la URL para el servicio de SSL Labs
	url := "https://api.ssllabs.com/api/v3/analyze?host=" + host

	// Realiza la solicitud HTTP al servicio de SSL Labs
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get SSL info"})
		return
	}
	defer resp.Body.Close()

	// Decodifica la respuesta JSON
	var sslInfo SSLLabsResponse
	err = json.NewDecoder(resp.Body).Decode(&sslInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse SSL info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sslInfo": sslInfo})
}

func parseWhoisResponse(response string) (DomainInfo, error) {
	var domainInfo DomainInfo

	createdDate, err := internal.ExtractValue(response, "Creation Date:")
	if err != nil {
		return domainInfo, err
	}
	domainInfo.Created, err = time.Parse("2006-01-02", createdDate)
	if err != nil {
		return domainInfo, err
	}

	// Busca la cadena "Registry Expiry Date:" y extrae la fecha de expiración
	expiryDate, err := internal.ExtractValue(response, "Registry Expiry Date:")
	if err != nil {
		return domainInfo, err
	}
	domainInfo.Expires, err = time.Parse("2006-01-02", expiryDate)
	if err != nil {
		return domainInfo, err
	}

	// Busca la cadena "Registrant Name:" y extrae el nombre del registrado
	domainInfo.OwnerName, err = internal.ExtractValue(response, "Registrant Name:")
	if err != nil {
		return domainInfo, err
	}

	// Busca la cadena "Registrant Email:" y extrae el email del registrado
	domainInfo.OwnerEmail, err = internal.ExtractValue(response, "Registrant Email:")
	if err != nil {
		return domainInfo, err
	}

	return domainInfo, nil
}

func getDomainInfo(c *gin.Context) {
	domain := c.Query("host")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain parameter is required"})
		return
	}

	// Realiza la consulta de whois para obtener la información del dominio
	result, err := whois.Whois(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get domain info"})
		return
	}

	// Parsea la respuesta del whois para obtener la información deseada
	domainInfo, err := parseWhoisResponse(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse whois response"})
		return
	}

	c.JSON(http.StatusOK, domainInfo)
}
