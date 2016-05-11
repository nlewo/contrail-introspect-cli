package main

import "fmt"
import "log"

import "github.com/moovweb/gokogiri/xml"

type Collection struct {
	url      string
	descCol  DescCol
	doc      *xml.XmlDocument
	node     xml.Node
	elements []Element
}

type DescCol struct {
	PageArgs    []string
	PageBuilder (func([]string) Sourcer)
	BaseXpath   string
	DescElt     DescElement
	SearchAttribute string
	SearchXpath (func(string) string)
}

type DescElement struct {
	ShortDetailXpath string
	LongDetail       LongAble
}

type Element struct {
	node xml.Node
	desc DescElement
}

func (col *Collection) Init() {
	ss, _ := col.node.Search(col.descCol.BaseXpath + "/*")
	col.elements = make([]Element, len(ss))
	for i, s := range ss {
		col.elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
}

func (col *Collection) Search(pattern string) Elements {
	ss, _ := col.node.Search(col.descCol.BaseXpath + "/" + col.descCol.SearchXpath(pattern))
	var elements []Element = make([]Element, len(ss))
	for i, s := range ss {
		elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
	return Elements(elements)
}

type Show interface {
	Long()
	Short()
	Xml()
}

type Elements []Element

func (e Element) Xml() {
	fmt.Printf("%s", e.node)
}
func (elts Elements) Xml() {
	for _, e := range elts {
		e.Xml()
	}
}
func (c Collection) Xml() {
	fmt.Printf("%s", c.node)
}

func (e Element) Short() {
	s, _ := e.node.Search(e.desc.ShortDetailXpath)
	if len(s) != 1 {
		log.Fatal("Xpath '" + e.desc.ShortDetailXpath + "' is not valid")
	}
	fmt.Printf("%s\n", s[0])
}
func (col Collection) Short() {
	Elements(col.elements).Short()
}
func (elts Elements) Short() {
	for _, e := range elts {
		e.Short()
	}
}
func (e Element) Long() {
	e.desc.LongDetail.Long(e)
}
func (col Collection) Long() {
	Elements(col.elements).Long()
}
func (elts Elements) Long() {
	for _, e := range elts {
		e.Long()
		fmt.Printf("\n")
	}
}

type LongFunc (func(Element))
type LongXpaths []string

type LongAble interface {
	Long(e Element)
}

func (lf LongFunc) Long(e Element) {
	lf(e)
}
func (xpaths LongXpaths) Long(e Element) {
	for _, xpath := range xpaths {
		s, _ := e.node.Search(xpath)
		fmt.Printf("%s ", s[0])
	}
}

