package http

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ref-message-hub/common/ecode"
	http2 "ref-message-hub/common/http"
	"ref-message-hub/common/log"
	"ref-message-hub/conf"
	g "ref-message-hub/http/gin"
	"strconv"
	"time"
)

var (
	srv *http.Server
)

func InitHttpServer() {
	r := gin.New()
	if conf.Conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(g.Recovery())
	r.Use(g.Log())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     conf.Conf.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "locale"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	app := r.Group("/api")
	{
		app.GET("/status", ping)
		app.POST("/message/send", sendMessage)
	}

	srv = &http.Server{
		Addr:    ":" + strconv.Itoa(conf.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Http Server listen error: %s", err)
		}
	}()
}

func ShutdownHttpServer() {
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Http Server Shutdown error:", err)
		}
		log.Info("Http Server exiting")
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, http2.Resp(ecode.OK, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}))
}
