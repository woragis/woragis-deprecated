package security

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware sets security headers to protect against common attacks
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// XSS Protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")

		// Referrer Policy
		c.Set("Referrer-Policy", "no-referrer")

		// Content Security Policy
		// Allow same-origin, inline scripts/styles for development
		// In production, consider stricter CSP
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self' https:; " +
			"frame-ancestors 'none';"
		c.Set("Content-Security-Policy", csp)

		// Permissions Policy (formerly Feature Policy)
		permissionsPolicy := "geolocation=(), microphone=(), camera=()"
		c.Set("Permissions-Policy", permissionsPolicy)

		// Strict Transport Security (only if using HTTPS)
		// Uncomment in production with HTTPS
		// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		return c.Next()
	}
}

