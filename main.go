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
	Wait()
	//

}

func Wait() {
	t := time.Now()
	if t.Hour() != 12 {
		waitFor := 0
		if t.Hour() > 12 {
			waitFor = 24 - t.Hour() + 12
		} else {
			waitFor = 12 - t.Hour()
		}
		time.Sleep(time.Hour * time.Duration(waitFor))
		Wait()
	} else {
		mrWithDiffs := parser.Parser("gpe/ais-upu/ais-upu-frontend", "ymrsGzzNEofRKhoX2f5G")
		telegram.SendMessage(mrWithDiffs)
		time.Sleep(time.Hour * 24)
		Wait()
	}
}
