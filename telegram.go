package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendTelegram(results []BackupResult) {
	bot, err := telego.NewBot(telegramBotToken, telego.WithDefaultLogger(false, true))
	if err != nil {
		slog.Error("connect telegram failed", slog.Any("error", err))
		return
	}

	botUser, err := bot.GetMe()
	if err != nil {
		slog.Error("get me telegram failed", slog.Any("error", err))
		return
	}

	slog.Info("connect telegram success", slog.Any("user", botUser))

	for _, result := range results {
		msg := fmt.Sprintf(`*VM Backup*
VMID: %s
Name: %s
Status: %s
Time: %s
Size: %s`, result.VMID, result.Name, strings.ToUpper(result.Status), result.Time, result.Size)
		if _, err = bot.SendMessage(tu.Message(tu.ID(int64(telegramChatID)), msg)); err != nil {
			slog.Error("send telegram failed", slog.Any("error", err))
		}
	}
}
