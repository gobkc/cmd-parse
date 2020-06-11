package CmdParse

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	ExplainMain = iota
	ExplainItem
)

type CmdVal struct {
	K     string
	V     interface{}
	Def   interface{}
	Usage string
}

type Explain struct {
	Info string
	Type int
}

type MyParse struct {
	CurrentItem string
	ParseResult map[string]CmdVal
	Origin      []string
	ExplainRow  Explain
	Explains    []Explain
}

var Cli *MyParse
var cliOnce sync.Once

func New() *MyParse {
	cliOnce.Do(func() {
		Cli = new(MyParse)
		Cli.ParseAll()
		//始终创建一条help记录
		Cli.SetItem("-h").SetUsage("查看帮助").SetDefault(false).SaveSet()
	})
	return Cli
}

//保存所有命令行参数，并初始化解析结果
func (m *MyParse) ParseAll() *MyParse {
	m.Origin = os.Args
	//初始化
	m.ParseResult = make(map[string]CmdVal)
	return m
}

func (m *MyParse) SetItem(item string) *MyParse {
	if row, ok := m.ParseResult[item]; ok {
		log.Fatalln("重复的选项：", item)
	} else {
		row.K = item
		m.ParseResult[item] = row
	}
	m.CurrentItem = item
	return m
}

func (m *MyParse) SetDefault(v interface{}) *MyParse {
	if m.CurrentItem == "" {
		log.Fatalln("请先执行SetItem方法")
	}
	if row, ok := m.ParseResult[m.CurrentItem]; ok {
		row.Def = v
		m.ParseResult[m.CurrentItem] = row
	}
	return m
}

func (m *MyParse) SetUsage(usage string) *MyParse {
	if m.CurrentItem == "" {
		log.Fatalln("请先执行SetItem方法")
	}
	if row, ok := m.ParseResult[m.CurrentItem]; ok {
		row.Usage = usage
		m.ParseResult[m.CurrentItem] = row
	}
	return m
}

func (m *MyParse) SaveSet() CDataI {
	m.ParseCmd()
	cData := new(CData)
	cData.data = m.ParseResult[m.CurrentItem].V
	return cData
}

func (m *MyParse) Explain(title string) *MyParse {
	m.ExplainRow.Type = ExplainMain
	m.ExplainRow.Info = title
	m.Explains = append(m.Explains, m.ExplainRow)
	m.ExplainRow = Explain{}
	return m
}

func (m *MyParse) SetExplainItem(info string) *MyParse {
	m.ExplainRow.Type = ExplainItem
	m.ExplainRow.Info = info
	m.Explains = append(m.Explains, m.ExplainRow)
	m.ExplainRow = Explain{}
	return m
}

func (m *MyParse) ParseCmd() {
	//1.遍历origin 找出MAP中是否还有这个结果
	oriLen := len(m.Origin)
	for i, row := range m.Origin {
		if item, ok := m.ParseResult[row]; ok {
			if oriLen > i+1 {
				v := m.Origin[i+1]
				if !m.valueIsKey(v) {
					item.V = v
				}
				m.ParseResult[row] = item
			}
			//如果断言是BOOL类型则赋值为TRUE
			valType := fmt.Sprintf("%T", item.Def)
			if valType == "bool" {
				item.V = true
				m.ParseResult[row] = item
			}
		}
	}
}

func (m *MyParse) valueIsKey(value string) (isKey bool) {
	if _, ok := m.ParseResult[value]; ok {
		isKey = true
	}
	return isKey
}

func (m *MyParse) SaveExplain() bool {
	var hasExplain = false
	var eNum int
	if len(m.Explains) > 0 {
		hasExplain = true
	}
	if hasHelp, ok := m.ParseResult["-h"]; ok && hasHelp.V != nil && hasHelp.V.(bool) {
		//遍历ParseResult并且输出命令提示
		fmt.Fprintf(os.Stderr, "\n")
		for _, row := range m.ParseResult {
			repeatNum := 40 - len(row.K)
			if repeatNum < 0 {
				repeatNum = 5
			}
			row.K = row.K + strings.Repeat(" ", repeatNum)
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;35;40m%s%c[0m", 0x1B, row.K, 0x1B))
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;36;40m%s(默认：%v)%c[0m\n", 0x1B, row.Usage, row.Def, 0x1B))
		}
		for _, v := range m.Explains {
			if v.Type == ExplainMain {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;32;40m\n%s：%c[0m\n", 0x1B, v.Info, 0x1B))
				eNum = 1
			} else {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (%v) %s%c[0m\n", 0x1B, eNum, v.Info, 0x1B))
				eNum++
			}
		}
		os.Exit(0)
	}
	return hasExplain
}

type CDataI interface {
	GetInt() int
	GetString() string
	GetBool() bool
}

type CData struct {
	data interface{}
}

func (c *CData) GetInt() int {
	var result int
	if c.data != nil {
		data := c.data.(string)
		tmpInt, err := strconv.Atoi(data)
		if isNum := err == nil; isNum {
			result = tmpInt
		}
	}
	return result
}

func (c *CData) GetString() string {
	var result string
	if c.data != nil {
		result = c.data.(string)
	}
	return result
}

func (c *CData) GetBool() bool {
	var result bool
	if c.data != nil {
		dataType := fmt.Sprintf("%T", c.data)
		if dataType == "bool" {
			result = c.data.(bool)
		}
	}
	return result
}
