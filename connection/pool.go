package connection

import (
	"time"
)

var nowFunc = time.Now // for testing

type ConnPool struct {
	//获取连接函数
	Dial func() (interface{}, error)
	//最大连接数
	MaxActive int
	idle      chan interface{}
}

type idleConn struct {
	connection interface{}
	timer      time.Time
}

// 批量生成连接，并把连接放到连接池channel里面
//带缓冲channel，无阻塞
func (this *ConnPool) InitPool() error {

	this.idle = make(chan interface{}, this.MaxActive)
	for x := 0; x < this.MaxActive; x++ {
		db, err := this.Dial()
		if err != nil {
			return err
		}
		this.idle <- idleConn{timer: nowFunc(), connection: db}
	}
	return nil

}

// 从连接池里取出连接

func (this *ConnPool) Get() interface{} {

	// 如果空闲连接为空，初始化连接池
	if this.idle == nil {
		this.InitPool()
	}
	ic := <-this.idle
	conn := ic.(idleConn).connection
	defer this.Release(conn)
	return conn

}

// 回收连接到连接池
func (this *ConnPool) Release(conn interface{}) {
	this.idle <- idleConn{timer: nowFunc(), connection: conn}
}
