package main

import (
	"GO/cmd/ip-notify/utils"
	"fmt"
	"time"

	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
)

// INFO webhook.urlにIPアドレス通知先のURLを記述

//go:embed webhook.url
var discordWebhookURL string

func main() {
	// INFO DNSサービス起動とIP取得を遅延させる
	time.Sleep(10 * time.Second)

	ips, err := utils.GetIPAddress()
	if err != nil {
		fmt.Println("IPアドレスの取得に失敗しました")
		return
	}

	message := map[string]string{
		"content": "IP Address:\n" + strings.Join(ips, "\n"),
	}
	payload, _ := json.Marshal(message)

	http.Post(strings.TrimSpace(discordWebhookURL), "application/json", bytes.NewBuffer(payload))
}
