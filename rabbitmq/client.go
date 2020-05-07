package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// url 格式 amqp://账号：密码@rabbit的服务器地址：端口号/vhost
const (
	MQURL = "amqp://guest:guest@127.0.0.1:5672"
)

var (
	RabbitMqConsumer *RabbitMQ
	RabbitMqProducer *RabbitMQ
)

func InitMQ() {
	RabbitMqProducer = NewSimplePattern("SecKill")
	RabbitMqConsumer = NewSimplePattern("SecKill")
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind key 名称
	Key   string
	MqUrl string
	sync.Mutex
}

// 创建RabbitMq实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	var err error
	rabbitMq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     MQURL,
	}
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	if err != nil {
		log.Println("fail to create MQ connection:",err)
	}

	rabbitMq.channel, err = rabbitMq.conn.Channel()
	if err != nil {
		log.Println("fail to get MQ channel:",err)
	}

	return rabbitMq
}

func (m *RabbitMQ) Destroy() {
	m.channel.Close()
	m.conn.Close()
}

 

