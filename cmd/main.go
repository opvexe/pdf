package main

import (
	"flag"
	"fmt"
	"json-excle/core"
)

var (
	inPut  string
	outPut string
)

func init() {
	flag.StringVar(&inPut, "i", "", "请输入pdf文件名，例如 ../pdf/test.pdf")
	flag.StringVar(&outPut, "o", "", "请输入传出文件名,例如 ../pdf/out")
	flag.Parse()
	if len(inPut) == 0 || len(outPut) == 0 {
		flag.PrintDefaults()
	}
}

func main() {
	err := core.Generate(inPut, outPut)
	if err != nil {
		fmt.Println(err)
		return
	}
}
