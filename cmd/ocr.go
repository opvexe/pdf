package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"json-excle/core"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	input string
	tid   string
	jde   string
)

const (
	sk = "0c406f80fb8362b14d9bd000c8403c754ba08509"
	ak = "0e4d9f01a1874e669d781ab620d566761b190abf"
)

func init() {
	flag.StringVar(&tid, "t", "2673e1b4-9589-48d0-a3d0-bf4939f28d83", "请输入模板id")
	flag.StringVar(&jde, "o", "a", "输出格式，默认是全部,j:json,e:excel")
	flag.StringVar(&input, "f", "", "输入文件名,例如: ./images")
	flag.Parse()
	if len(input) == 0 {
		flag.PrintDefaults()
	}
}

type Token struct {
	Access_token string `json:"access_token"`
}

func main() {
	err := core.CreateDir(core.SenseTime_JsonDir)
	if err != nil {
		fmt.Println("创建文件夹失败")
		return
	}
	err = core.CreateDir(core.SenseTime_PdfDir)
	if err != nil {
		fmt.Println("创建文件夹失败")
		return
	}
	err = core.CreateDir(core.SenseTime_InputDir)
	if err != nil {
		fmt.Println("创建文件夹失败")
		return
	}
	err = core.CreateDir(core.SenseTime_Excel)
	if err != nil {
		fmt.Println("创建文件夹失败")
		return
	}

	token := getToken()
	if len(token) != 0 || token !=""{

		fmt.Println("获取token成功:", token)
		var (
			ch          = make(chan *core.Resp, 1000)
			wg          sync.WaitGroup
			wgReceiving sync.WaitGroup
		)
		wgReceiving.Add(1)
		go core.ReceivingResults(ch, &wgReceiving, jde)
		//获取文件目录
		list, err := ioutil.ReadDir(input)
		if err != nil {
			fmt.Println("ReadDir:", err)
			return
		}
		for i := range list {
			fmt.Println(list[i].Name())
			path := path.Ext(list[i].Name())
			if path == ".jpg" || path == ".bmp" || path == ".png" || path == ".jpeg" {
				fmt.Println(input + "/" + list[i].Name())
				wg.Add(1)
				go core.SendHttp(ch, &wg, token, input+"/"+list[i].Name(), "image", list[i].Name(), tid)
			} else if path == ".pdf" {
				s := strings.Split(list[i].Name(), ".")[0]
				f := fmt.Sprintf("%s/%s/", core.SenseTime_InputDir, s) //生成文件目
				err = core.Generate(input+"/"+list[i].Name(), f)
				if err != nil {
					fmt.Println("Generate:", err)
					return
				}
				//打开目录
				pdir, err := ioutil.ReadDir(f)
				if err != nil {
					return
				}
				for j := range pdir {
					wg.Add(1)
					go core.SendHttp(ch, &wg, token, f+pdir[j].Name(), "pdf", pdir[j].Name(), tid)
				}
			}
		}
		wg.Wait()

		time.Sleep(1 * time.Millisecond)
		//关闭管道
		close(ch)

		wgReceiving.Wait()
	}
}

func getToken() string {
	url := fmt.Sprintf("https://ai-test.sensetime.com/user/acquireToken/?ak=%s&sk=%s", ak, sk)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求异常", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("获取token失败", err)
		return ""
	}
	return token.Access_token
}
