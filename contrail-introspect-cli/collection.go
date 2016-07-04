// An introspect page is mapped into a Collection which is basically a
// list of Element.

package main

import "fmt"
import "log"

import "github.com/moovweb/gokogiri/xml"
import "github.com/gosuri/uitable"

// A Collection describes and contains a list of Elements.
type Collection struct {
	url     string
	descCol DescCollection
	// The node containing the whole XML, for instance the whole
	// loaded XML Page.
	rootNode xml.Node
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
	PageArgs []string
	// A function that takes the list of arguments specified by
	// the user
	PageBuilder (func([]string) Sourcer)
	// The root Xpath
	BaseXpath string
	// Description of Collection's Elements
	DescElt DescElement
	// Name of the attribute used to search in the collection
	PrimaryField string
}

type DescElement struct {
	// Xpath used to generate the short version of an element
	ShortDetailXpath string
	// Used to generate the long version of an element
	LongDetail LongFormatter
}

func (col *Collection) Init() {
	ss, _ := col.rootNode.Search(col.descCol.BaseXpath + "/*")
	col.elements = make([]Element, len(ss))
	for i, s := range ss {
		col.elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
}

func (col *Collection) SearchXpathFuzzy(pattern string) string {
	return col.descCol.BaseXpath + "/*/" + col.descCol.PrimaryField + "[contains(text(),'" + pattern + "')]/.."
}

func (col *Collection) SearchXpathStrict(pattern string) string {
	return col.descCol.BaseXpath + "/*/" + col.descCol.PrimaryField + "[text()='" + pattern + "']/.."
}

func (col *Collection) SearchFuzzy(pattern string) Elements {
	return col.Search(col.SearchXpathFuzzy, pattern)
}

func (col *Collection) SearchStrict(pattern string) Elements {
	return col.Search(col.SearchXpathStrict, pattern)
}

func (col *Collection) SearchFuzzyUnique(pattern string) Element {
	res := col.SearchFuzzy(pattern)
	if len(res) > 1 {
		fmt.Printf("Pattern %s matches:", pattern)
		for _, e := range res {
			fmt.Printf("\t%s", e)
		}
		log.Fatal("Pattern must match exactly one element")
	}
	return res[0]
}

func (col *Collection) Search(searchPredicate func(string) string, pattern string) Elements {
	ss, _ := col.rootNode.Search(searchPredicate(pattern))
	var elements []Element = make([]Element, len(ss))
	for i, s := range ss {
		elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
	return Elements(elements)
}

func (e Element) GetField(field string) string {
	s, _ := e.node.Search(fmt.Sprintf("%s/text()", field))
	for i, _ := range s {
		if s[i].String() != s[0].String() {
			log.Fatal(fmt.Sprintf("All fields values must be equal (values: %s)", s))
		}
	}
	if len(s) < 1 {
		log.Fatal(fmt.Sprintf("Field %s has not be found.", field))
	}
	return s[0].String()
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
	fmt.Printf("%s", c.rootNode)
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
	e.desc.LongDetail.LongFormat(FORMAT_TEXT, e)
}
func (col Collection) Long() {
	Elements(col.elements).Long()
}
func (elts Elements) Long() {
	for i, e := range elts {
		format := FORMAT_TABLE
		if i == 0 {
			format = FORMAT_TABLE_HEADER
		}
		e.desc.LongDetail.LongFormat(format, e)
		fmt.Printf("\n")
	}
}

// This is used to show the long version of an Element.
type LongFormatter interface {
	LongFormat(f Format, e Element)
}

type LongFormatFn (func(Element))
type LongFormatXpaths []string

type Format uint8

const (
	FORMAT_TEXT         Format = 1
	FORMAT_TABLE_HEADER Format = 2
	FORMAT_TABLE        Format = 3
)

func (fn LongFormatFn) LongFormat(format Format, e Element) {
	fn(e)
}

func (xpaths LongFormatXpaths) LongFormat(format Format, e Element) {
	if format == FORMAT_TABLE_HEADER || format == FORMAT_TABLE {
		longFormatTable(format, e, xpaths)
	} else {
		for _, xpath := range xpaths {
			s, _ := e.node.Search(xpath + "/text()")
			if len(s) == 1 {
				fmt.Printf("%s ", Pretty(s))
			}
		}
	}
}

func longFormatTable(format Format, e Element, xpaths LongFormatXpaths) {
	table := uitable.New()
	table.MaxColWidth = 80

	if format == FORMAT_TABLE_HEADER {
		tmp := make([]interface{}, len(xpaths))
		for i, v := range xpaths {
			tmp[i] = v
		}
		table.AddRow(tmp...)
	}

	tmp := make([]interface{}, len(xpaths))
	for i, xpath := range xpaths {
		s, _ := e.node.Search(xpath + "/text()")
		if len(s) == 1 {
			tmp[i] = Pretty(s)
		}
	}
	table.AddRow(tmp...)

	fmt.Print(table)
}
