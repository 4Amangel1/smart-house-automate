package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router     *gin.Engine
	repo       *database.Repository
	logger     *log.Logger
	httpServer *http.Server
	config     config.APIConfig
}

func NewServer(repo *database.Repository, logger *log.Logger, cfg config.APIConfig) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	server := &Server{
		router: r,
		repo:   repo,
		logger: logger,
		config: cfg,
	}

	server.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(metricsMiddleware())

	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := s.router.Group("/api/v1")
	{
		api.GET("/sensors", s.getAllSensors)
		api.GET("/readings/latest", s.getLatestReadings)
		api.GET("/readings/latest/:type", s.getLatestReadingsByType)
		api.GET("/readings/history", s.getReadingHistory)
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) getAllSensors(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sensors": []string{"temp_1", "motion_1", "air_1"},
	})
}

func (s *Server) getLatestReadings(c *gin.Context) {
	readings, err := s.repo.GetLatestReadings()
	if err != nil {
		s.logger.Printf("Error getting latest readings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get readings"})
		return
	}

	c.JSON(http.StatusOK, readings)
}

func (s *Server) getLatestReadingsByType(c *gin.Context) {
	sensorType := c.Param("type")

	if sensorType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor type is required"})
		return
	}

	readings, err := s.repo.GetLatestReadingsByType(sensorType)
	if err != nil {
		s.logger.Printf("Error getting latest readings by type %s: %v", sensorType, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get readings"})
		return
	}

	c.JSON(http.StatusOK, readings)
}

func (s *Server) getReadingHistory(c *gin.Context) {
	sensorID := c.Query("sensorId")
	if sensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	readings, err := s.repo.GetReadingHistory(sensorID, limit)
	if err != nil {
		s.logger.Printf("Error getting reading history for sensor %s: %v", sensorID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reading history"})
		return
	}

	c.JSON(http.StatusOK, readings)
}
