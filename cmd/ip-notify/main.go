package main

import (
	"GO/cmd/ip-notify/utils"

	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const webhookURL = "https://discord.com/api/webhooks/11111111111"

func main() {
	ips, err := utils.GetIPAddress()
	if err != nil {
		return
	}

	message := map[string]string{
		"content": "IP Address:\n" + strings.Join(ips, "\n"),
	}
	payload, _ := json.Marshal(message)

	http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
}
