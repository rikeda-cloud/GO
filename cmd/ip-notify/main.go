package main

import (
	"GO/cmd/ip-notify/utils"
	"GO/internal/config"

	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	cfg := config.GetConfig()
	ips, err := utils.GetIPAddress()
	if err != nil {
		return
	}

	message := map[string]string{
		"content": "IP Address:\n" + strings.Join(ips, "\n"),
	}
	payload, _ := json.Marshal(message)

	http.Post(cfg.App.IpNotify.DiscordWebhookUrl, "application/json", bytes.NewBuffer(payload))
}
