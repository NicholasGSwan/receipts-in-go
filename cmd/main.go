package main

import (
	"NicholasGSwan/receipts-in-go/internal/models"
	"NicholasGSwan/receipts-in-go/internal/pointsservice"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db = make(map[string]string)

var Database, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.POST("/receipt", func(c *gin.Context) {
		var receipt models.Receipt

		if err := c.Bind(&receipt); err == nil {
			Database.Create(&receipt)
			c.JSON(http.StatusOK, gin.H{"id": receipt.ID})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.GET("/receipt/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")

		var receipt models.Receipt
		result := Database.Model(&models.Receipt{}).Preload("Items").First(&receipt, id)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"receipt": receipt})
		}

	})

	r.GET("/receipt/:id/points", func(c *gin.Context) {
		id := c.Params.ByName("id")

		var receipt models.Receipt
		result := Database.Model(&models.Receipt{}).Preload("Items").First(&receipt, id)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		} else {
			totalPoints := pointsservice.CalcPoints(receipt)
			c.JSON(http.StatusOK, gin.H{"totalPoints": totalPoints})
			//points service -> calc points
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	if err != nil {
		panic(err.Error())
	}
	Database.AutoMigrate(&models.Receipt{})
	Database.AutoMigrate(&models.Item{})
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
