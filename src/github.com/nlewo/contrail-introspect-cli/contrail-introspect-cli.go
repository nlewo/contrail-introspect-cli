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
		resp, _ := http.Get(url)
		data, _ = ioutil.ReadAll(resp.Body)
	}
	doc, _ := gokogiri.ParseXml(data)
	return(doc)
}

func itf(vrouter string) {
	url := "http://" + vrouter + ":8085/Snh_ItfReq"
	var doc = load(url, false)
	xps := xpath.Compile("//name/text()")
	ss, _ := doc.Root().Search(xps)

	for _, s := range ss {
		fmt.Printf("%s\n", s)
	}
}

func multiple(vrouter string, vrf_name string) {
	url := "http://" + vrouter + ":8085" + "/Snh_PageReq?x=begin:-1,end:-1,table:" + vrf_name + ".uc.route.0,"
	
	var doc = load(url, false)
	defer doc.Free()
	xps := xpath.Compile("//route_list/list/RouteUcSandeshData/path_list/list/PathSandeshData/nh/NhSandeshData/mc_list/../../../../../../src_ip/text()")
	ss, _ := doc.Root().Search(xps)
	for _, s := range ss {
		fmt.Printf("%s\n", s)
	}
}

func vrf(vrouter string) {
	var url = "http://" + vrouter + ":8085" + "/Snh_VrfListReq"
	var doc = load(url, false)
	defer doc.Free()
	xps := xpath.Compile("//vrf_list/list//name/text()")
	ss, _ := doc.Root().Search(xps)
	for _, s := range ss {
		fmt.Printf("%s\n", s)
	}
}

func routeFromFile(filePath string) Collection {
	var doc = load(filePath, true)
	return route(doc)
}

func routeFromUrl(vrouter string, vrfName string) Collection {
	var url = "http://" + vrouter + ":8085" + "/Snh_PageReq?x=begin:-1,end:-1,table:" + vrfName + ".uc.route.0,"
	var doc = load(url, false)
	return route(doc)
}


func route(doc *xml.XmlDocument) Collection {
	ss, err := doc.Root().Search("/__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list")
	if err != nil {
		log.Fatal(err)
	}
	col := Collection{node: ss[0]}
	return col
}


type Collection struct {
	doc *xml.XmlDocument
	node xml.Node
}

func routeList(col Collection) {
	ss, _ := col.node.Search("RouteUcSandeshData/src_ip/text()")	
	for _, s := range ss {
		fmt.Printf("%s\n", s)
	}
}

func routeGet(c Collection, srcIp string) xml.Node {
 	route, _ := c.node.Search("RouteUcSandeshData/src_ip[text()='" + srcIp + "']/..")
	return route[0]
}

func routeDetail(n xml.Node) {
	srcIp, _ := n.Search("src_ip/text()")
	fmt.Printf("%s\n", srcIp[0])
	paths, _ := n.Search("path_list/list/PathSandeshData")
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
				itf(vrouter)
			},
		},
		{
			Name:      "multiple",
			Usage:     "vrouter vrf_name",
			Action: func(c *cli.Context) {
				if c.NArg() != 2 {
					log.Fatal("Wrong argument number!")
				}
				vrouter := c.Args()[0]
				vrf_name := c.Args()[1]
				multiple(vrouter, vrf_name)
			},
		},
		{
			Name:      "vrf",
			Usage:     "vrf <vrouterUrl>",
			Action: func(c *cli.Context) {
				if c.NArg() != 1 {
					log.Fatal("Wrong argument number!")
				}
				vrouter := c.Args()[0]
				vrf(vrouter)
			},
		},
		{
			Name:      "route",
			Usage:     "route <vrouterUrl> <vrfName> [<srcIp>]",
			Action: func(c *cli.Context) {
				col := routeFromUrl(c.Args()[0], c.Args()[1])
				switch c.NArg() {
				case 2:
					if showAsXml {
						fmt.Printf("%s\n", col.node)
						return
					}
					routeList(col)
				case 3:
					route := routeGet(col, c.Args()[2])
					if showAsXml {
						fmt.Printf("%s\n", route)
						return
					}
					routeDetail(route)
				}
			},
		},
		{
			Name:      "route-from-file",
			Action: func(c *cli.Context) {
				col := routeFromFile(c.Args()[0])
				switch c.NArg() {
				case 1:
					if showAsXml {
						fmt.Printf("%s\n", col.node)
						return
					}
					routeList(col)
				case 2:
					route := routeGet(col, c.Args()[1])
					if showAsXml {
						fmt.Printf("%s\n", route)
						return
					}
					routeDetail(route)
				}
			},
		},
	}
	app.Run(os.Args)
}
