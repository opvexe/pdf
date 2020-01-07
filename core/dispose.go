package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Resp struct {
	data  []byte
	types string
	name  string
}

const SenseTime_JsonDir = "./senseTime_JsonDir"
const SenseTime_PdfDir = "./senseTime_PdfDir"
const SenseTime_InputDir = "./senseTime_InputDir"
const SenseTime_Excel = "./senseTime_Excel"

//接收回调数据
func ReceivingResults(ch <-chan *Resp, wg *sync.WaitGroup, jde string) {
	defer func() {
		wg.Done()
	}()
	for data := range ch {
		var obj Jsons
		err := json.Unmarshal(data.data, &obj)
		if err == nil && obj.Code == 1000 {
			s := strings.Split(data.name, ".")[0]
			n := SenseTime_JsonDir + "/" + s + ".json"
			b, _ := PathExists(n)
			if !b {
				fp, err := os.Create(n)
				if err != nil {
					fmt.Println("create:", err)
				}
				fp.Close()
			}
			f, err := os.OpenFile(n, os.O_WRONLY|os.O_TRUNC, 0600)
			if err != nil {
				fmt.Println("写入文件失败:", err)
			}
			f.WriteString(string(data.data))
			f.Close()

			ts, err := TransJson(&obj)
			if err != nil {
				fmt.Println("TransJson:", err)
				return
			}
			err = Create(ts, SenseTime_Excel, data.name)
			if err != nil {
				return
			}
		}
	}
}

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建文件夹
func CreateDir(des string) error {
	if _, err := os.Stat(des); !os.IsNotExist(err) {
		err := os.RemoveAll(des)
		if err != nil {
			return err
		}
	}
	return os.Mkdir(des, os.ModePerm)
}

//发起请求
func SendHttp(ch chan<- *Resp, wg *sync.WaitGroup, token string, fileName string, types string, name string,tid string) {
	defer func() {
		wg.Done()
	}()

	//发起请求
	target_url := fmt.Sprintf("https://ai-test.sensetime.com/ocr/tablematch?access_token=%s", token)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, _ := bodyWriter.CreateFormFile("image_file", fileName)
	fh, _ := os.Open(fileName)
	_, _ = io.Copy(fileWriter, fh)
	_ = bodyWriter.WriteField("tid",tid)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, _ := http.Post(target_url, contentType, bodyBuf)
	resp_body, _ := ioutil.ReadAll(resp.Body)
	reslut := &Resp{
		data:  resp_body,
		types: types,
		name:  name,
	}
	fmt.Println(string(resp_body))
	ch <- reslut
	return
}
