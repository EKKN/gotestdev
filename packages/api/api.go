package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/EKKN/gotestdev/packages/jobid"
	"github.com/EKKN/gotestdev/packages/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type APIServer struct {
	ListenAddr string
}

func NewAPI(ListenAddr string) *APIServer {
	return &APIServer{
		ListenAddr: ListenAddr,
	}
}

func (s *APIServer) Run() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	URL_ALLOWORIGIN := os.Getenv("ALLOW_ORIGIN")
	TRUSTED_PROXY := os.Getenv("TRUSTED_PROXY")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{URL_ALLOWORIGIN},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Atur proxy hanya untuk localhost
	router.SetTrustedProxies([]string{TRUSTED_PROXY})

	router.GET("/marketing", s.handleResponse(s.GetAllMarketing))
	router.POST("/marketing", s.handleResponse(s.CreateMarketing))
	router.PUT("/marketing/:id/toggle", s.handleResponse(s.ToggleMarketingActive))

	router.GET("/komisi", s.handleResponse(s.GetAllKomisi))
	router.GET("/laporan-komisi", s.handleResponse(s.GetLaporanKomisi))

	router.GET("/penjualan", s.handleResponse(s.GetAllPenjualan))
	router.POST("/penjualan", s.handleResponse(s.CreatePenjualan))

	router.GET("/kredit", s.handleResponse(s.GetAllKredit))

	// GetAllPenjualanNotPembayaran
	router.GET("/pembayaran/list-penjualan", s.handleResponse(s.GetAllPenjualanNotPembayaran))
	router.GET("/pembayaran", s.handleResponse(s.GetAllPembayaran))
	router.POST("/pembayaran", s.handleResponse(s.CreatePembayaran))
	router.GET("/pembayaran-details/:pembayaranid", s.handleResponse(s.GetAllPembayaranDetailByPembayaranId))

	log.Fatal(router.Run(s.ListenAddr))
}

func (s *APIServer) filterResponseLog(response map[string]interface{}) map[string]interface{} {
	filteredResponse := make(map[string]interface{})

	for k, v := range response {
		if k != "token" && k != "password" && k != "pwd" {
			filteredResponse[k] = v
		}
	}
	return filteredResponse
}

func (s *APIServer) filterResponse(jobId string, isError bool, response map[string]interface{}) map[string]interface{} {
	filteredResponse := make(map[string]interface{})

	var status string
	if isError {
		status = "error"
	} else {
		status = "success"
	}
	filteredResponse["status"] = status
	filteredResponse["jobid"] = jobId

	for k, v := range response {
		if k == "password" || k == "pwd" || k == "token" || k == "error" {
			continue
		}

		// Check if the value is a nested JSON object
		if reflect.TypeOf(v).Kind() == reflect.Map {
			nestedResponse := v.(map[string]interface{})
			filteredResponse[k] = s.filterResponse(jobId, isError, nestedResponse)
		} else {
			filteredResponse[k] = v
		}
	}

	return filteredResponse
}

func (s *APIServer) handleResponse(handlerFunc func(*gin.Context) (map[string]interface{}, map[string]interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			c.JSON(http.StatusOK, gin.H{"error": "Failed to read request body"})
			return
		}
		// Restore the request body to ensure it's not consumed
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		t := time.Now()

		data, handlerErr := handlerFunc(c)

		latency := time.Since(t)
		jobId := jobid.JobID()

		responseError := s.filterResponse(jobId, true, handlerErr)
		responseSuccess := s.filterResponse(jobId, false, data)

		var responseLog map[string]interface{}
		var response map[string]interface{}
		if handlerErr != nil {
			response = responseError
			responseLog = handlerErr
		} else {
			response = responseSuccess
			responseLog = data
		}

		s.requestResponseLog(c, jobId, responseLog, latency, bodyBytes)

		if handlerErr != nil {
			// c.JSON(http.StatusBadRequest, response)
			c.JSON(http.StatusOK, response)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func (s *APIServer) requestResponseLog(c *gin.Context, jobId string, responseLog map[string]interface{}, latency time.Duration, bodyBytes []byte) {
	clientIP := c.ClientIP()
	method := c.Request.Method
	url := c.Request.URL.String()
	userAgent := c.Request.UserAgent()
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")

	// Process and sanitize request body
	bodyRequestString := string(bodyBytes)
	if bodyRequestString == "" {
		bodyRequestString = "{}" // Set to an empty JSON object if the body is empty
	}

	// Unmarshal JSON to map
	bodyRequestMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bodyRequestString), &bodyRequestMap); err == nil {
		for k := range bodyRequestMap {
			if k == "password" || k == "pwd" || k == "token" {
				bodyRequestMap[k] = "***"
			}
		}
	} else {
		fmt.Printf("Failed to unmarshal request body: %v\n", err)
	}

	// Mask Authorization header
	headers := c.Request.Header
	for k := range headers {
		if k == "Authorization" {
			headers.Set(k, "***")
		}
	}

	// Create request log
	requestLog := map[string]interface{}{
		"body":      bodyRequestMap, // Use the map directly
		"method":    method,
		"url":       url,
		"headers":   headers,
		"client_ip": clientIP,
		"time":      currentTime,
		"agent":     userAgent,
	}

	// Combine request and response logs
	requestResponseLog := map[string]interface{}{
		"request":  requestLog,
		"response": s.filterResponseLog(responseLog),
		"latency":  latency.String(),
		"jobid":    jobId,
	}

	// Convert log to JSON
	logJSON, jsonErr := json.Marshal(requestResponseLog)
	if jsonErr != nil {
		fmt.Printf("Failed to encode log to JSON: %v\n", jsonErr)
		return
	}

	// Write log using the logger package
	logger.Log(string(logJSON))
}

func SetJobId() gin.HandlerFunc {
	return func(c *gin.Context) {
		jobId := jobid.JobID()
		c.Header("JobID", jobId)

		c.Next()
	}
}
