package elm

import (
	"bytes"
	"log"
	"os"

	"golang.org/x/net/html"
)

//Doc struct
type Doc struct {
	root *Elm
}

//NewDoc func
func NewDoc(path string) *Doc {
	d := new(Doc)

	fs, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}

	node, err := html.Parse(fs)
	if err != nil {
		log.Panicln(err)
	}

	d.root = &Elm{node}

	return d
}

func (d *Doc) String() string {
	buf := new(bytes.Buffer)
	html.Render(buf, d.root.Node)

	return buf.String()
}

//Find 查找子元素
func (d *Doc) Find(query string) []*Elm {
	return d.root.Find(query)
}

//FindFirst 查找匹配的第一个子元素
func (d *Doc) FindFirst(query string) *Elm {
	return d.root.FindFirst(query)
}

//FindEach 查找子元素并输出
func (d *Doc) FindEach(query string, f func(node *Elm)) {
	d.root.FindEach(query, f)
}
