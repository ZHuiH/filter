package filter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//WordUnit 保存词语的基本结构
type WordUnit struct {
	IsEnd bool
	Next  map[rune]*WordUnit
}

//WordTree 查找根节点
var WordTree = make(map[rune]*WordUnit, 100)

//insertQueue 词语插入队列
var insertQueue = make(chan []rune, 10)

//Set 设置关键词
func Set(word string) bool {
	if len(word) <= 1 {
		fmt.Println("关键词不能少于1个单词！")
		return false
	}
	insertQueue <- slice(word)
	return true
}

//slice 将词语切分成单独的
func slice(word string) (array []rune) {
	word = strings.TrimSpace(word)
	array = []rune(word)
	return
}

//push 将关键词推进结构树
func push(word []rune) {
	var next map[rune]*WordUnit
	//验证头部
	if _, ok := WordTree[word[0]]; !ok {
		WordTree[word[0]] = createNode(word[1])
	}
	//截取头部
	next = WordTree[word[0]].Next
	target := word[1:]
	length := len(target) - 1
	//查询插入
	for index, item := range target {
		var nextWord rune = ' '

		if length > (index + 1) {
			nextWord = target[index+1]
		} else if length == index {
			if _, verify := next[item]; verify {
				next[item].IsEnd = true
			}
		}
		//插入节点&&去往下一个节点
		if _, ok := next[item]; !ok {
			next[item] = createNode(nextWord)
			next = next[item].Next
		} else {
			next = next[item].Next
		}
	}
}

//createNode 生成一个节点
func createNode(next rune) *WordUnit {
	unit := &WordUnit{
		IsEnd: false,
		Next:  make(map[rune]*WordUnit),
	}
	if next != ' ' {

		unit.Next[next] = &WordUnit{
			IsEnd: false,
			Next:  make(map[rune]*WordUnit),
		}
	} else {
		unit.IsEnd = true
	}
	return unit
}

//Insert 插入
func Insert() {
	var words []rune = make([]rune, 1)
	for {
		words = <-insertQueue
		push(words)
	}
}

//Search 查找
func Search(text string) (newText string) {
	//去掉空格 转换为[]rune
	target := slice(text)
	//初始化
	words := WordTree
	newText = ""
	verify := ""
	//查找
	for _, item := range target {
		//存在下一个节点
		if _, ok := words[item]; ok {
			verify += string(item)
			//当前的词语组合有一组匹配
			if words[item].IsEnd {
				verify = sign(verify)
			}
			words = words[item].Next
		} else {
			unit := string(item)
			//从根节点再找下有没有匹配的
			if _, again := WordTree[item]; again {
				verify += unit
				words = WordTree[item].Next
				unit = ""
			} else {
				words = WordTree
				newText += verify + unit
				verify = ""
			}
		}
	}
	//最后合并下
	newText += verify
	return
}

//sign 根据字符串的长度返回对应长度的*
func sign(str string) (sign string) {

	for range str {
		sign += "*"
	}
	return
}

//LoadConf 加载敏感词
func LoadConf(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件加载出错")
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("正在加载敏感词配置...")
	for scanner.Scan() {
		lineText := scanner.Text()
		Set(lineText)
	}
	fmt.Println("加载完成")
}
