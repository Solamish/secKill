package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"strings"
)

/**
  Simple 模式
*/

type SecResult struct {
	ProductId uint   `json:"product_id"`
	UserId    string `json:"user_id"`
	Code      int    `json:"code"` //状态码
}

// 创建simple模式的mq
func NewSimplePattern(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

func (m *RabbitMQ) PublishSimple(message string) {
	m.Lock()
	defer m.Unlock()
	// 1.申请队列，如果队列不存在会自动创建
	_, err := m.channel.QueueDeclare(m.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否排他
		false,
		// 是否阻塞
		false,
		// 额外参数
		nil)
	if err != nil {
		// TODO
		log.Panic("fail to get queue", err)
		fmt.Println(err)
	}
	// 2.发送消息到队列

	_ = m.channel.Publish(
		m.Exchange,
		m.QueueName,
		// 如果设置为true， 会根据exchange的类型和routekey规则寻找队列，
		// 如果无法找到符合规则的队列，那么发送的消息会返回给发送者
		false,
		// 如果设置为true, 当exchange发送消息到消息队列后发现该队列上没有绑定消费者，则会把消息返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

func (m *RabbitMQ) ConsumeSimple(createOrder func(string, uint)) {
	// 1.声明队列
	q, err := m.channel.QueueDeclare(m.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否排他
		false,
		// 是否阻塞
		false,
		// 额外参数
		nil)
	if err != nil {
		// TODO
		fmt.Println(err)
	}
	// TODO
	////消费者流控
	//m.channel.Qos(
	//	1, //当前消费者一次能接受的最大消息数量
	//	0, //服务器传递的最大容量（以八位字节为单位）
	//	false, //如果设置为true 对channel可用
	//)
	// 2.接受消息
	msgs, err := m.channel.Consume(
		q.Name,
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true, 表示不能将同一Connection中发送的消息，不能发给这个Connection中的消费者
		false,
		// 是否阻塞
		false,
		nil, )
	forever := make(chan bool)
	go func() {
		for d := range msgs {

			message := SecResult{}

			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Println("json.Unmarshal failed, Error:", err)
				continue
			}
			writeLog(message.UserId+"&"+strconv.Itoa(int(message.ProductId)), "./log/mq.log")
			// 创建订单
			createOrder(message.UserId, message.ProductId)

			// 如果为true表示确认所有未确认的消息，
			// 为false表示确认当前消息
			// TODO
			//d.Ack(false)
		}
	}()

	<-forever
}

func writeLog(msg string, logPath string) {
	fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()
	content := strings.Join([]string{msg, "\r\n"}, "")
	buf := []byte(content)
	fd.Write(buf)
}
