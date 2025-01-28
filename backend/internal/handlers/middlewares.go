package handlers

import (
	"net/http"
	"strings"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/gin-gonic/gin"
	// "github.com/redis/go-redis/v9"
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

// func loggerMiddleware() gin.HandlerFunc {
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

// 	return func(c *gin.Context) {
// 		start := time.Now()

// 		c.Next()

// 		duration := time.Since(start)

// 		var errors []error
// 		for _, err := range c.Errors {
// 			errors = append(errors, err)
// 		}

// 		logger := log.Info()
// 		if len(c.Errors) > 0 {
// 			logger = log.Error().Errs("errors", errors)
// 		}

// 		logger.
// 			Str("method", c.Request.Method).
// 			Str("path", c.Request.RequestURI).
// 			Int("status_code", c.Writer.Status()).
// 			Str("status_text", http.StatusText(c.Writer.Status())).
// 			Dur("duration", duration)
// 	}
// }

// func redisCacheMiddleware(redis *redis.Client) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		requestPath := ctx.Request.URL.Path

// 		cacheData, err := redis.Get(ctx, requestPath).Bytes()
// 		if err == nil {
// 			log.Info().
// 				Msgf("cached hit for: %v", requestPath)

// 			var jsonData any
// 			if err := json.Unmarshal(cacheData, &jsonData); err != nil {
// 				ctx.AbortWithError(http.StatusInternalServerError, errors.New("could not unmarshal redis cache"))

// 				return
// 			}

// 			// ctx.Data(http.StatusOK, "application/json", cacheData)
// 			ctx.AbortWithStatusJSON(http.StatusOK, jsonData)

// 			return
// 		}

// 		ctx.Next()
// 	}
// }

func CORSmiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(204)

			return
		}

		ctx.Next()
	}
}