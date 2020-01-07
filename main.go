package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"json-excle/core"
	"sync"
)

var (
	inPut  string
	outPut string
)

func init() {
	flag.StringVar(&inPut, "i", "", "请输入传入文件名，例如 ./test（不包含多级子目录）")
	flag.StringVar(&outPut, "o", "", "请输入传出文件名,例如 ./sensetime/out (out文件名需已创建)")
	flag.Parse()
	if len(inPut) == 0 || len(outPut) == 0 {
		flag.PrintDefaults()
	}
}

func main() {
	var (
		ch = make(chan *RecvieObj, 1000)
		wg sync.WaitGroup
	)
	wg.Add(1)
	go writeResult(ch, &wg)

	wg.Add(1)
	go readFile(ch, &wg)

	wg.Wait()
}

type RecvieObj struct {
	data []byte //数据
	fn   string //文件名
}

func writeResult(c <-chan *RecvieObj, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for oj := range c {
		obj, err := core.ParseJson(oj.data)
		if err != nil {
			fmt.Println("ParseJson:", err)
			return
		}
		ts, err := core.TransJson(obj)
		if err != nil {
			fmt.Println("TransJson:", err)
			return
		}
		err = core.Create(ts, outPut, oj.fn)
		if err != nil {
			return
		}
	}
}

func readFile(c chan<- *RecvieObj, wg *sync.WaitGroup) {
	defer func() {
		close(c)
		wg.Done()
	}()
	list, err := ioutil.ReadDir(inPut)
	if err != nil {
		fmt.Println("ReadDir:", err)
		return
	}
	for i := range list {
		fmt.Println(list[i].Name())
		b, err := ioutil.ReadFile(inPut + "/" + list[i].Name())
		if err != nil {
			fmt.Println("ReadFile:", err)
			return
		}
		obj := &RecvieObj{
			data: b,
			fn:   list[i].Name(),
		}
		c <- obj
	}
}
