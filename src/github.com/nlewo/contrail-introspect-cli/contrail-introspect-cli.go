package main

import "fmt"
import "os"
import "log"

import "github.com/moovweb/gokogiri/xpath"
import "github.com/moovweb/gokogiri/xml"
import "github.com/codegangsta/cli"
import "github.com/gosuri/uitable"

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

func DescPeering() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		BaseXpath: "AgentXmppConnectionStatus/peer/list",
		PageBuilder: func(args []string) Sourcer {
			return Webui{Path: "Snh_AgentXmppConnectionStatusReq", VrouterUrl: args[0]}
		},
		DescElt: DescElement{
			ShortDetailXpath: "controller_ip/text()",
			LongDetail:       LongFormatXpaths([]string{"controller_ip", "state", "flap_count"}),
		},
	}
}

func DescItf() DescCollection {
	return DescCollection{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"uuid", "name", "vrf_name", "vm_uuid"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.interface.0", VrouterUrl: args[0]}
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "ItfSandeshData/name[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescRoute() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn", "vrf-name"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{VrouterUrl: args[0], Table: args[1] + ".uc.route.0,"}
		},
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "src_ip/text()",
			LongDetail:       LongFormatFn(routeDetail)},
		SearchAttribute: "source IP",
		SearchXpath: func(pattern string) string {
			return "RouteUcSandeshData/src_ip[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescVrf() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vrf.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__VrfListResp_list/VrfListResp/vrf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name"}),
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "VrfSandeshData/name[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescVn() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vn.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__VnListResp_list/VnListResp/vn_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name", "vrf_name"}),
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "VnSandeshData/name[contains(text(),'" + pattern + "')]/.."
		},
	}
}

func DescMpls() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.mpls.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__MplsResp_list/MplsResp/mpls_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "label/text()",
			LongDetail:       LongFormatFn(mplsDetail),
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "MplsSandeshData/label[contains(text(),'" + pattern + "')]/.."
		},
	}
}

func routeDetail(e Element) {
	srcIp, _ := e.node.Search("src_ip/text()")
	fmt.Printf("Src %s\n", srcIp[0])
	paths, _ := e.node.Search("path_list/list/PathSandeshData")

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("    Dst", "Peers", "MPLS label", "Interface", "Dest VN")
	for _, path := range paths {
		nhs, _ := path.Search("nh/NhSandeshData//dip/text()")
		peers, _ := path.Search("peer/text()")
		label, _ := path.Search("label/text()")
		destvn, _ := path.Search("dest_vn/text()")
		itf, _ := path.Search("nh/NhSandeshData/itf/text()")
		table.AddRow("    "+Pretty(nhs), Pretty(peers), Pretty(label), Pretty(itf), Pretty(destvn))
	}
	fmt.Println(table)
}

func mplsDetail(e Element) {
	nexthopDetail(e.node)
}

func nexthopDetail(node xml.Node) {
        table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("    Type", "Interface", "Nexthop index")
	nhs, _ := node.Search("nh/NhSandeshData/type/text()")
	itf, _ := node.Search("nh/NhSandeshData/itf/text()")
	nhIdx, _ := node.Search("nh/NhSandeshData/nh_index/text()")
	table.AddRow("    "+Pretty(nhs), Pretty(itf), Pretty(nhIdx))
	fmt.Println(table)
}

func main() {
	var count bool
	var hosts_file string

	app := cli.NewApp()
	app.Name = "contrail-introspect-cli"
	app.Usage = "CLI on ContraiL Introspects"
	app.Version = "0.0.3"
	app.Before = func(c *cli.Context) error {
		if c.GlobalIsSet("hosts") {
			hosts = LoadHostsFile(c.GlobalString("hosts"))
			return nil
		}
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "hosts",
			Usage:       "host file to do DNS resolution",
			Destination: &hosts_file,
		}}
	app.Commands = []cli.Command{
		GenCommand(DescRoute(), "route", "Show routes"),
		GenCommand(DescItf(), "itf", "Show interfaces"),
		GenCommand(DescVrf(), "vrf", "Show vrfs"),
		GenCommand(DescPeering(), "peering", "Peering with controller"),
		GenCommand(DescVn(), "vn", "Show virtual network"),
		GenCommand(DescMpls(), "mpls", "Show mpls"),
		{
			Name:      "multiple",
			Usage:     "List routes with multiple nexthops",
			ArgsUsage: "vrouter vrf_name",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "count",
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
