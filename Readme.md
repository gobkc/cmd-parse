# ginger
Package of command

> DEMO

    package main

    import (
        "fmt"
        "github.com/gobkc/cmd-parse"
    )
    var a = CmdParse.New().SetItem("a").SetUsage("aaa").SetDefault(1).SaveSet()
    var b = CmdParse.New().SetItem("b").SetUsage("bbb").SetDefault("bbb").SaveSet()
    var c = CmdParse.New().SetItem("c").SetUsage("ccc").SetDefault(false).SaveSet()
    var d = CmdParse.New().Explain("注意事项").
    SetExplainItem("本程序必须在网络联通环境下使用").
    SetExplainItem("本程序必须在必须配置数据库").
    SetExplainItem("本程序必须在必须配置XXX").
    Explain("使用示例").
    SetExplainItem("abc a aaa b bbb").
    SetExplainItem("abc a aaa2 b bbb2 c").
    SaveExplain()
    func main() {
        fmt.Println(a.GetInt(), b.GetString(),c.GetBool(),d)
    }

