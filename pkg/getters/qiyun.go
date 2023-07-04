package getters

import (
	"encoding/json"
	"fmt"
	"go-proxypool/pkg/global"
	"go-proxypool/pkg/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Qiyun fetches ips from qiyun.
func Qiyun() (result []*models.Ip) {
	apikey := global.Config.GetString("qiyun_apikey")
	if len(apikey) == 0 {
		global.Logger.Error("unable to get qiyun apikey from config")
		return nil
	}
	url := fmt.Sprintf("http://dev.qydailiip.com/api/?apikey=%s&num=100&type=json&line=unix&proxy_type=putong&sort=1&model=all&protocol=http&address=&kill_address=&port=&kill_port=&today=false&abroad=1&isp=&anonymity=", apikey)
	ips, err := fetchAndParseIps(url)
	if err != nil {
		global.Logger.Errorf("unable to fetch and parse ips from qiyun: %v", err)
	}
	return ips
}

// fetchAndParseIps fetches and parses ips from url.
func fetchAndParseIps(url string) (result []*models.Ip, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse ips from body.
	var ipsAsStrings []string
	if err := json.Unmarshal(body, &ipsAsStrings); err != nil {
		return nil, err
	}

	var ips []*models.Ip
	for _, ipAsString := range ipsAsStrings {
		parts := strings.Split(ipAsString, ":")
		if len(parts) != 2 {
			continue
		}
		ip := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		ips = append(ips, &models.Ip{
			Ip:   ip,
			Port: port,
		})
	}

	return ips, nil
}
