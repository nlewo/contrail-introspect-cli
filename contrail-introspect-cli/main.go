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
		PageArgs:  []string{"vrouter-fqdn"},
		BaseXpath: "AgentXmppConnectionStatus/peer/list",
		PageBuilder: func(args []string) Sourcer {
			return Webui{Path: "Snh_AgentXmppConnectionStatusReq", VrouterUrl: args[0], Port: 8085}
		},
		DescElt: DescElement{
			ShortDetailXpath: "controller_ip/text()",
			LongDetail:       LongFormatXpaths([]string{"controller_ip", "state", "flap_count", "cfg_controller"}),
		},
		PrimaryField: "name",
	}
}

func DescItf() DescCollection {
	return DescCollection{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"uuid", "name", "vrf_name", "vm_uuid", "mdata_ip_addr"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.interface.0", VrouterUrl: args[0], Port: 8085}
		},
		PrimaryField: "name",
	}
}

func DescSi() DescCollection {
	return DescCollection{
		BaseXpath: "__ServiceInstanceResp_list/ServiceInstanceResp/service_instance_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "uuid/text()",
			// LongDetailHelp: []string{"Service instance uuid", "Type of service instance", "Virtual machine uuid"},
			LongDetail:       LongFormatXpaths([]string{"uuid", "service_type", "instance_id"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.service-instance.0", VrouterUrl: args[0], Port: 8085}
		},
		PrimaryField: "uuid",
	}
}

func DescRoute() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn", "vrf-name"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{VrouterUrl: args[0], Table: args[1] + ".uc.route.0,", Port: 8085}
		},
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "src_ip/text()",
			LongDetail:       LongFormatFn(routeDetail)},
		PrimaryField: "src_ip",
	}
}
func DescVrf() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vrf.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__VrfListResp_list/VrfListResp/vrf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name", "uc_index"}),
		},
		PrimaryField: "name",
	}
}
func DescVn() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.vn.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__VnListResp_list/VnListResp/vn_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name", "vrf_name"}),
		},
		PrimaryField: "name",
	}
}

func DescRiSummary() DescCollection {
	return DescCollection{
		PageArgs: []string{"controller-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Webui{Path: "Snh_ShowRoutingInstanceSummaryReq?search_string=", VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRoutingInstanceSummaryResp/instances/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name", "virtual_network"}),
		},
		PrimaryField: "name",
	}
}


func DescCtrlRouteSummary() DescCollection {
	return DescCollection{
		PageArgs: []string{"controller-fqdn", "search"},
		PageBuilder: func(args []string) Sourcer {
			path := fmt.Sprintf("Snh_ShowRouteSummaryReq?search_string=%s", args[1])
			return Webui{Path: path, VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRouteSummaryResp/tables/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"name", "prefixes", "paths", "primary_paths", "secondary_paths", "pending_updates"}),
		},
		PrimaryField: "name",
	}
}

func DescCtrlRoute() DescCollection {
	return DescCollection{
		PageArgs: []string{"controller-fqdn", "routing-instance"},
		PageBuilder: func(args []string) Sourcer {
			path := fmt.Sprintf("Snh_ShowRouteReq?x=%s.inet.0", args[1])
			return Webui{Path: path, VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRouteResp/tables/list/ShowRouteTable/routes/list",
		DescElt: DescElement{
			ShortDetailXpath: "prefix/text()",
			LongDetail:       LongFormatFn(controllerRoutePath),
		},
		PrimaryField: "prefix",
	}
}

func DescMpls() DescCollection {
	return DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.mpls.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__MplsResp_list/MplsResp/mpls_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "label/text()",
			LongDetail:       LongFormatFn(mplsDetail),
		},
		PrimaryField: "label",
	}
}


func routeSummaryDetail(e Element) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("Name", "Prefixes", "Paths", "Primary paths", "Secondary paths", "Pending Updates")
	fields := []string{"name", "prefixes", "paths", "primary_paths",
		"secondary_paths", "pending_updates"}
	paths, _ := e.node.Search(".")
	for _, path := range paths {
		values := [6]string{}
		for i, field := range fields {
			value, _ := path.Search(fmt.Sprintf("%s/text()", field))
			values[i] = Pretty(value)
		}
		table.AddRow(values[0], values[1], values[2], values[3], values[4], values[5])
	}
	fmt.Println(table)
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
	fmt.Printf("Label: %s\n", e.GetField("label"))
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

func controllerRoutePath(e Element) {
	srcIp, _ := e.node.Search("prefix/text()")
	fmt.Printf("Prefix %s\n", srcIp[0])
	paths, _ := e.node.Search("paths/list/ShowRoutePath")

	table := uitable.New()
	table.MaxColWidth = 80
	table.AddRow("    Protocol", "Nexthop", "Local Pref", "Peers", "MPLS label")
	for _, path := range paths {
		protocol, _ := path.Search("protocol/text()")
		nhs, _ := path.Search("next_hop/text()")
		peers, _ := path.Search("source/text()")
		label, _ := path.Search("label/text()")
		localPref, _ := path.Search("local_preference/text()")
		table.AddRow("    "+Pretty(protocol), Pretty(nhs), Pretty(localPref), Pretty(peers), Pretty(label))
	}
	fmt.Println(table)
}

func main() {
	var count bool
	var hosts_file string

	app := cli.NewApp()
	app.Name = "contrail-introspect-cli"
	app.Usage = "CLI on ContraiL Introspects"
	app.Version = "0.0.4"
	app.EnableBashCompletion = true
	app.Before = func(c *cli.Context) error {
		if c.GlobalIsSet("hosts") {
			var err error
			hosts, err = LoadHostsFile(c.GlobalString("hosts"))
			return err
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
		GenCommand(DescRoute(), "agent-route", "Show routes on agent"),
		GenCommand(DescItf(), "agent-itf", "Show interfaces on agent"),
		GenCommand(DescSi(), "agent-si", "Show service instances on agent"),
		GenCommand(DescVrf(), "agent-vrf", "Show vrfs on agent "),
		GenCommand(DescPeering(), "agent-peering", "Peering with controller on agent"),
		GenCommand(DescVn(), "agent-vn", "Show virtual networks on agent"),
		GenCommand(DescMpls(), "agent-mpls", "Show mpls on agent"),
		Follow(),
		Path(),
		GenCommand(DescRiSummary(), "controller-ri", "Show routing instances on controller"),
		GenCommand(DescCtrlRoute(), "controller-route", "Show routes on controller"),
		GenCommand(DescCtrlRouteSummary(), "controller-route-summary", "Show routes summary on controller"),
		{
			Name:      "agent-multiple",
			Usage:     "List routes with multiple nexthops",
			ArgsUsage: "vrouter vrf_name",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "count",
					Destination: &count,
				}},
			Action: func(c *cli.Context) error {
				if c.NArg() != 2 {
					log.Fatal("Wrong argument number!")
				}
				vrouter := c.Args()[0]
				vrf_name := c.Args()[1]
				multiple(vrouter, vrf_name, count)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
