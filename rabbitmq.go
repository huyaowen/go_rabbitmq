package go_rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

//消费配置类
type ConsumerConfig struct {
	QueueName string //queue name
	Consumer  string //consumer tag
	AutoAck   bool   //no ack
	Exclusive bool   //Check Queue struct documentation
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
	Done      chan error
	Handler   func(<-chan amqp.Delivery, chan error) //处理函数
}

type MessageConfig struct {
	Body         []byte
	ContentType  string
	DeliveryMode uint8  //1:非持久化 2：持久化
	Exchange     string //交换机
	RoutingKey   string //路由键
	Mandatory    bool   //将该消息至少放到一个队列，如果不能则返回
	Immediate    bool   //当与消息routeKey关联的所有queue(一个或多个)都没有消费者时，该消息会通过basic.return方法返还给生产者
	Confirm      bool   //可信赖消息加入确认机制
}

type QueueConfig struct {
	// The queue name may be empty, in which the server will generate a unique name
	// which will be returned in the Name field of Queue struct.
	Name string

	// Check Exchange comments for durable
	Durable bool

	// Check Exchange comments for autodelete
	AutoDelete bool

	// Exclusive queues are only accessible by the connection that declares them and
	// will be deleted when the connection closes.  Channels on other connections
	// will receive an error when attempting declare, bind, consume, purge or delete a
	// queue with the same name.
	Exclusive bool

	// When noWait is true, the queue will assume to be declared on the server.  A
	// channel exception will arrive if the conditions are met for existing queues
	// or attempting to modify an existing queue from a different connection.
	NoWait bool

	// Check Exchange comments for Args
	Args amqp.Table
}

type BindingOptions struct {
	// Publishings messages to given Queue with matching -RoutingKey-
	// Every Queue has a default binding to Default Exchange with their Qeueu name
	// So you can send messages to a queue over default exchange
	RoutingKey string

	// Do not wait for a consumer
	NoWait bool

	// App specific data
	Args amqp.Table
}

//错误处理
func FailOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
