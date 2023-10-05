package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pretty66/websocketproxy"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	WebsocketHost string `json:"websocket_host"`
	WebsocketPort int    `json:"websocket_port"`
	WebsocketSSL  bool   `json:"websocket_ssl"`
	WebsocketPath string `json:"websocket_path"`
}

func main() {
	configPath := ""
	flag.StringVar(&configPath, "config", "config.json", "Path to config file")

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
		return
	}

	url := strings.Replace(config.WebsocketHost, ".", "-", -1)
	url = fmt.Sprintf("ws://%s-%d-p", url, config.WebsocketPort)
	if config.WebsocketSSL {
		url += "-s"
	}
	url = fmt.Sprintf("%s.webvpn.zju.edu.cn:8001%s", url, config.WebsocketPath)
	log.Println("Converted websocket URL:", url)

	twdId, err := login(config.Username, config.Password)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Login success, TWFID:", twdId)

	wp, err := websocketproxy.NewProxy(url, func(r *http.Request) error {
		r.Header.Set("Cookie", "TWFID="+twdId)
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/", wp.Proxy)

	log.Println("Listening on", config.Host+":"+strconv.Itoa(config.Port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
