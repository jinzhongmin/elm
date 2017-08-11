# elm 
简单的Golang Template
## example
``` go
package main

import (
	"net/http"

	"git.oschina.net/jinzm/elm"

	"github.com/labstack/echo"
)

//使用echo框架时的例子
func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		doc := elm.NewDoc("index.html")
        //...         你可以在此处理网页文件的内容
		return c.HTML(http.StatusOK, doc.String())
	})

	e.Start(":80")

}
```
## api
``` go
/*************************\\
||       struct DOC       ||
\\************************/
//读取html文件
func NewDoc (file string) Doc
//输出html内容
func (d *Doc) String() string

//Find 查找子元素
func (d *Doc) Find(query string) []*Elm
//FindFirst 查找匹配的第一个子元素
func (d *Doc) FindFirst(query string) *Elm 
//FindEach 查找子元素并输出
func (d *Doc) FindEach(query string, f func(node *Elm))


/*************************\\
||       struct Elm       ||
\\************************/
//NewElm 从字符串创建*Elm
//必须是只有一个源的树形结构 否则只能获取到第一个源
func NewElm(str string) *Elm 
//NewElmFromFile 从文件创建*Elm
//必须是只有一个源的树形结构 否则只能获取到第一个源
func NewElmFromFile(file string) *Elm
//Attr 设置或返回attr的值
func (e *Elm) Attr(key string, val ...interface{}) string 
//Text 设置或返回text的值
func (e *Elm) Text(text ...string) string 
//Each 遍历elm所有子、孙元素
func (e *Elm) Each(f func(node *Elm))
//Find 查找子元素
func (e *Elm) Find(query string) []*Elm
//FindFirst 查找匹配的第一个子元素
func (e *Elm) FindFirst(query string) *Elm
//FindEach 查找子元素并输出
func (e *Elm) FindEach(query string, f func(node *Elm))
//AppendChild 添加子元素
func (e *Elm) AppendChild(node *Elm)
//RemoveChild 移除子元素
func (e *Elm) RemoveChild(node *Elm)
//Copy 复制元素元素
func (e *Elm) Copy() *Elm
//Repace 替换掉元素
func (e *Elm) Repace(node *Elm)
//Remove 移除元素
func (e *Elm) Remove()
//InsertBefore 在当前元素前插入元素
func (e *Elm) InsertBefore(node *Elm)
//InsertAfter 在当前元素后插入元素
func (e *Elm) InsertAfter(node *Elm) 
```