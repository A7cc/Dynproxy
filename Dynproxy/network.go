package Dynproxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

// 初始化http客户端连接
func InitClient(timeout int) error {
	// 设置tr
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 设置客户端
	Client = &http.Client{
		// 设置请求信息
		Transport: tr,
		// 设置超时时间
		Timeout: time.Duration(timeout) * time.Second,
	}
	return nil
}

// 开启代理
func Strartsocks(proxyhost []string, proxyport string) {
	// 开启监听本地1080
	listener, err := net.Listen("tcp", "127.0.0.1:"+proxyport)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// 输出
	fmt.Println(RIGHT, "开始监听socks端口: 127.0.0.1:"+proxyport)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go changesocks(proxyhost, conn.(*net.TCPConn))
	}
}

// 使用代理
func changesocks(proxyhost []string, ws *net.TCPConn) {
	socksproxy := proxyhost[rand.Intn(len(proxyhost))]
	fmt.Printf("\u001B[2K\r%s 当前使用代理%s", RIGHT, socksproxy)
	defer ws.Close()
	socks, err := net.DialTimeout("tcp", socksproxy, 5*time.Second)
	if err != nil {
		log.Println("dial socks error:", err)
		for i := 0; i < len(proxyhost); i++ {
			if proxyhost[i] == socksproxy {
				proxyhost = append(proxyhost[:i], proxyhost[i+1:]...)
			}
		}
		changesocks(proxyhost, ws)
		return
	}
	defer socks.Close()
	var wg sync.WaitGroup
	ioCopy := func(dst io.Writer, src io.Reader) {
		defer wg.Done()
		io.Copy(dst, src)
	}
	wg.Add(2)
	go ioCopy(socks, ws)
	go ioCopy(ws, socks)
	wg.Wait()
}

// New 新建一个协程池
func PoolNew(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add 新增一个执行
func (p *Pool) Add(delta int) {
	// delta为正数就添加
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	p.wg.Add(delta)
}

// Done 执行完成减一
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

// 等待全部协程结束
func (p *Pool) Wait() {
	p.wg.Wait()
}
