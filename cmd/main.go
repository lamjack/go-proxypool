package main

import (
	"context"
	"fmt"
	"go-proxypool/pkg/api"
	"go-proxypool/pkg/getters"
	"go-proxypool/pkg/global"
	"go-proxypool/pkg/models"
	"go-proxypool/pkg/utils"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	global.Initialize()

	ipChan := make(chan *models.Ip, 100)

	var wg sync.WaitGroup
	wg.Add(1)
	go utils.ValidateAndStoreProxies(ipChan, &wg)

	// Periodically validate proxies every minute.
	go utils.PeriodicProxyValidation(5 * time.Second)

	apiServer := api.NewServer(8080)
	go func() {
		err := apiServer.Run()
		if err != nil {
			panic(fmt.Sprintf("unable to start API server: %v", err))
		}
		global.Logger.Infof("API server started")
	}()

	for {
		n, _ := global.Storage.GetAll(context.Background())
		global.Logger.Infof("current ipChan length: %d, storage length: %d", len(ipChan), len(n))
		if len(ipChan) < 100 {
			global.Logger.Debugf("length of ipChan is less than 100, fetching ips from getters")
			go putIpsToChan(ipChan)
		}
		time.Sleep(30 * time.Second)
	}
}

func putIpsToChan(ipChan chan<- *models.Ip) {
	var wg sync.WaitGroup
	funs := []func() []*models.Ip{
		getters.Qiyun,
	}
	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.Ip) {
			ips := f()
			for _, ip := range ips {
				ipChan <- ip
			}
			wg.Done()
		}(f)
	}
	wg.Wait()
}
