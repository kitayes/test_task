package delivery

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	"log"
	"net/http"
	_ "test_task/docs"
	"test_task/internal/application"
	"time"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Config struct {
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}

type Handler struct {
	cfg        *Config
	services   *application.Service
	httpServer *http.Server
	router     *gin.Engine
	logger     Logger
}

func NewHandler(services *application.Service, cfg *Config, logger Logger) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
		logger:   logger,
	}
}

func (h *Handler) Run(_ context.Context) {

	h.httpServer = &http.Server{
		Addr:         ":" + h.cfg.Port,
		Handler:      h.router,
		ReadTimeout:  h.cfg.ReadTimeOut,
		WriteTimeout: h.cfg.WriteTimeOut,
	}
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			log.Println("listen: %s\n", err.Error())
			return
		}
	}()
}

func (h *Handler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Handler) Init() error {
	router := gin.Default()

	sub := router.Group("/subscriptions")
	{
		sub.POST("/", h.createSubscription)
		sub.GET("/", h.listSubscriptions)
		sub.GET("/summary", h.getSummary)
		sub.GET("/:id", h.getSubscriptionByID)
		sub.PUT("/:id", h.updateSubscription)
		sub.DELETE("/:id", h.deleteSubscription)
	}

	h.router = router
	return nil
}
