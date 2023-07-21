package Dynproxy

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/liushuochen/gotable"
)

// 用户初始化
func Fofauserinfo(email, key string) (*Fofauser, error) {
	// 构造URL
	URL := "https://fofa.info/api/v1/info/my?email=" + email + "&key=" + key
	// 请求
	requ, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	requ.Header.Set("Accept", "*/*;q=0.8")
	requ.Header.Set("Connection", "close")
	requ.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36")
	resp, err := Client.Do(requ)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 初始化用户信息
	var fofauser *Fofauser
	json.Unmarshal(body, &fofauser)
	if fofauser.Err {
		return nil, errors.New("用户登录错误！")
	}
	// 由于fofa接口没有key值需要获取key值
	fofauser.Key = key
	return fofauser, nil
}

// 输出用户基本信息
func UserInfo(user *Fofauser) int {
	fmt.Println(RED + "======================个人信息======================" + END)
	fmt.Println(WARNING+" 邮箱：", user.Email, "\n"+WARNING+" 用户名：", user.Usename, "\n"+WARNING+" VIP等级：", user.Vip)
	return user.Vip
}

// 查询
func (fofauser *Fofauser) Fofaquery(query, fields string, size int) (*Fofadata, error) {
	query = base64.StdEncoding.EncodeToString([]byte(query))
	// 构造URL
	URL := fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&page=1&size=%d&key=%s&qbase64=%s&fields=%s", fofauser.Email, size, fofauser.Key, query, fields)
	requ, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	requ.Header.Set("Accept", "*/*;q=0.8")
	requ.Header.Set("Connection", "close")
	requ.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36")
	resp, err := Client.Do(requ)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 赋值
	var requfofa *Fofadata
	json.Unmarshal(body, &requfofa)
	if requfofa.Err {
		return nil, errors.New("接口查询错误！")
	}
	return requfofa, nil
}

// 列出输出字段
func OutQuery(num int) {
	switch num {
	case 1:
		fmt.Println(WARNING + " 以下是搜索出的标签字段（参数title）：")
		fmt.Println(`+-----------------+--------------------------------------------------+----------+`)
		fmt.Println(`|     字段名      |                       描述                       |   权限   |`)
		fmt.Println(`+-----------------+--------------------------------------------------+----------+`)
		fmt.Println(`|     字段名      |                       描述                       |   权限   |`)
		fmt.Println(`|       ip        |                      ip地址                      |    无    |`)
		fmt.Println(`|      port       |                       端口                       |    无    |`)
		fmt.Println(`|    protocol     |                      协议名                      |    无    |`)
		fmt.Println(`|     country     |                     国家代码                     |    无    |`)
		fmt.Println(`|  country_name   |                      国家名                      |    无    |`)
		fmt.Println(`|     region      |                       区域                       |    无    |`)
		fmt.Println(`|      city       |                       城市                       |    无    |`)
		fmt.Println(`|    as_number    |                     asn编号                      |    无    |`)
		fmt.Println(`| as_organization |                     asn组织                      |    无    |`)
		fmt.Println(`|      host       |                      主机名                      |    无    |`)
		fmt.Println(`|     domain      |                       域名                       |    无    |`)
		fmt.Println(`|       os        |                     操作系统                     |    无    |`)
		fmt.Println(`|     server      |                    网站server                    |    无    |`)
		fmt.Println(`|       icp       |                    icp备案号                     |    无    |`)
		fmt.Println(`|      title      |                     网站标题                     |    无    |`)
		fmt.Println(`|     header      |                    网站header                    |    无    |`)
		fmt.Println(`|      cert       |                       证书                       |    无    |`)
		fmt.Println(`|      body       |                   网站正文内容                   | 企业会员 |`)
		fmt.Println(`|       fid       |                       fid                        | 企业会员 |`)
		fmt.Println(`|   structinfo    | 结构化信息 (部分协议支持、比如elastic、mongodb}) | 企业会员 |`)
		fmt.Println(`+-----------------+--------------------------------------------------+----------+`)
	case 2:
		fmt.Println(WARNING, "您输入的模式不存在（参数m），请使用下面其中一个模式：")
		fmt.Println("    - search	该功能是通过fofa搜索自定义资产（默认使用模式）")
		fmt.Println("    - getproxy	该功能是通过资产收集获取公开代理")
		fmt.Println("    - runproxy	该功能是读取自定义的代理文件进行代理")
		fmt.Println("    - outtitle	该功能是列出输出的结果标签")
	case 3:
		fmt.Println(WARNING, "可设置的地区如下（参数m）：")
		fmt.Println("    - abroad	该地区是使用国外代理（默认使用模式）")
		fmt.Println("    - china	该地区是国内代理")
		fmt.Println("    - all	该地区是随机地区代理")
	}
}

// fofa查询
func ProcesQuery(res *Fofadata) error {
	// 创建文件
	filePath = "fofascan_" + time.Now().Format("20060102_15_04_05") + ".csv"
	fw, err := os.Create(filePath)
	if err != nil {
		fmt.Println(ERR+"  创建文件错误 Error：", err.Error())
		return err
	}
	defer fw.Close()
	// 写入UTF-8 BOM
	fw.WriteString("\xEF\xBB\xBF")
	// 定义表格，用于格式终端输出
	title := strings.Split(Title, ",")
	// fmt.Println(title)
	table, err := gotable.Create(title...)
	if err != nil {
		fmt.Println(ERR+"  表格创建错误 Error：", err.Error())
		return err
	}
	// 创建一个新的写入文件流
	w := csv.NewWriter(fw)
	// 写入查询语句
	w.Write([]string{"查询语句：", res.Query})
	// 写入标题
	w.Write(title)
	// 遍历页面
	for _, datat := range res.Results {
		table.AddRow(datat)
	}
	// 写入多行数据
	w.WriteAll(res.Results)
	// 写入数据清除缓存
	w.Flush()
	// 输出
	fmt.Println("\r" + RED + "======================查询结果======================" + END)
	fmt.Print(table)
	fmt.Printf(RIGHT+" 搜索结果保存在"+ORANGE+" %s "+END+"中！", filePath)
	return nil
}
