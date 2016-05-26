// An introspect page is mapped into a Collection which is basically a
// list of Element.

package main

import "fmt"
import "log"

import "github.com/moovweb/gokogiri/xml"

// A Collection describes and contains a list of Elements.
type Collection struct {
	url      string
	descCol  DescCollection
	doc      *xml.XmlDocument
	node     xml.Node
	elements []Element
}

type Elements []Element

type Element struct {
	node xml.Node
	desc DescElement
}

// This contains informations to generate and query a Collection
type DescCollection struct {
	// Names of arguments required from the user to get datas from
	// introspect
	PageArgs    []string
	// A function that takes the list of arguments specified by
	// the user
	PageBuilder (func([]string) Sourcer)
	// The root Xpath
	BaseXpath   string
	// Description of Collection's Elements
	DescElt     DescElement
	// Name of the attribute used to search in the collection
	PrimaryField string
}

type DescElement struct {
	// Xpath used to generate the short version of an element
	ShortDetailXpath string
	// Used to generate the long version of an element
	LongDetail       LongFormatter
}

func (col *Collection) Init() {
	ss, _ := col.node.Search(col.descCol.BaseXpath + "/*")
	col.elements = make([]Element, len(ss))
	for i, s := range ss {
		col.elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
}


func (col *Collection) SearchXpathFuzzy(pattern string) string{
	return col.descCol.BaseXpath + "/*/" + col.descCol.PrimaryField + "[contains(text(),'" + pattern + "')]/.."
}

func (col *Collection) SearchXpathStrict(pattern string) string{
	return col.descCol.BaseXpath + "/*/" + col.descCol.PrimaryField + "[text()='" + pattern + "']/.."
}

func (col *Collection) SearchFuzzy(pattern string) Elements {
	return col.Search(col.SearchXpathFuzzy, pattern)
}

func (col *Collection) SearchStrict(pattern string) Elements {
	return col.Search(col.SearchXpathStrict, pattern)
}

func (col *Collection) Search(searchPredicate (func(string) string), pattern string) Elements {
	ss, _ := col.node.Search(searchPredicate(pattern))
	var elements []Element = make([]Element, len(ss))
	for i, s := range ss {
		elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
	return Elements(elements)
}

// Several representations of resources
type Shower interface {
	Long()
	Short()
	Xml()
}

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
		log.Fatal("Xpath '" + e.desc.ShortDetailXpath + "' is not valid (verify ShortDetailXpath)")
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
	e.desc.LongDetail.LongFormat(e)
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

// This is used to show the long version of an Element.
type LongFormatter interface {
	LongFormat(e Element)
}

type LongFormatFn (func(Element))
type LongFormatXpaths []string

func (fn LongFormatFn) LongFormat(e Element) {
	fn(e)
}
func (xpaths LongFormatXpaths) LongFormat(e Element) {
	for _, xpath := range xpaths {
		s, _ := e.node.Search(xpath + "/text()")
		if len(s) == 1 {
			fmt.Printf("%s ", Pretty(s))
		}
	}
}

