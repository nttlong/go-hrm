package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"vn.ghrm/internal/config"
	"vn.ghrm/internal/db"
	"vn.ghrm/internal/models"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up the database and user
	if err := db.SetupDatabase(cfg); err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}

	// Initialize database connection with GORM
	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Set up the server
	server := &Server{
		DB:     db,
		Router: gin.Default(),
	}

	// Define routes
	server.setupRoutes()

	// Start the server
	addr := cfg.AppPort
	log.Printf("Server starting on %s", addr)
	if err := server.Router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) setupRoutes() {
	s.Router.GET("/api/employees", s.getEmployees)
	s.Router.POST("/api/employees", s.createEmployee)
}

func (s *Server) getEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := s.DB.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (s *Server) createEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, employee)
}
