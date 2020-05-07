package main

import (
	_ "net/http/pprof"
	"secKill/model"
	"secKill/rabbitmq"
	"secKill/redis"
	"secKill/route"
	"secKill/logic"
)

func main() {
	logic.InitMap()
	for _, item := range logic.ItemMap {
		item.Monitor()
		go item.OffShelve()
	}
	model.InitDB()
	rabbitmq.InitMQ()
	go rabbitmq.RabbitMqConsumer.ConsumeSimple(logic.CreateOrder)
	redis.InitRedis()
	router := route.InitRouter()
	//ginpprof.Wrap(router)
	router.Run(":8080")

}


