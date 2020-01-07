package core

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

/*
	source :来源文件
	destance:目标源文件
*/
func Generate(source, destance string) error {
	s := strings.Split(source, "/")
	str := s[len(s)-1]
	hex := strings.Split(str, ".")[1]
	doc, err := fitz.New(source)
	if err != nil {
		return err
	}
	defer doc.Close()
	err = CreateOutPut(destance)
	if err != nil {
		return err
	}
	page := doc.NumPage() //获取pdf页面
	for i := 0; i < page; i++ {
		err := extractImages(doc, i, destance, hex)
		if err != nil {
			return err
		}
	}
	return nil
}

//生成文件创建
func  CreateOutPut(destance string) error {
	if _, err := os.Stat(destance); !os.IsNotExist(err) {
		return os.RemoveAll(destance)
	}
	fmt.Println("创建文件夹成功")
	return os.Mkdir(destance, os.ModePerm)
}

//生成image
func  extractImages(doc *fitz.Document, number int, destance string, hex string) error {
	img, err := doc.ImageDPI(number, 72)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s_%03d.png", hex, number)
	f, err := os.Create(filepath.Join(destance, name))
	defer f.Close()

	if err != nil {
		return err
	}

	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	fmt.Println("生成图片完成")
	return nil
}
