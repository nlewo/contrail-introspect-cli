package main

import "fmt"
import "os"
import "log"

import "github.com/moovweb/gokogiri/xpath"
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

func DescPeering() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn"},
		BaseXpath: "AgentXmppConnectionStatus/peer/list",
		PageBuilder: func(args []string) Sourcer {
			return Webui{Path: "Snh_AgentXmppConnectionStatusReq", VrouterUrl: args[0]}
		},
		DescElt: DescElement{
			ShortDetailXpath: "controller_ip/text()",
			LongDetail:       LongXpaths([]string{"controller_ip/text()", "state/text()", "flap_count/text()"}),
		},
	}
}


func DescItf() DescCol {
	return DescCol{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongXpaths([]string{"uuid/text()", "name/text()", "vrf_name/text()"}),
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
func DescRoute() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn", "vrf-name"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{VrouterUrl: args[0], Table: args[1] + ".uc.route.0,"}
		},
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "src_ip/text()",
			LongDetail:       LongFunc(routeDetail)},
		SearchAttribute: "source IP",
		SearchXpath: func(pattern string) string {
			return "RouteUcSandeshData/src_ip[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescVrf() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vrf.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__VrfListResp_list/VrfListResp/vrf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongXpaths([]string{"name/text()"}),
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "VrfSandeshData/name[contains(text(),'" + pattern + "')]/.."
		},
	}
}
func DescVn() DescCol {
	return DescCol{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vn.0", VrouterUrl: args[0]}
		},
		BaseXpath: "__VnListResp_list/VnListResp/vn_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongXpaths([]string{"name/text()", "vrf_name/text()"}),
		},
		SearchAttribute: "name",
		SearchXpath: func(pattern string) string {
			return "VnSandeshData/name[contains(text(),'" + pattern + "')]/.."
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

func main() {
	var count bool
	var hosts_file string

	app := cli.NewApp()
	app.Name = "contrail-introspect-cli"
	app.Usage = "CLI on ContraiL Introspects"
	app.Version = "0.0.2"
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
