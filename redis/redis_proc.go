package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"secKill/proxy"
	"secKill/logic"
	"time"
)

const (
	Proxy2layerQueueName = "sec_queue"
	Layer2proxyQueueName = "recv_queue"
)

// 将请求从接入层写入逻辑层
func Proxy2layerWriteHandle() {
	for {
		req := <-proxy.SecReqChan

		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marshal req failed. Error : %v, req : %v", err, req)
			continue
		}

		err = RedisClent.RedisProxyClient.LPush(Proxy2layerQueueName, string(data))
		if err != nil {
			log.Printf("lpush req failed. Error : %v, req : %v", err, req)
			continue
		}
		//log.Printf("lpush req success. req : %v", string(data))
	}
}

// 逻辑层读取接入层发来的请求
func Proxy2layerReadHandle() {
	for {
		data, err := RedisClent.RedisLayerClient.BRPop(Proxy2layerQueueName, 60)
		if err != nil {
			log.Printf("brpop proxy2layer failed. Error : %v", err)
			continue
		}


		var req logic.SecRequest
		err = json.Unmarshal([]byte(data[1]), &req)

		if err != nil {
			log.Printf("json.Unmarshal failed. Error : %v", err)
			continue
		}

		// 请求超时
		nowTime := time.Now().Unix()
		if nowTime-req.AccessTime.Unix() >= 30 {
			log.Printf("req[%v] is expire", req)
			continue
		}

		// channel处理 超时时间
		ticker := time.NewTicker(time.Millisecond * 100)
		select {
		case <-ticker.C:
			log.Printf("send to handle chan timeout, req : %v", req)
			break
		case logic.SecHandleChan <- &req:
		}
	}
}

func Layer2proxyWriteHandle() {
	for {
		res := <-logic.SecWriteChan
		data, err := json.Marshal(res)
		if err != nil {
			log.Printf("json.Marshal res failed. Error : %v, req : %v", err, res)
			continue
		}
		err = RedisClent.RedisLayerClient.LPush(Layer2proxyQueueName,string(data))
		if err != nil {
			log.Printf("lpush req failed. Error : %v, res : %v", err, res)
			continue
		}
	}
}

// 接入层读取逻辑层的秒杀结果
func Layer2proxyReadHandle() {
	for {

		//阻塞弹出
		data, err := RedisClent.RedisProxyClient.BRPop(Layer2proxyQueueName, 60)
		if err != nil {
			log.Printf("brpop layer2proxy failed. Error : %v", err)
			continue
		}

		var result *proxy.SecResult
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			log.Printf("json.Unmarshal failed. Error : %v", err)
			continue
		}

		userKey := fmt.Sprintf("%s_%d", result.UserId, result.ProductId)

		proxy.UserCon.UserConnMapLock.Lock()
		resultChan, ok := proxy.UserCon.UserConnMap[userKey]
		proxy.UserCon.UserConnMapLock.Unlock()

		if !ok {
			log.Printf("user not found : %v", userKey)
			continue
		}

		resultChan <- result
		//log.Printf("request result send to chan success, userKey : %v", userKey)
	}
}

//  初始化redis进程
func initRedisProcess() {
    for i := 0 ; i <30 ; i++ {
 		go  Proxy2layerWriteHandle()
	}

	for i := 0 ; i <20 ; i++ {
		go  Proxy2layerReadHandle()
	}

	for i := 0 ; i <30 ; i++ {
		go  Layer2proxyWriteHandle()
	}

	for i := 0 ; i <20 ; i++ {
		go  Layer2proxyReadHandle()
	}
}
