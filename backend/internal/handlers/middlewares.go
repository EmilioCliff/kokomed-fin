package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
	authorizationPayloadKey       = "payload"
)

func authMiddleware(maker pkg.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "No header was passed")))

			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Invalid or Missing Bearer Token")))

			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != authorizationHeaderBearerType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Authentication Type Not Supported")))

			return
		}

		token := fields[1]

		payload, err := maker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Access Token Not Valid")))

			return
		}

		ctx.Set(authorizationPayloadKey, payload)

		ctx.Next()
	}
}

func redisCacheMiddleware(cache services.CacheService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path
		queryParams := ctx.Request.URL.Query()

		cacheKey := constructCacheKey(requestPath, queryParams)

		var target any

		log.Println(cacheKey)

		exists, _ := cache.Get(ctx, cacheKey, &target)
		if exists {
			log.Println("Cache hit for: ", cacheKey)

			ctx.AbortWithStatusJSON(http.StatusOK, target)

			return
		}

		ctx.Next()
	}
}

func CORSmiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
        "https://frontend-production-ca6a.up.railway.app",
		"https://frontend-development-fdd2.up.railway.app",
		"http://localhost:3001",
    }

	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
        
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }

		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true") 
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400") 
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")       

		// Handle preflight (OPTIONS) requests
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}