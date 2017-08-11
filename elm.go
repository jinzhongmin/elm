package elm

import (
	"bytes"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//node类型
const (
	errorNode html.NodeType = iota
	textNode
	documentNode
	elementNode
	commentNode
	doctypeNode
	scopeMarkerNode
)

//Elm struct
type Elm struct {
	*html.Node
}

//NewElm 从字符串创建*Elm
//必须是只有一个源的树形结构 否则只能获取到第一个源
func NewElm(str string) *Elm {
	buf := bytes.NewBufferString(str)

	node, err := html.ParseFragment(buf, &html.Node{
		Type: elementNode,
	})
	if err != nil {
		log.Panicln(err)
	}

	return &Elm{node[0]}
}

//NewElmFromFile 从文件创建*Elm
//必须是只有一个源的树形结构 否则只能获取到第一个源
func NewElmFromFile(file string) *Elm {
	fs, err := os.Open(file)
	if err != nil {
		log.Panicln(err)
	}

	node, err := html.ParseFragment(fs, &html.Node{
		Type: elementNode,
	})
	if err != nil {
		log.Panicln(err)
	}

	return &Elm{node[0]}
}

//Attr 设置或返回attr的值
func (e *Elm) Attr(key string, val ...interface{}) string {

	if val == nil {
		//没有val参数，返回attr的值
		attr := e.Node.Attr
		for i := range attr {
			if attr[i].Key == key {
				return attr[i].Val
			}
		}

	} else if v, ok := val[0].(string); ok == true {
		//有val参数，就设置attr的值
		for i := range e.Node.Attr {
			if e.Node.Attr[i].Key == key {
				e.Node.Attr[i].Val = v
				return ""
			}
		}

		//没有找到
		e.Node.Attr = append(e.Node.Attr, html.Attribute{Key: key, Val: v})

	} else {
		log.Panicln("Elm.Attr(key, val): val must be nil or string")
	}
	return ""
}

//Text 设置或返回text的值
func (e *Elm) Text(text ...string) string {
	if e.Node.Type == textNode {
		if len(text) > 0 {
			e.Node.Data = text[0]
		} else {
			return e.Node.Data
		}
	}
	return ""
}

//Each 遍历elm所有子、孙元素
func (e *Elm) Each(f func(node *Elm)) {
	flag := false

	elmEach(e, func(node *Elm) {
		if flag == true {
			f(node)
		} else {
			flag = true
		}
	})
}

//Find 查找子元素
func (e *Elm) Find(query string) []*Elm {

	querySplit := strings.Split(query, " ")
	elmLits := make([]*Elm, 0)
	tmpLits := make([]*Elm, 0)

	elmLits = append(elmLits, e)
	for i := range querySplit {
		for ii := range elmLits {

			elmLits[ii].Each(func(node *Elm) {
				if elmTest(node, querySplit[i]) == true {
					tmpLits = append(tmpLits, node)
				}

			})

		}

		elmLits = append(make([]*Elm, 0), tmpLits...)
		tmpLits = make([]*Elm, 0)

	}

	return elmLits
}

//FindFirst 查找匹配的第一个子元素
func (e *Elm) FindFirst(query string) *Elm {

	querySplit := strings.Split(query, " ")
	elmLits := make([]*Elm, 0)
	tmpLits := make([]*Elm, 0)

	elmLits = append(elmLits, e)
	for i := range querySplit {
		for ii := range elmLits {

			elmLits[ii].Each(func(node *Elm) {
				if elmTest(node, querySplit[i]) == true {
					tmpLits = append(tmpLits, node)
				}

			})

		}

		elmLits = append(make([]*Elm, 0), tmpLits...)
		tmpLits = make([]*Elm, 0)

	}

	if len(elmLits) > 0 {
		return elmLits[0]
	}
	return nil
}

//FindEach 查找子元素并输出
func (e *Elm) FindEach(query string, f func(node *Elm)) {
	es := e.Find(query)
	for i := range es {
		f(es[i])
	}
}

//AppendChild 添加子元素
func (e *Elm) AppendChild(node *Elm) {
	e.Node.AppendChild(nodeCopy(node.Node))
}

//RemoveChild 移除子元素
func (e *Elm) RemoveChild(node *Elm) {
	e.Node.RemoveChild(node.Node)
	node.Node = nil
	node = nil
}

//Copy 复制元素元素
func (e *Elm) Copy() *Elm {
	node := nodeCopy(e.Node)
	return &Elm{node}
}

//Repace 替换掉元素
func (e *Elm) Repace(node *Elm) {
	parent := e.Node.Parent
	parent.InsertBefore(node.Node, e.Node)
	parent.RemoveChild(e.Node)

	e.Node = nil
	e = nil
}

//Remove 移除元素
func (e *Elm) Remove() {
	parent := e.Node.Parent
	parent.RemoveChild(e.Node)

	e.Node = nil
	e = nil
}

//InsertBefore 在当前元素前插入元素
func (e *Elm) InsertBefore(node *Elm) {
	parent := e.Node.Parent
	parent.InsertBefore(nodeCopy(node.Node), e.Node)

}

//InsertAfter 在当前元素后插入元素
func (e *Elm) InsertAfter(node *Elm) {
	parent := e.Node.Parent

	//当前元素是否还有后一个兄弟元素
	if e.Node.NextSibling != nil {
		parent.InsertBefore(nodeCopy(node.Node), e.Node.NextSibling)
	} else {
		parent.AppendChild(nodeCopy(node.Node))
	}

}

//forEach 遍历elm本身、所有子、孙元素
func elmEach(e *Elm, f func(node *Elm)) {
	f(e)
	lastNode := e.Node.FirstChild
	if lastNode == nil {
		return
	}

	for {
		elmEach(&Elm{lastNode}, f)
		if lastNode.NextSibling != nil {

			lastNode = lastNode.NextSibling
			continue
		} else {
			return
		}
	}
}

//test 测试是否符合query选择器
func elmTest(e *Elm, query string) bool {

	index := make([]int, 0)
	index = append(index, 0)
	for i := 1; i < len(query); i++ {
		if query[i] == '.' || query[i] == '#' {
			index = append(index, i)
		}
	}
	index = append(index, len(query))

	//复合选择器拆分成多个单个选择器，逐一测试
	for i := 1; i < len(index); i++ {
		start := index[i-1]
		end := index[i]

		q := query[start:end]
		if q[0] == '.' {
			flag := false

			class := e.Attr("class")
			classSplit := strings.Split(class, " ")

			for ii := range classSplit {
				if classSplit[ii] == q[1:] {
					flag = true
					break
				}
			}

			if flag == false {
				return false
			}
		} else if q[0] == '#' {
			id := e.Attr("id")

			if id == q[1:] {
				continue
			} else {
				return false
			}

		} else if q[0] == '$' {
			if e.Node.Type == 1 && e.Node.Data == q[1:] {
				continue
			} else {
				return false
			}
		} else {
			if e.Node.Type == 3 && e.Node.Data == q {
				continue
			} else {
				return false
			}
		}
	}

	return true
}

//nodeCopy 复制node
//AppendChild 等操作时将改变添加的节点的一系列指针
//会造成指针混乱，因此需要复制节点
func nodeCopy(src *html.Node) *html.Node {
	dst := new(html.Node)

	dst.Type = src.Type
	dst.DataAtom = src.DataAtom
	dst.Data = src.Data
	dst.Namespace = src.Namespace
	dst.Attr = append(dst.Attr, src.Attr...)

	si := src.FirstChild
	if si != nil {
		for {
			dst.AppendChild(nodeCopy(si))
			si = si.NextSibling

			if si != nil {
				continue
			} else {
				break
			}
		}
	}

	return dst
}
