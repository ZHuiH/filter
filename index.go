package main

import (
	"bufio"
	"filter"
	"flag"
	"fmt"
	"os"
)

var file string

func init() {
	flag.StringVar(&file, "c", "", "请使用 -c <file name> 设置敏感词")
}

func main() {
	//先开一个协程专门检查并插入词语
	go filter.Insert()
	flag.Parse()

	if file == "" {
		fmt.Println("请加载敏感词配置")
		os.Exit(0)
	}

	//加载敏感词配置
	filter.LoadConf(file)

	for {
		f := bufio.NewReader(os.Stdin)
		fmt.Print("请输入文本>")
		if text, err := f.ReadString('\n'); err == nil {
			fmt.Println(filter.Search(text))
		} else {
			fmt.Println("读取失败！")
		}
	}
}
