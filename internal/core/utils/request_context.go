package utils

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestContext representa informações de contexto da requisição
type RequestContext struct {
	IPAddress string
	UserAgent string
	Method    string
	Path      string
}

// ExtractRequestContext extrai informações de contexto da requisição HTTP
func ExtractRequestContext(c *gin.Context) *RequestContext {
	if c == nil || c.Request == nil {
		return &RequestContext{}
	}

	return &RequestContext{
		IPAddress: extractRealIP(c.Request),
		UserAgent: c.Request.UserAgent(),
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
	}
}

// ExtractRequestContextFromContext extrai informações do contexto se disponível
func ExtractRequestContextFromContext(ctx context.Context) *RequestContext {
	if ginCtx, exists := ctx.(*gin.Context); exists {
		return ExtractRequestContext(ginCtx)
	}

	// Se não é um contexto Gin, retorna valores vazios
	return &RequestContext{}
}

// extractRealIP extrai o IP real do cliente, considerando proxies e load balancers
func extractRealIP(r *http.Request) string {
	// Verifica headers de proxy primeiro
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" && realIP != "unknown" {
		return realIP
	}

	// Verifica X-Forwarded-For (pode conter múltiplos IPs)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" && forwarded != "unknown" {
		// Pega o primeiro IP da lista (cliente original)
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Verifica outros headers comuns
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	if trueIP := r.Header.Get("True-Client-IP"); trueIP != "" {
		return trueIP
	}

	// Se nada encontrado, usa RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// IsLocalIP verifica se um IP é local/privado
func IsLocalIP(ip string) bool {
	if ip == "" {
		return false
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Verifica se é loopback
	if parsedIP.IsLoopback() {
		return true
	}

	// Verifica se é IP privado
	return parsedIP.IsPrivate()
}

// SanitizeUserAgent remove caracteres potencialmente perigosos do User-Agent
func SanitizeUserAgent(userAgent string) string {
	if len(userAgent) > 500 {
		userAgent = userAgent[:500]
	}

	// Remove caracteres de controle
	sanitized := strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, userAgent)

	return strings.TrimSpace(sanitized)
}
