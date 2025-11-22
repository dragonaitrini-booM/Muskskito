package firewall

import (
	"github.com/gin-gonic/gin"
)

// FirewallAgent represents the v2 firewall.
type FirewallAgent struct {
	// In a real implementation, this would hold state, configuration, etc.
}

// New creates and returns a new FirewallAgent.
func New() *FirewallAgent {
	return &FirewallAgent{}
}

// Middleware returns a Gin middleware function that enforces firewall rules.
func (fa *FirewallAgent) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rule 3: Protocol & Port Guard (Simplified)
		// This is a simplified check for HTTPS. A real implementation would
		// use eBPF for deeper packet inspection.
		if c.Request.TLS == nil {
			// Block non-HTTPS traffic
			c.AbortWithStatusJSON(403, gin.H{"error": "port_guard_violation", "reason": "insecure connection"})
			return
		}

		// All checks passed
		c.Next()
	}
}
