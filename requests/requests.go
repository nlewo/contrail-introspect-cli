package requests

import "fmt"
import "log"
import "github.com/moovweb/gokogiri/xml"
import "github.com/gosuri/uitable"

import "github.com/nlewo/contrail-introspect-cli/utils"

func DescItf() DescCollection {
	return DescCollection{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       LongFormatXpaths([]string{"uuid", "name", "vrf_name", "vm_uuid", "ip_addr", "mdata_ip_addr"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) Sourcer {
			return Remote{Table: "db.interface.0", VrouterUrl: args[0], Port: 8085}
		},
		PrimaryField: "name",
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

func routeDetail(t *uitable.Table, e Element) {
	srcIp, _ := e.node.Search("src_ip/text()")
	t.AddRow(fmt.Sprintf("Src %s", srcIp[0]))
	paths, _ := e.node.Search("path_list/list/PathSandeshData")

	t.AddRow("    Dst", "Peers", "MPLS label", "Interface", "Dest VN")
	for _, path := range paths {
		nhs, _ := path.Search("nh/NhSandeshData//dip/text()")
		peers, _ := path.Search("peer/text()")
		label, _ := path.Search("label/text()")
		destvn, _ := path.Search("dest_vn/text()")
		itf, _ := path.Search("nh/NhSandeshData/itf/text()")
		t.AddRow("    "+utils.Pretty(nhs), utils.Pretty(peers), utils.Pretty(label), utils.Pretty(itf), utils.Pretty(destvn))
	}
	t.AddRow("")
}

func mplsDetail(t *uitable.Table, e Element) {
	label, err := e.GetField("label")
	if err != nil {
		log.Fatal(err)
	}
	t.AddRow(fmt.Sprintf("Label: %s", label))
	nexthopDetail(t, e.node)
	t.AddRow("")
}

func nexthopDetail(t *uitable.Table, node xml.Node) {
	t.AddRow("    Type", "Interface", "Nexthop index")
	nhs, _ := node.Search("nh/NhSandeshData/type/text()")
	itf, _ := node.Search("nh/NhSandeshData/itf/text()")
	nhIdx, _ := node.Search("nh/NhSandeshData/nh_index/text()")
	t.AddRow("    "+utils.Pretty(nhs), utils.Pretty(itf), utils.Pretty(nhIdx))
}

func controllerRoutePath(t *uitable.Table, e Element) {
	srcIp, _ := e.node.Search("prefix/text()")
	t.AddRow(fmt.Sprintf("Prefix %s", srcIp[0]))
	paths, _ := e.node.Search("paths/list/ShowRoutePath")

	t.AddRow("    Protocol", "Nexthop", "Local Pref", "Peers", "MPLS label")
	for _, path := range paths {
		protocol, _ := path.Search("protocol/text()")
		nhs, _ := path.Search("next_hop/text()")
		peers, _ := path.Search("source/text()")
		label, _ := path.Search("label/text()")
		localPref, _ := path.Search("local_preference/text()")
		t.AddRow("    "+utils.Pretty(protocol), utils.Pretty(nhs), utils.Pretty(localPref), utils.Pretty(peers), utils.Pretty(label))
	}
	t.AddRow("")
}
