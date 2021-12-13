package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/SakuraBurst/gitlab-bot/parser"
	"github.com/SakuraBurst/gitlab-bot/telegram"
)

const CAN_BE_MERGED = "can_be_merged"

func main() {
	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			// UNSAFE!
			// DON'T USE IN PRODUCTION!
			InsecureSkipVerify: true,
		},
	}
	mrWithDiffs := parser.Parser("gpe/ais-upu/ais-upu-frontend", "ymrsGzzNEofRKhoX2f5G")
	telegram.SendMessage(mrWithDiffs)

}
