package utils

import (
	"context"
	"fmt"
	"go-proxypool/pkg/global"
	"go-proxypool/pkg/models"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var ctx = context.Background()

// ValidateAndStoreProxies validates and stores proxies.
func ValidateAndStoreProxies(ipChan <-chan *models.Ip, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range ipChan {
		proxyURL := fmt.Sprintf("%s:%d", ip.Ip, ip.Port)
		if isValidProxy(proxyURL) {
			err := global.Storage.Add(*ip, ctx)
			if err != nil {
				global.Logger.Errorf("failed to store proxy %s: %v", proxyURL, err)
			}
		}
	}
}

// PeriodicProxyValidation periodically validates proxies.
func PeriodicProxyValidation(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		ips, err := global.Storage.GetAll(ctx)
		if err != nil {
			global.Logger.Errorf("failed to get all proxies: %v", err)
			continue
		}

		var validationWg sync.WaitGroup
		for _, ip := range ips {
			validationWg.Add(1)
			go func(ip models.Ip) {
				defer validationWg.Done()

				proxyURL := fmt.Sprintf("%s:%d", ip.Ip, ip.Port)
				if !isValidProxy(proxyURL) {
					err := global.Storage.Remove(ip, ctx)
					if err != nil {
						global.Logger.Errorf("failed to remove proxy %s: %v", proxyURL, err)
					}
				}
			}(ip)
		}

		validationWg.Wait()
	}
}

func isValidProxy(proxyURL string) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{Scheme: "http", Host: proxyURL}),
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(global.Config.GetString("test_url"))
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
