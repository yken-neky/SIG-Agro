package http2

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sig-agro/services/user-service/internal/service"
)

type UserHTTPHandler struct {
	Service *service.UserService
}

func RegisterHandlers(r *gin.Engine, srv *service.UserService) {
	block := r.Group("/users")
	block.POST("/register", Register(srv))
	block.POST("/login", Login(srv))
	block.GET("/users/:id", GetUser(srv))
	block.GET("/users", ListUsers(srv))
}

func Register(srv *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
			FullName string `json:"full_name" binding:"required"`
			Phone    string `json:"phone" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := srv.Register(c.Request.Context(), req.Email, req.Password, req.FullName, req.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"user_id":   user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"message":   "User registered successfully",
		})
	}
}

func Login(srv *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, token, err := srv.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":    user.ID,
			"email":      user.Email,
			"token":      token,
			"expires_at": user.CreatedAt.Add(time.Hour).Unix(),
		})
	}
}

func GetUser(srv *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := srv.GetUserByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func ListUsers(srv *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := srv.ListUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
