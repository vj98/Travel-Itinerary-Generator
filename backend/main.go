package main

import (
	"backend/internal/database"
	"backend/middleware"
	"backend/pkg/api"
	"backend/pkg/store"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, session")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	dbClient = database.ConnectMongoDB()
	IndexEmailId(dbClient)

	// Auth
	userStore := store.NewUserStore(dbClient, "travel", "users")
	handler := api.NewHandler(userStore)

	router.POST("/api/register", handler.Register)
	router.POST("/api/login", handler.Login)
	router.POST("/api/refresh", handler.RefreshSession)
	router.GET("/api/token", api.VerifyToken)

	protected := router.Group("/user").Use(ReadBodyIntoContext()).Use(middleware.JWTAuthMiddleware(dbClient))
	{
		protected.POST("/fetch", handler.FetchUser)
		protected.POST("/search", api.CompletionHandler)
		protected.POST("/dummy", api.DummyCall)
		protected.POST("/all", handler.FetchAllUsers)
		protected.POST("/update/status", handler.UpdateStatus)
		protected.POST("/contact", handler.SendEmail)
	}

	// 	// Start the server
	router.Run(":8080")

	defer dbClient.Disconnect(context.TODO())
}

func IndexEmailId(dbClient *mongo.Client) {
	db := dbClient.Database("travel")
	collection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created index:", indexName)
}

var ginLambda *ginadapter.GinLambda

// func init() {
// 	router := gin.New()

// 	router.Use(gin.Logger())

// 	router.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pong",
// 		})
// 	})

// 	router.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	})

// 	dbClient := database.ConnectMongoDB()
// 	IndexEmailId(dbClient)
// 	userStore := store.NewUserStore(dbClient, "travel", "users")
// 	handler := api.NewHandler(userStore)

// 	router.POST("/api/register", handler.Register)
// 	router.POST("/api/login", handler.Login)
// 	router.POST("/api/refresh", handler.RefreshSession)
// 	router.GET("/api/token", api.VerifyToken)

// 	protected := router.Group("/user").Use(ReadBodyIntoContext()).Use(middleware.JWTAuthMiddleware(dbClient))
// 	{
// 		protected.POST("/fetch", handler.FetchUser)
// 		protected.POST("/search", api.CompletionHandler)
// 		protected.POST("/dummy", api.DummyCall)
// 		protected.POST("/all", handler.FetchAllUsers)
// 		protected.POST("/update/status", handler.UpdateStatus)
// 		protected.POST("/contact", handler.SendEmail)
// 	}

// 	// Initialize ginadapter
// 	ginLambda = ginadapter.New(router)
// }

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log the incoming request (be cautious of logging sensitive information in production)
	log.Printf("Received request: %+v\n", req)

	// Use the adapter to handle the request with Gin
	resp, err := ginLambda.ProxyWithContext(ctx, req)
	if err != nil {
		log.Printf("Error handling request: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	// Log the response (again, be cautious in production)
	log.Printf("Generated response: %+v\n", resp)

	return resp, nil
}

// func main() {
// 	lambda.Start(Handler)
// }

func ReadBodyIntoContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Unmarshal using json.Unmarshal
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request body is not valid JSON"})
			c.Abort()
			return
		}

		// Store the body in the context for later use in handlers
		c.Set("body", body)

		// Call the next middleware/handler in chain
		c.Next()
	}
}
