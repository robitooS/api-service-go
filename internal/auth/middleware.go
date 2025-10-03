package auth

import (
	"bytes"
	"context"

	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/domain/user"
)

func AuthenticateHMAC(hmacKey []byte, repository user.UserRepository, cache cache.NonceStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		tsStr := ctx.GetHeader("X-Timestamp")
		authHeader := ctx.GetHeader("Authorization")
		nonce := ctx.GetHeader("X-Nonce")
		userIDStr := ctx.GetHeader("X-User-ID")

		if tsStr == "" || nonce == "" || authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "headers de autenticação ausentes"})
			return
		}

		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "X-Timestamp inválido"})
			return
		}

		if err := ValidateTimeStamp(ts); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "timestamp inválido ou expirado"})
			return
		}
		if err := cache.CacheNonce(nonce); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "nonce já utilizado"})
			return
		}
		
		msg := BuildMessage(method, path, ts, string(bodyBytes), nonce)
		

		ok, err := ValidateSignature(msg, authHeader, hmacKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "erro ao validar assinatura"})
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "assinatura inválida"})
			return
		}
		
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		if !verifyUserExists(repository, ctx.Request.Context(), userID) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

// (O resto do arquivo permanece igual)
func verifyUserExists(r user.UserRepository, ctx context.Context, userID int64) bool {
	_, err := r.FindByID(ctx, userID)
	return err == nil
}