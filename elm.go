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
	node *html.Node
}

//NewElm 从字符串创建*Elm
//必须是只有一个根的树形结构 否则只能获取到第一个根
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
//必须是只有一个根的树形结构 否则只能获取到第一个根
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
		attr := e.node.Attr
		for i := range attr {
			if attr[i].Key == key {
				return attr[i].Val
			}
		}

	} else if v, ok := val[0].(string); ok == true {
		//有val参数，就设置attr的值
		for i := range e.node.Attr {
			if e.node.Attr[i].Key == key {
				e.node.Attr[i].Val = v
				return ""
			}
		}

		//没有找到
		e.node.Attr = append(e.node.Attr, html.Attribute{Key: key, Val: v})

	} else {
		log.Panicln("Elm.Attr(key, val): val must be nil or string")
	}
	return ""
}

//Text 设置或返回text的值
func (e *Elm) Text(text ...string) string {
	if e.node.Type == textNode {
		if len(text) > 0 {
			e.node.Data = text[0]
		} else {
			return e.node.Data
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

//Parent 获取父元素
func (e *Elm) Parent() *Elm {
	if e.node.Parent != nil {
		p := new(Elm)
		p.node = e.node.Parent
		return p
	}
	return nil
}

//Child 获取子元素
//
func (e *Elm) Child(query ...interface{}) []*Elm {

	childs := make([]*Elm, 0)
	for child := e.node.FirstChild; child != nil; child = child.NextSibling {
		e := new(Elm)
		e.node = child
		childs = append(childs, e)
	}

	if len(query) > 0 {
		q, ok := query[0].(string)
		if ok {
			_childs := make([]*Elm, 0)
			for i := range childs {
				if elmTest(childs[i], q) {
					_childs = append(_childs, childs[i])
				}
			}

			return _childs
		}
	}

	return childs
}

//ChildAppend 添加子元素
func (e *Elm) ChildAppend(node *Elm) {
	e.node.AppendChild(nodeCopy(node.node))
}

//ChildRemove 移除子元素
//可以是选择器，可以是*Elm或[]*Elm，node为空时全部移除
func (e *Elm) ChildRemove(node ...interface{}) {
	if len(node) > 0 {
		switch node[0].(type) {
		case string:
			n, _ := node[0].(string)
			childs := e.Child(n)
			for i := range childs {
				e.node.RemoveChild(childs[i].node)
			}

		case *Elm:
			n, _ := node[0].(*Elm)
			e.node.RemoveChild(n.node)

		case []*Elm:
			n, _ := node[0].([]*Elm)
			for i := range n {
				e.node.RemoveChild(n[i].node)
			}
		default:
		}
	} else {
		for child := e.node.FirstChild; child != nil; {
			_child := child.NextSibling
			e.node.RemoveChild(child)
			child = _child
		}
	}

}

//Copy 复制元素元素
func (e *Elm) Copy() *Elm {
	node := nodeCopy(e.node)
	return &Elm{node}
}

//Repace 替换掉元素
func (e *Elm) Repace(node *Elm) {
	parent := e.node.Parent
	parent.InsertBefore(node.node, e.node)
	parent.RemoveChild(e.node)

	e.node = nil
	e = nil
}

//Remove 移除元素
func (e *Elm) Remove() {
	parent := e.node.Parent
	parent.RemoveChild(e.node)

	e.node = nil
	e = nil
}

//Before 在当前元素前插入元素
func (e *Elm) Before(node *Elm) {
	parent := e.node.Parent
	parent.InsertBefore(nodeCopy(node.node), e.node)

}

//After 在当前元素后插入元素
func (e *Elm) After(node *Elm) {
	parent := e.node.Parent

	//当前元素是否还有后一个兄弟元素
	if e.node.NextSibling != nil {
		parent.InsertBefore(nodeCopy(node.node), e.node.NextSibling)
	} else {
		parent.AppendChild(nodeCopy(node.node))
	}

}

//String 以文本输出
func (e *Elm) String() string {
	buf := new(bytes.Buffer)
	html.Render(buf, e.node)

	return buf.String()
}

//forEach 遍历elm本身、所有子、孙元素
func elmEach(e *Elm, f func(node *Elm)) {
	f(e)
	lastNode := e.node.FirstChild
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
			if e.node.Type == 1 && e.node.Data == q[1:] {
				continue
			} else {
				return false
			}
		} else {
			if e.node.Type == 3 && e.node.Data == q {
				continue
			} else {
				return false
			}
		}
	}

	return true
}

//nodeCopy 复制node
//ChildAppend 等操作时将改变添加的节点的一系列指针
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
