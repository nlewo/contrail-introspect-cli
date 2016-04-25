package main

import "fmt"
import "os"
import "net/http"
import "io/ioutil"
import "log"

import "github.com/moovweb/gokogiri"
import "github.com/moovweb/gokogiri/xml"
import "github.com/moovweb/gokogiri/xpath"
import "github.com/codegangsta/cli"


func load(url string, fromFile bool) *xml.XmlDocument{
	var data []byte
	if fromFile{
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
	return(doc)
}

func multiple(vrouter string, vrf_name string, count bool) {
	url := "http://" + vrouter + ":8085" + "/Snh_PageReq?x=begin:-1,end:-1,table:" + vrf_name + ".uc.route.0,"
	
	var doc = load(url, false)
	defer doc.Free()
	xps := xpath.Compile("//route_list/list/RouteUcSandeshData/path_list/list/PathSandeshData/nh/NhSandeshData/mc_list/../../../../../../src_ip/text()")
	ss, _ := doc.Root().Search(xps)
	if count {
		fmt.Printf("%d\n", len(ss))
	} else {
		for _, s := range ss {
			fmt.Printf("%s\n", s)
		}
	}
}

type Page struct {
	VrouterUrl string;
	Table string;
}

type File struct {
	Path string;
}

type LoadAble interface {
	Load(descCol DescCol) Collection;
}

// Parse data to XML
func fromDataToCollection(data []byte, descCol DescCol) Collection {
	doc, _ := gokogiri.ParseXml(data)
	ss, _ := doc.Search(descCol.BaseXpath)
	col := Collection{node: ss[0], descCol: descCol}
	col.Init()
	return col
}

func (file File) Load(descCol DescCol) Collection {
	f, _ := os.Open(file.Path)
	data, _ := ioutil.ReadAll(f)
	return fromDataToCollection(data, descCol)
}

func (page Page) Load(descCol DescCol) Collection {
	url := "http://" + page.VrouterUrl + ":8085/Snh_PageReq?x=begin:-1,end:-1,table:" + page.Table + ","
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return fromDataToCollection(data, descCol)
}

type Collection struct {
	descCol DescCol;
	doc *xml.XmlDocument;
	node xml.Node;
	elements []Element;
}

type DescCol struct {
	BaseXpath string;
	ShortDetailXpath string;
	LongDetailXpath []string;
	DescElt DescElement;
}

type DescElement struct {
	ShortDetailXpath string;
	LongDetailFunc (func (Element));
}

type Element struct {
	node xml.Node;
	desc DescElement;
}

func (col *Collection) Init() {
	ss, _ := col.node.Search("*")
	col.elements = make([]Element, len(ss))
	for i, s := range ss {
		col.elements[i] = Element{node: s, desc: col.descCol.DescElt}
	}
}

func (e Element) Short() {
	s, _ := e.node.Search(e.desc.ShortDetailXpath)
	if len(s) != 1 {
		log.Fatal("Xpath '" + e.desc.ShortDetailXpath + "' is not valid")
	}
	fmt.Printf("%s\n", s[0])
}

func (e Element) Long() {
	e.desc.LongDetailFunc(e)
}

func test() {
	col := BuildRouteFile("tmp.out").Load(DescRoute())
	col.Init()
	for _, e := range col.elements {
		e.Long()
	}
}

func (col Collection) Short() {
	for _, e := range col.elements {
		e.Short()
	}
}
func (col Collection) Long() {
	for _, e := range col.elements {
		e.Long()
	}
}

func BuildItfPage(vrouter string) Page {
	return Page{Table: "db.interface.0", VrouterUrl: vrouter}
}
func BuildVrfPage(vrouter string) Page {
	return Page{Table: "db.vrf.0", VrouterUrl: vrouter}
}
func BuildRoutePage(vrouter string, vrfName string) Page {
	return Page{Table: vrfName + ".uc.route.0,", VrouterUrl: vrouter}
}
func BuildRouteFile(path string) File {
	return File{Path: path}
}

func DescItf() DescCol {
	return DescCol{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		ShortDetailXpath: "ItfSandeshData/name/text()",
		DescElt: DescElement {
			ShortDetailXpath: "ItfSandeshData/name/text()"},
	}
}

func DescRoute() DescCol {
	return DescCol{
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		ShortDetailXpath: "RouteUcSandeshData/src_ip/text()",
		DescElt: DescElement {
			ShortDetailXpath: "src_ip/text()",
			LongDetailFunc: routeDetail},
	}
}

func DescVrf() DescCol {
	return DescCol{
		BaseXpath: "//vrf_list/list",
		ShortDetailXpath: "//name/text()",
	}
}

func routeGet(c Collection, srcIp string) xml.Node {
 	route, _ := c.node.Search("RouteUcSandeshData/src_ip[text()='" + srcIp + "']/..")
	if len(route) == 0 {
		log.Fatal("Route to " + srcIp + " was not found")
	}
	return route[0]
}

func routeDetail(e Element) {
	srcIp, _ := e.node.Search("src_ip/text()")
	fmt.Printf("%s\n", srcIp[0])
	paths, _ := e.node.Search("path_list/list/PathSandeshData")
	for _, path := range paths {
		nhs, _ := path.Search("nh/NhSandeshData//dip/text()")
		peers, _ := path.Search("peer/text()")
		label, _ := path.Search("label/text()")
		itf, _ := path.Search("nh/NhSandeshData/itf/text()")
		fmt.Printf("  %s %s %s %s\n", nhs, peers, label, itf)
	}
}

func main() {
	var vrouter string;
	var showAsXml bool;
	var count bool;

	app := cli.NewApp()
	app.Name = "contrail-introspect-cli"
	app.Usage = "CLI on contrail introspects"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "format-xml",
			Destination: &showAsXml,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "itf",
			Aliases:     []string{"a"},
			Usage:     "list interfaces",
			Flags: []cli.Flag{
				 cli.StringFlag{
					 Name: "vrouter",
					 Destination: &vrouter,
				 },
			},
			Action: func(c *cli.Context) {
				if c.NArg() != 1 {
					log.Fatal("Wrong argument number!")
				}
				vrouter := c.Args()[0]
				BuildItfPage(vrouter).Load(DescItf()).Short()
			},
		},
		{
			Name:      "multiple",
			Usage:     "vrouter vrf_name",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "count",
					Destination: &count,
				}},
			Action: func(c *cli.Context) {
				if c.NArg() != 2 {
					log.Fatal("Wrong argument number!")
				}
				vrouter := c.Args()[0]
				vrf_name := c.Args()[1]
				multiple(vrouter, vrf_name, count)
			},
		},
		{
			Name:      "vrf",
			Usage:     "vrf <vrouterUrl>",
			Action: func(c *cli.Context) {
				if c.NArg() != 1 {
					log.Fatal("Wrong argument number!")
				}
				BuildVrfPage(c.Args()[0]).Load(DescVrf()).Short()
			},
		},
		{
			Name:      "route",
			Usage:     "route <vrouterUrl> <vrfName> [<srcIp>]",
			Action: func(c *cli.Context) {
				col := BuildRoutePage(c.Args()[0], c.Args()[1]).Load(DescRoute())
				switch c.NArg() {
				case 2:
					if showAsXml {
						fmt.Printf("%s\n", col.node)
						return
					}
					col.Short()
				case 3:
					route := routeGet(col, c.Args()[2])
					if showAsXml {
						fmt.Printf("%s\n", route)
						return
					}
					routeDetail(Element{node: route})
				}
			},
		},
		{
			Name:      "route-from-file",
			Action: func(c *cli.Context) {
				col := BuildRouteFile(c.Args()[0]).Load(DescRoute())
				switch c.NArg() {
				case 1:
					if showAsXml {
						fmt.Printf("%s\n", col.node)
						return
					}
					col.Short()
				case 2:
					route := routeGet(col, c.Args()[1])
					if showAsXml {
						fmt.Printf("%s\n", route)
						return
					}
					routeDetail(Element{node: route})
				}
			},
		},
	}
	app.Run(os.Args)
}
