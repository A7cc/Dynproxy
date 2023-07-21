package Dynproxy

import (
	"fmt"
	"net/http"
	"sync"
)

// fofa用户信息
type Fofauser struct {
	Err     bool   `json:"error"`
	Email   string `json:"email"`
	Key     string `json:"key"`
	Usename string `json:"username"`
	Vip     int    `json:"vip_level"`
}

// 数据类型
type Fofadata struct {
	Err     bool       `json:"error"`
	Size    int        `json:"size"`
	Page    int        `json:"page"`
	Mode    string     `json:"mode"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
}

// 终端颜色
var (
	ERR     string = "[\033[1;31m✘\033[0m]"
	RIGHT   string = "[\033[1;32m✓\033[0m]"
	WARNING string = "[\033[0;38;5;214m!\033[0m]"
	YELLOW  string = "\033[1;33m"
	MAIN    string = "\033[38;5;50m"
	ORANGE  string = "\033[0;38;5;214m"
	GREEN   string = "\033[1;32m"
	BLUE    string = "\033[1;34m"
	RED     string = "\033[1;31m"
	END     string = "\033[0m"
)

// 设置协程结构体
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// 基础配置
var (
	// 设置请求时的信息
	Client *http.Client
	// 输入Email
	Email string
	// 输入key
	Key string
	// 搜索的语句
	Query string
	// 查询条数
	Size int
	// 搜索内容
	Title string
	// 文件存放路径
	filePath string
	// 文件存放路径
	ProxyfilePath string
	// 超时时间
	Timeout int
	// 模式
	Mode string
	// 设置线程
	Threadnum int
	// 设置要代理的端口
	Proxyport int
	// 设置地址
	Region string
	// 版本
	version string = "1.0.0"
)

// 标记
func TagPrint(version string) {
	tag := `     __    __ 
    \  \__/  /
     \_*  *_/_    _ ____ __ __________
     /|| :||  \ | |  . \  * ___ _ . _|\ /|
  __| ||  || | \|*|  __/|  /_^_ \ V /|.|.| 
 / .. ||  || \\   | |   | |{ * } > < |_  |
 \____\\  /|__\\  | |   |_|{_|_}/_^_\ _| |
        \/      \_|/ version: %8s/__|  
`
	fmt.Printf(tag, version)
}
