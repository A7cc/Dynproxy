package Dynproxy

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 处理命令行
func ParseFlag() {
	// 输出tag
	TagPrint(version)
	// 解析命令行参数
	flag.IntVar(&Size, "size", 100, "需要通过网络测绘查询的数据条数")
	flag.StringVar(&Key, "k", "", "登录fofa的key")
	flag.StringVar(&Email, "e", "", "登录fofa的邮箱")
	flag.StringVar(&Query, "q", "title='bilibili'", "要搜索的语句，在getproxy模式下该参数无效")
	flag.IntVar(&Threadnum, "n", 200, "设置验证存活代理的协程数量，默认为200")
	flag.StringVar(&Title, "title", "host,ip,port,country,city,title", "搜索输出的内容，默认追加host，如果需要搜索其他输出内容可以使用模式mode=outtitle")
	flag.IntVar(&Timeout, "time", 5, "请求超时时间")
	flag.StringVar(&ProxyfilePath, "pf", "", "读取文件中的代理信息并使用")
	flag.StringVar(&Mode, "m", "search", "使用的功能模式:getproxy/search/runproxy/outtitle")
	flag.IntVar(&Proxyport, "port", 1080, "设置要代理的端口")
	flag.StringVar(&Region, "rg", "abroad", "要使用的代理地区：")
	flag.Parse()
}

// 处理代理
func Procesproxy(proxyhosts []string) error {
	// 检测存活性
	lives, chinas, abroads := Proxylivecheck(proxyhosts, Threadnum)
	fmt.Println()
	fmt.Println(RIGHT, "一共获取存活代理:", len(lives), "条")
	fmt.Println(RIGHT, "一共获取地区为中国的代理:", len(chinas), "条")
	fmt.Println(RIGHT, "一共获取地区为国外的代理:", len(abroads), "条")
	// 写入文件
	for _, live := range lives {
		err := WriteFile(live, "dynproxy.txt")
		if err != nil {
			return err
		}
	}
	// 写入文件
	for _, live := range lives {
		err := WriteFile(live, "dynproxychinas.txt")
		if err != nil {
			return err
		}
	}
	// 写入文件
	for _, live := range lives {
		err := WriteFile(live, "dynproxyabroads.txt")
		if err != nil {
			return err
		}
	}
	fmt.Println()
	var dynproxy []string
	switch Region {
	case "china":
		dynproxy = chinas
	case "abroad":
		dynproxy = abroads
	default:
		dynproxy = lives
	}

	// 写入代理成功的IP
	Strartsocks(dynproxy, strconv.Itoa(Proxyport))
	return nil
}

// 处理getproxy
func Getproxy(results [][]string) ([]string, error) {
	if len(results) < 1 {
		return nil, errors.New("没有可用的内容")
	}
	var proxyhosts []string
	// 循环获取所有的内容
	for _, value := range results {
		proxyhosts = append(proxyhosts, value[0])
	}
	return proxyhosts, nil
}

// 检测存活性
func Proxylivecheck(proxyhosts []string, threadnum int) (lives, chinas, abroads []string) {
	// 输出内容
	fmt.Println(RED + "======================探测存活======================" + END)
	// 创建协程池
	pool := PoolNew(threadnum)
	// 记录已经检测的IP
	currentdata := 0
	lenth := len(proxyhosts)
	fmt.Print("\033[s")
	// 循环遍历内容
	for i := 0; i < lenth; i++ {
		// 添加一个协程
		pool.Add(1)
		go func(tempsocks string) {
			// 完成一个协程
			defer pool.Done()
			// 判断是否存活，如果存活那么就加入
			region, flag := IsProxy(tempsocks)
			if flag {
				if find := strings.Contains(region, "中国"); find {
					chinas = append(chinas, tempsocks)
				} else {
					abroads = append(abroads, tempsocks)
				}
				lives = append(lives, tempsocks)
			}
			// 添加已经测试的
			currentdata = currentdata + 1
			fmt.Printf("\u001B[2K\r%s 已检测%.2f%%,当前检测IP为:%s", RIGHT, float32(currentdata*100)/float32(lenth), tempsocks)
		}(proxyhosts[i])
	}
	// 完成测试
	pool.Wait()
	return
}

// 探测代理是否正常
// TODO:可用验证是否为国内的
func IsProxy(proxyIp string) (from string, isProxy bool) {
	// 设置socker5
	proxyUrl := fmt.Sprintf("socks5://%s", proxyIp)
	// 判断url格式
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return "", false
	}
	// 设置客户端信息
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// 忽略对证书的校验
			InsecureSkipVerify: true,
		},
		// 设置代理
		Proxy: http.ProxyURL(proxy),
	}
	Client.Transport = transport
	// 设置检测代理
	res, err := Client.Get("http://myip.ipip.net")
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode == 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err == nil && strings.Contains(string(body), "当前 IP") {
				// 正则表达式
				var reBody = regexp.MustCompile(`来自于：(.*)`)
				titletmp := reBody.FindStringSubmatch(string(body))
				fmt.Printf("\u001B[2K\r[+] 当前 IP：%s，来自于：%s", proxyIp, titletmp[1])
				return titletmp[1], true
			} else {
				return "", false
			}
		} else {
			return "", false
		}
	}
	return "", false
}

// 文件写入
func WriteFile(result string, filename string) error {
	var text = []byte(result + "\n")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	_, err = f.Write(text)
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

// 读取文件并以[]string形式返回内容
func ReadingLines(filename string) []string {
	// 设置结果变量
	var result []string
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer file.Close()
	// 开始读取文件
	scanner := bufio.NewScanner(file)
	// 循环读取
	for scanner.Scan() {
		val := scanner.Text()
		// 判断是否为空
		if val == "" {
			continue
		}
		// 拼接
		result = append(result, val)
	}
	// 判断是否存在终止错误
	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}
