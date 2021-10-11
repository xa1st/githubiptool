package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/yinheli/qqwry"
)

const (
	GitApi    string = "https://api.github.com/meta"
	SiteUrl   string = "github.com"
	UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36"
)

var (
	// 协程管理
	wg sync.WaitGroup

	// 请求成功数量
	SuccNum int = 0
)

// 获取github iplist
func GetConf(info []byte) (res []string, err error) {
	json, err := simplejson.NewJson(info)
	if err != nil {
		fmt.Println(err, "远程信息格式不正确")
		return res, err
	}
	res = json.Get("web").MustStringArray()
	return res, err
}

// HTTP请求
func DownUrl(requrl string) (body []byte, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return body, err
	}
	// 添加头部协议
	request.Header.Add("User-Agent", UserAgent)
	response, err := client.Do(request)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return body, err
	}
	return body, err
}

func checkIp(addr string) (url string, match bool) {
	// 先把IP中带的掩码去掉
	taddr := strings.Split(addr, "/")
	if len(taddr) != 2 {
		return "", false
	}
	// 只验证IPV4
	match, err := regexp.MatchString(`^\d{0,3}\.\d{0,3}\.\d{0,3}\.\d{0,3}$`, taddr[0])
	if err != nil {
		return "", false
	}
	return taddr[0], match
}

func getPaddr(ip string) string {
	match, err := regexp.MatchString(`^\d{0,3}\.\d{0,3}\.\d{0,3}\.\d{0,3}$`, ip)
	if err != nil || !match {
		return ip
	}
	paddr := qqwry.NewQQwry("qqwry.dat")
	paddr.Find(ip)
	return paddr.Country
}

func TestUrl(saddr string, taddr string, port string) {
	// 结束时完成线程
	defer wg.Done()
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
					DualStack: true,
				}).DialContext(ctx, network, taddr+":"+port)
			},
		},
	}
	// 获取IP地址所在的信息
	paddr := getPaddr(taddr)
	t1 := time.Now()
	request, _ := http.NewRequest("GET", "https://"+saddr, nil)
	// 添加UA,防止网站拒绝访问
	request.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("%8s\t%s\t%s\t%6.2f秒\t%s\n", paddr, taddr, "https", time.Since(t1).Seconds(), err)
		return
	}
	defer resp.Body.Close()
	SuccNum++
	fmt.Printf("%8s\t%s\t%s\t%6.2f秒\t连接正常,响应代码:%s\n", paddr, taddr, "https", time.Since(t1).Seconds(), resp.Status)
}

func main() {
	fmt.Printf("开始获取所有github.com的公开IP地址...")
	body, err := DownUrl(GitApi)
	if err != nil {
		log.Panicln(err, "发生错误")
		os.Exit(0)
	}
	fmt.Printf("ok...")
	addrs, err := GetConf(body)
	if err != nil {
		log.Panicln(err, "发生错误!")
		os.Exit(0)
	}
	fmt.Printf("获取成功，共获取%d个有效ip地址，开始速度测试...\n", len(addrs))
	for _, addr := range addrs {
		if ip, ok := checkIp(addr); ok {
			wg.Add(1)
			go TestUrl(SiteUrl, ip, "443")
		}
	}
	wg.Wait()
	fmt.Printf("测试完成，所有的IP共%d个，可用的IP地址有%d个，欢迎下次使用...\n", len(addrs), SuccNum)
}
