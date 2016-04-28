package main

import "fmt"
import "os"
import "net/http"
import "io/ioutil"
import "log"
import "strings"
import "bufio"
import "sort"

import "github.com/moovweb/gokogiri"
import "github.com/moovweb/gokogiri/xml"
import "github.com/moovweb/gokogiri/xpath"
import "github.com/codegangsta/cli"

// A map from IP to FQDN
type Hosts map[string]string

// Ok, it's a horrible hack... but I don't know yet how to propagated
// this variable from Arguments to Printf!
var hosts Hosts

// Take a hosts file in the same format than /etc/hosts file.
// Currently, the only two first elements are used.
func LoadHostsFile(filepath string) Hosts {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	m := make(Hosts)
	for scanner.Scan() {
		line := scanner.Text()
		ips := strings.Split(line, " ")
		m[ips[0]] = ips[1]
	}
	return m
}

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
func fromDataToCollection(data []byte, descCol DescCol, url string) Collection {
	doc, _ := gokogiri.ParseXml(data)
	ss, _ := doc.Search(descCol.BaseXpath)
	if len(ss) < 1 {
		log.Fatal(fmt.Sprintf("%d Failed to search xpath '%s'", len(ss), descCol.BaseXpath))
	}
	col := Collection{node: ss[0], descCol: descCol, url: url}
	col.Init()
	return col
}

func (file File) Load(descCol DescCol) Collection {
	f, _ := os.Open(file.Path)
	data, _ := ioutil.ReadAll(f)
	return fromDataToCollection(data, descCol, "file://" + file.Path)
}

func (page Page) Load(descCol DescCol) Collection {
	url := "http://" + page.VrouterUrl + ":8085/Snh_PageReq?x=begin:-1,end:-1,table:" + page.Table + ","
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return fromDataToCollection(data, descCol, url)
}

type Collection struct {
	url string;
	descCol DescCol;
	doc *xml.XmlDocument;
	node xml.Node;
	elements []Element;
}

type DescCol struct {
	PageArgs []string;
	PageBuilder (func([]string) LoadAble);
	BaseXpath string;
	DescElt DescElement;
	SearchXpath (func(string) string);
}

type DescElement struct {
	ShortDetailXpath string;
	LongDetail LongAble;
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

func (col *Collection) Search(pattern string) Elements{
	ss, _ := col.node.Search(col.descCol.SearchXpath(pattern))
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



func DescItf() DescCol {
	return DescCol{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: DescElement {
			ShortDetailXpath: "name/text()",
			LongDetail: LongXpaths([]string{"uuid/text()", "name/text()"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) LoadAble{
			return Page{Table: "db.interface.0", VrouterUrl: args[0]}
		},
	}
}
func DescRoute() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn", "vrf-name"},
		PageBuilder: func(args []string) LoadAble{
			return Page{VrouterUrl: args[0], Table: args[1] + ".uc.route.0,"}
		},
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		DescElt: DescElement {
			ShortDetailXpath: "src_ip/text()",
			LongDetail: LongFunc(routeDetail)},
		SearchXpath: func(pattern string) string {
			return "RouteUcSandeshData/src_ip[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescVrf() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) LoadAble{
			return Page{Table: "db.vrf.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__VrfListResp_list/VrfListResp/vrf_list/list",
		DescElt: DescElement {
			ShortDetailXpath: "name/text()",
			LongDetail: LongXpaths([]string{"name/text()"}),
		},
		SearchXpath: func(pattern string) string {
			return "VrfSandeshData/name[contains(text(),'" + pattern + "')]/.."
		},
	}
}

type LongFunc (func (Element))
type LongXpaths []string

type LongAble interface {
	Long(e Element);
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

func ResolveIp(ip string) string {
	fqdn, ok := hosts[ip]
	if ok {
		return fqdn
	} else {
		return ip
	}
}

func Pretty(nodes []xml.Node) string {
	ret := make([]string, len(nodes))
	for i, n := range nodes {
		ret[i] = ResolveIp(n.Content())
	}
	sort.Strings(ret)
	return strings.Join(ret, " ; ")
}

func routeDetail(e Element) {
	srcIp, _ := e.node.Search("src_ip/text()")
	fmt.Printf("%s\n", srcIp[0])
	paths, _ := e.node.Search("path_list/list/PathSandeshData")
	fmt.Printf("  Dest_ip ; Peers ; Label ; Itfs\n")
	for _, path := range paths {
		nhs, _ := path.Search("nh/NhSandeshData//dip/text()")
		peers, _ := path.Search("peer/text()")
		label, _ := path.Search("label/text()")
		itf, _ := path.Search("nh/NhSandeshData/itf/text()")
		fmt.Printf("  %s %s %s %s\n", Pretty(nhs), Pretty(peers), Pretty(label), Pretty(itf))
	}
}

func GenCommand(descCol DescCol, name string, usage string) cli.Command {
	return cli.Command{
		Name: name,
		Aliases: []string{"a"},
		Usage: usage,
		ArgsUsage: fmt.Sprintf("%s\n", strings.Join(descCol.PageArgs, " ")),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "long, l",
				Usage: "Long version format",
			},
			cli.BoolFlag{
				Name: "xml, x",
				Usage: "XML output format",
			},
			cli.BoolFlag{
				Name: "from-file",
				Usage: "Load file instead URL (for debugging)",
			},
			cli.BoolFlag{
				Name: "url, u",
				Usage: "Just show used URL",
			},
			cli.StringFlag{
				Name: "search, s",
				Usage: "Search pattern",
				Value: "",
			},
		},
		Action: func(c *cli.Context) {
			var page LoadAble;
			if c.IsSet("from-file") {
				page = File{Path: c.Args()[0]}
			} else {
				if c.NArg() < len(descCol.PageArgs) {
					cli.ShowSubcommandHelp(c)
					os.Exit(1)
				}
				page = descCol.PageBuilder(c.Args())
			}
			col := page.Load(descCol)
			if c.IsSet("url") {
				fmt.Println(col.url)
				return
			}

			var list Show;

			if c.String("s") != "" {
				list = col.Search(c.String("s"))
			} else {
				list = col
			}

			if c.IsSet("xml") {
				list.Xml()
				return
			}
			if c.IsSet("long") {
				list.Long()
				return
			}
			list.Short()
		},
	}
}

func main() {
	hosts = LoadHostsFile("hosts")

	var count bool;
	var hosts_file string;
	
	app := cli.NewApp()
	app.Name = "contrail-introspect-cli"
	app.Usage = "CLI on ContraiL Introspects"
	app.Version = "0.0.1"
	app.Before = func(c *cli.Context) error {
		if c.GlobalIsSet("hosts") {
			hosts = LoadHostsFile(c.GlobalString("hosts"))
			return nil
		}
		hosts = LoadHostsFile("hosts")
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "hosts",
			Usage: "host file to do DNS resolution",
			Destination: &hosts_file,
		}}
	app.Commands = []cli.Command{
		GenCommand(DescRoute(), "route", "Show routes"),
		GenCommand(DescItf(), "itf", "Show interfaces"),
		GenCommand(DescVrf(), "vrf", "Show vrfs"),
		{
			Name:      "multiple",
			Usage:     "List routes with multiple nexthops",
			ArgsUsage: "vrouter vrf_name",
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
	}
	app.Run(os.Args)
}
