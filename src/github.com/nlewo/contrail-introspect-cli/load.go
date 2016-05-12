package main

import "fmt"
import "os"
import "net/http"
import "io/ioutil"
import "log"

import "github.com/moovweb/gokogiri"
import "github.com/moovweb/gokogiri/xml"

type Remote struct {
	VrouterUrl string
	Table      string
}

type Webui struct {
	VrouterUrl string
	Path      string
}

type File struct {
	Path string
}

type Sourcer interface {
	Load(descCol DescCollection) Collection
}

func load(url string, fromFile bool) *xml.XmlDocument {
	var data []byte
	if fromFile {
		file, _ := os.Open(url)
		data, _ = ioutil.ReadAll(file)
	} else {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		data, _ = ioutil.ReadAll(resp.Body)
	}
	doc, _ := gokogiri.ParseXml(data)
	return (doc)
}

// Parse data to XML
func fromDataToCollection(data []byte, descCol DescCollection, url string) Collection {
	doc, _ := gokogiri.ParseXml(data)
	ss, _ := doc.Search("/")
	if len(ss) < 1 {
		log.Fatal(fmt.Sprintf("%d Failed to search xpath '%s'", len(ss), descCol.BaseXpath))
	}
	col := Collection{node: ss[0], descCol: descCol, url: url}
	col.Init()
	return col
}

func (file File) Load(descCol DescCollection) Collection {
	f, _ := os.Open(file.Path)
	data, _ := ioutil.ReadAll(f)
	return fromDataToCollection(data, descCol, "file://"+file.Path)
}

func (page Remote) Load(descCol DescCollection) Collection {
	url := "http://" + page.VrouterUrl + ":8085/Snh_PageReq?x=begin:-1,end:-1,table:" + page.Table + ","
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return fromDataToCollection(data, descCol, url)
}

func (page Webui) Load(descCol DescCollection) Collection {
	url := "http://" + page.VrouterUrl + ":8085/" + page.Path
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return fromDataToCollection(data, descCol, url)
}
