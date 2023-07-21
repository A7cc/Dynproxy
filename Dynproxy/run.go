package Dynproxy

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func Run() {
	// 处理命令行
	ParseFlag()
	// 初始化客户端
	InitClient(Timeout)
	// 根据模式选择
	switch {
	case Mode == "search":
		title := strings.Split(Title, ",")
		if len(title) == 0 {
			fmt.Println(WARNING, "您输入的查询title的为空，所以默认使用'host,ip'标签")
			Title = "host,ip"
		} else if len(title) == 1 {
			fmt.Println(WARNING, "您输入的查询title的过少，所以默认追加'host'标签")
			Title = Title + ",host"
		}
	case Mode == "getproxy":
		Query = "protocol='socks5' && 'Method:No Authentication(0x00)'"
		Title = "host,ip"
		if Size < 10000 {
			fmt.Println(WARNING, "您输入的获取页面数量不超过 10000 条，为了防止获取的代理过少，所以使用默认数量 10000 条")
			Size = 10000
		}
	case Mode == "runproxy":
		proxyhosts := ReadingLines(ProxyfilePath)
		if len(proxyhosts) == 0 {
			fmt.Println(ERR, "您输入代理文件没有合适的代理主机，如果没有指定文件请使用pf参数进行指定")
			return
		} else {
			err := Procesproxy(proxyhosts)
			if err != nil {
				fmt.Println(ERR, "处理代理失败 Error：", err.Error())
				return
			}
		}
	case Mode == "outtitle":
		// 列出输出字段
		OutQuery(1)
		return
	default:
		OutQuery(2)
		return
	}
	// 判断有没有key
	if Key == "" || Email == "" {
		fmt.Println(ERR, "没有输入邮箱或者key！")
		flag.Usage()
		return
	}
	// 初始化用户
	user, err := Fofauserinfo(Email, Key)
	if err != nil {
		fmt.Println(ERR, "用户登录错误 Error：", err.Error())
		return
	}
	// 输出用户信息
	ok := UserInfo(user)
	if ok == 0 {
		fmt.Println(ERR, "您没有权限调用fofa的API！即将退出！")
		return
	}

	// 查询
	results, err := user.Fofaquery(Query, Title, Size)
	if err != nil {
		fmt.Println(ERR, "查询失败 Error：", err.Error())
		return
	}
	// 输出内容
	fmt.Println(RED + "======================查询内容======================" + END)
	fmt.Println(WARNING+" 查询语句："+results.Query, "\n"+WARNING+" 查询参数："+Title, "\n"+WARNING+" 查询条数：1 -", strconv.Itoa(Size))
	switch {
	case Mode == "getproxy":
		// 处理proxy
		proxyhosts, err := Getproxy(results.Results)
		if err != nil {
			fmt.Println(ERR, "处理查询结果失败 Error：", err.Error())
			return
		}
		err = Procesproxy(proxyhosts)
		if err != nil {
			fmt.Println(ERR, "处理代理失败 Error：", err.Error())
		}
	case Mode == "search":
		err = ProcesQuery(results)
		if err != nil {
			fmt.Println(ERR, "处理查询结果失败 Error：", err.Error())
		}
	}
}
