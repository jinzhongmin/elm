package main

import (
	"fmt"

	"git.oschina.net/jinzm/elm"
)

func main() {
	doc := elm.NewDoc("index.html")
	fmt.Println(doc.String())

	fmt.Println("\n\n\nexample 1. 批量处理")
	text := elm.NewElm(". ")
	doc.FindEach("span", func(node *elm.Elm) {
		node.After(text)
	})
	fmt.Println(doc.String())

	fmt.Println("\n\n\nexample 2. 从模板获取元素")
	root := elm.NewElmFromFile("tpl.html")
	alink := root.Find("a")[0]
	doc.FindFirst("div").Repace(alink.Copy())
	fmt.Println(doc.String())

	fmt.Println("\n\n\nexample 3. 查找字符串元素")
	doc.Find("$mark")[0].Repace(alink.Copy())
	fmt.Println(doc.String())

}
