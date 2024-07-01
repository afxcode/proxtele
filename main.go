package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	telegramBotToken = ""
	telegramChatID   = 0
	port             = 0
)

func main() {
	signalQuit := make(chan os.Signal, 2)
	signal.Notify(signalQuit, os.Interrupt, os.Kill, syscall.SIGTERM)

	if telegramBotToken = os.Getenv("BOT_TOKEN"); telegramBotToken == "" {
		slog.Error("empty BOT_TOKEN")
		return
	}
	if telegramChatID, _ = strconv.Atoi(os.Getenv("CHAT_ID")); telegramChatID == 0 {
		slog.Error("empty CHAT_ID")
		return
	}
	if port, _ = strconv.Atoi(os.Getenv("PORT")); port == 0 {
		slog.Error("empty PORT")
		return
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("message", messageHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen failed", slog.Int("port", port))
			signalQuit <- os.Kill
		}
	}()
	<-signalQuit
	slog.Info("exiting")
}

func messageHandler(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("bad request", slog.Any("error", err))
		c.Status(http.StatusBadRequest)
		return
	}

	messages := strings.Split(req.Message, "\n")

	backupResults, err := parseMessages(messages)
	if err != nil {
		slog.Warn("error parsing message output", slog.Any("error", err))
		c.Status(http.StatusBadRequest)
		return
	}

	go sendTelegram(backupResults)

	c.Status(http.StatusOK)
}
