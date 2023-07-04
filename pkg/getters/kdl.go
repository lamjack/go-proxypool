package getters

import (
	"encoding/json"
	"go-proxypool/pkg/global"
	"go-proxypool/pkg/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type kdjResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Count          int      `json:"count"`
		ProxyList      []string `json:"proxy_list"`
		TodayLeftCount int      `json:"today_left_count"`
		DedupCount     int      `json:"dedup_count"`
	} `json:"data"`
}

func KuaiDaiLi() (result []*models.Ip) {
	api := global.Config.GetString("kdj_api")
	if api == "" {
		global.Logger.Warnf("unable to get 快代理 api from config")
		return nil
	}

	resp, err := http.Get(api)
	if err != nil {
		global.Logger.Errorf("unable to fetch from 快代理: %v", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.Logger.Errorf("unable to read body from 快代理: %v", err)
		return nil
	}

	var respBody kdjResponse
	if err := json.Unmarshal(body, &respBody); err != nil {
		global.Logger.Errorf("unable to parse body from 快代理: %v", err)
		return nil
	}

	var ips []*models.Ip
	if respBody.Data.Count > 0 {
		for _, ipAsString := range respBody.Data.ProxyList {
			parts := strings.Split(ipAsString, ":")
			if len(parts) != 2 {
				continue
			}
			ip := parts[0]
			port, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			ips = append(ips, &models.Ip{
				Ip:   ip,
				Port: port,
			})
		}
	}

	return ips
}
