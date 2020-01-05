package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
	"time"
)

// 定义json
type Jsons struct {
	Code int        `json:"code"`
	Data []JsonData `json:"data"`
}

//定义data
type JsonData struct {
	Objects []JsonObject `json:"objects"`
}

//定义object
type JsonObject struct {
	Type        string      `json:"type"`
	Areas       []JsonAreas `json:"areas"`
	Rectify_mat []float32   `json:"rectify_mat"`
}

//定义areas
type JsonAreas struct {
	Name  string      `json:"name"`
	Valid bool        `json:"valid"`
	Texts []JsonTexts `json:"texts"`
}

//定义texts
type JsonTexts struct {
	Name      string         `json:"name"`
	Subtexts  []JsonSubtexts `json:"subtexts"`
	Roi       []JsonRoi      `json:"roi"`
	Valid     bool           `json:"valid"`
	Multiline bool           `json:"multiline"`
}

//定义roi
type JsonRoi struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

//定义subtexts
type JsonSubtexts struct {
	Content string    `json:"content"`
	Roi     []JsonRoi `json:"roi"`
	Score   float64   `json:"score"`
	Valid   bool      `json:"valid"`
}

//解析数据
func ParseJson(j []byte) (*Jsons, error) {
	var jsonReslut Jsons
	err := json.Unmarshal(j, &jsonReslut)
	if err != nil {
		return nil, err
	}
	return &jsonReslut, nil
}

//表名数据
type TransData struct {
	Data []map[string][][]string
}

//转换数据 Excle表格想要的数据
func TransJson(obj *Jsons) (*TransData, error) {
	tsd := new(TransData)
	if len(obj.Data) == 0 {
		return nil, errors.New("data数据为空")
	}
	//获取data
	d := obj.Data[0]
	//获取object
	if len(d.Objects) == 0 {
		return nil, errors.New("object数据为空")
	}
	oj := d.Objects[0]
	//获取areas
	if len(oj.Areas) == 0 {
		return nil, errors.New("areas数据为空")
	}
	var tbc []map[string][][]string
	//获取area对象
	for _, v := range oj.Areas {
		tmap := make(map[string][][]string)
		//获取表名
		tb := v.Name
		//二维数组数据
		var s [][]string
		//获取行表体
		var h []string
		for _, hv := range v.Texts {
			if strings.HasPrefix(hv.Name, "item-0-") {
				im := strings.Split(hv.Name, "item-0-")[1]
				h = append(h, im)
			}
		}
		//第一列数据
		s = append(s, h)
		//获取列数据
		var (
			kc  []string
			dic = make(map[string]int)
		)
		for _, hv := range v.Texts {
			s1 := strings.Split(hv.Name, "-")[1]
			n, _ := strconv.ParseInt(s1, 0, 64)
			_, ok := dic[s1]
			if !ok {
				if kc != nil {
					s = append(s, kc)
					kc = nil
				}
				var c string = ""
				for k, subv := range hv.Subtexts {
					if k == 0 {
						c = subv.Content
					} else {
						if len(c) == 0 || c == "" {
							c = subv.Content
						} else {
							c = c + "/" + subv.Content
						}
					}
				}
				kc = append(kc, c)
				dic[s1] = int(n)
			} else {
				var c string = ""
				for k, subv := range hv.Subtexts {
					if k == 0 {
						c = subv.Content
					} else {
						if len(c) == 0 || c == "" {
							c = subv.Content
						} else {
							c = c + "/" + subv.Content
						}
					}
				}
				kc = append(kc, c)
			}
		}
		tmap[tb] = s
		//追加table
		tbc = append(tbc, tmap)
	}
	tsd.Data = tbc
	return tsd, nil
}

//创建Excle表格
func Create(ts *TransData, file string, name string) error {
	for k, v := range ts.Data {
		//根据key获取表名
		tb := fmt.Sprintf("table-%d", k+1)
		//获取表数据
		s := v[tb]
		if len(s) > 0 {
			//创建表格
			f := xlsx.NewFile()
			//创建表格标题
			sheet, err := f.AddSheet("SenseTime-" + tb)
			if err != nil {
				return err
			}
			for _, r := range s {
				row := sheet.AddRow()
				row.SetHeightCM(2)
				for _, h := range r {
					rn := row.AddCell()
					rn.Value = h
				}
			}
			//获取当前时间戳
			num := strconv.Itoa(k)
			n := strings.Split(name,".")[0] +"_" + time.Now().Format("20060102")+"_"+ num
			err = f.Save(file + "/" + n + ".xlsx")
			if err != nil {
				return err
			}
			fmt.Println("Excel表格保存成功", n)
		}
	}
	return nil
}
