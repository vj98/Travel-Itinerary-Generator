package api

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/pkg/model"
	"backend/pkg/store"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/resend/resend-go/v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Handler holds the endpoints for the user service.
type Handler struct {
	store *store.UserStore
}

// NewHandler creates a new Handler.
func NewHandler(userStore *store.UserStore) *Handler {
	return &Handler{
		store: userStore,
	}
}

// Login handles user login and sets an HTTP-only cookie with the JWT token.
func (h *Handler) Login(c *gin.Context) {
	var loginDetails model.Users
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		log.Println("issue in binding struct ", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	user, err := h.store.FindUserByEmail(c.Request.Context(), loginDetails.Email)
	log.Println("user ", user, err, h.store)
	if err != nil {
		log.Println("user not registerd ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "user not found"})
		return
	}

	if user.Status != "Active" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "user not active contact admin"})
		c.Abort()
		return
	}

	t := fmt.Sprintf("%v", time.Now().Local().Unix())
	err = h.store.UpdateSession(c.Request.Context(), loginDetails.Email, t)
	if err != nil {
		log.Println("error in session update ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error in session update"})
		return
	}

	match := checkPasswordHash(loginDetails.Password, user.Password)
	if !match {
		log.Println("password incorrect ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "password incorrect"})
		return
	}

	// Generate JWT token
	token, refreshToken, err := auth.GenerateTokens(user.Email)
	if err != nil {
		log.Println("issue in generating token ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "login successful", "access_token": token, "refresh_token": refreshToken, "session": t, "type": user.Type, "email": user.Email})
}

// RefreshSession handles the session refresh and updates the JWT token cookie.
func (h *Handler) RefreshSession(c *gin.Context) {
	// Extract the refresh token from the cookie
	refreshToken := c.GetHeader("Authorization")

	if refreshToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "refresh token is required"})
		return
	}

	// Refresh the access token
	newToken, err := auth.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session refreshed", "access_token": newToken})
}

// Register handles user registration.
func (h *Handler) Register(c *gin.Context) {
	var newUser model.Users
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Println("error in mapping struct ", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	_, err := h.store.FindUserByEmail(c.Request.Context(), newUser.Email)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "user already registered"})
		return
	}

	newUser.Type = "user"

	if err := h.store.CreateUser(c, newUser); err != nil {
		log.Println("user already registered ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "user already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Registered successfully"})
}

func checkPasswordHash(password, bcryptHash string) bool {
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(bcryptHash), passwordBytes)
	return err == nil
}

// fetch users data
func (h *Handler) FetchUser(c *gin.Context) {
	var newUser model.Users
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Println("error in mapping struct ", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	_, err := h.store.FindUserByEmail(c.Request.Context(), newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "user already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "user is registered"})
}

// verify token is it valid or not
func VerifyToken(c *gin.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "valid token", "status": http.StatusOK})
}

// fetch all users for admin
func (h *Handler) FetchAllUsers(c *gin.Context) {
	users, err := h.store.FetchUsers(c)

	if err != nil {
		log.Println("error in fetch users ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching users", "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users, "status": http.StatusOK})
}

// update status of user by admin
func (h *Handler) UpdateStatus(c *gin.Context) {
	body, exists := c.Get("body")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve body from context"})
		return
	}

	// Convert interface{} to map[string]interface{}
	data, ok := body.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Body is not in the expected format"})
		return
	}

	email, _ := data["email"].(string)
	status, _ := data["status"].(string)
	err := h.store.UpdateStatus(c.Request.Context(), email, status)
	if err != nil {
		log.Println("error in session update ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "error in session update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "valid token", "status": http.StatusOK})
}

// send mail to admin by user to contact regarding query
func (h *Handler) SendEmail(c *gin.Context) {
	body, exists := c.Get("body")

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve body from context"})
		return
	}

	// Convert interface{} to map[string]interface{}
	data, ok := body.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Body is not in the expected format"})
		return
	}

	email, _ := data["email"].(string)
	subj, _ := data["subject"].(string)
	content, _ := data["content"].(string)

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error", "status": http.StatusInternalServerError})
		c.Abort()
		return
	}

	apiKey := cfg.EmailAPIKey

	client := resend.NewClient(apiKey)

	content = fmt.Sprintf("<p> %v </p>", content)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"travelplanner527@gmail.com"},
		Subject: subj,
		Html:    content,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println("Error in sending mail ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in sending mail", "status": http.StatusInternalServerError})
		return
	}

	log.Println("sent successfully from ", email, sent)
	c.JSON(http.StatusOK, gin.H{"message": "sent successfully", "status": http.StatusOK})
}
