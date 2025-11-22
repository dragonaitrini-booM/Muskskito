package main

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/muskskito/muskskito/database"
    "github.com/muskskito/muskskito/handlers"
    "github.com/muskskito/muskskito/internal/firewall"
    "github.com/muskskito/muskskito/middleware"
    "log"
    "os"
)

func main() {
    if _, err := os.Stat(".env"); err == nil {
        err := godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }
    }

    database.Connect()
    database.Migrate()

    r := gin.Default()

    sessionSecret := os.Getenv("SESSION_SECRET")
    if sessionSecret == "" {
        log.Fatal("SESSION_SECRET is not set")
    }
    store := cookie.NewStore([]byte(sessionSecret))
    r.Use(sessions.Sessions("mysession", store))

    auth := r.Group("/auth")
    {
        auth.GET("/login", handlers.HandleLogin)
        auth.GET("/callback", handlers.HandleCallback)
        auth.GET("/me", handlers.GetMe)
        auth.POST("/logout", handlers.Logout)
    }

    firewallAgent := firewall.New()

    protected := r.Group("/")
    protected.Use(middleware.AuthMiddleware())
    protected.Use(firewallAgent.Middleware())

    settings := protected.Group("/settings")
    {
        settings.GET("/", handlers.GetSettings)
        settings.PUT("/", handlers.UpdateSettings)
    }

    vpn := protected.Group("/vpn")
    {
        vpn.GET("/locations", handlers.GetVpnLocations)
        vpn.POST("/connect", handlers.ConnectVpn)
    }

    browse := protected.Group("/browse")
    {
        browse.POST("/scan", handlers.ScanUrl)
        browse.POST("/start", handlers.StartBrowsing)
        browse.POST("/nuke", handlers.NukeSession)
        browse.GET("/sessions", handlers.GetSessions)
    }

    chat := protected.Group("/chat")
    {
        chat.GET("/messages", handlers.GetChatMessages)
        chat.POST("/send", handlers.SendChatMessage)
    }

    analytics := protected.Group("/analytics")
    {
        analytics.GET("/stats", handlers.GetAnalyticsStats)
    }

    subscription := protected.Group("/subscription")
    {
        subscription.POST("/create-payment", handlers.CreatePayment)
        subscription.POST("/confirm-payment", handlers.ConfirmPayment)
    }

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
