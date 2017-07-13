package connection

import (
	"fmt"

	"github.com/streadway/amqp"
)

var rabbitmqPool *ConnPool

func init() {
	rabbitmqPool = RabbitPool()
}

func RabbitPool() *ConnPool {

	qUser := "guest"
	qPass := "guest"
	qHost := "localhost"
	qPort, _ := "5276"
	rabbitmq := fmt.Sprintf("amqp://%s:%s@%s:%d/", qUser, qPass, qHost, qPort)
	poolNum := 10
	coonPool := &ConnPool{
		MaxActive: poolNum,
		Dial: func() (interface{}, error) {
			conn, err := amqp.Dial(rabbitmq)
			return conn, err
		},
	}

	return coonPool

}

func GetConnection() *amqp.Connection {
	connc := rabbitmqPool.Get().(*amqp.Connection)
	fmt.Println(connc)
	return connc
}
