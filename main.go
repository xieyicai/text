package main

import (
	"bufio"
	"fmt"
	"github.com/xieyicai/text"
	"os"
)

func main() {
	var line, result, nn string
	var arr []rune
	var list text.Nums
	fmt.Printf("请输入一句中文。\n>")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line = input.Text()
		if line == "exit" {
			return
		}
		arr = []rune(line)
		list = *text.ReadNums(&arr)
		nn = ""
		for index, cn := range list {
			if index != 0 {
				nn += ", "
			}
			nn += cn.ToString()
		}
		result = text.Replace(line)
		fmt.Printf("提取了%d个数字（%s）。替换之后的字符串：%s\r\n\r\n>", len(list), nn, result)
	}
}
