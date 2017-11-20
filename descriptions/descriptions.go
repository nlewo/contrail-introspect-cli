package descriptions

import "fmt"
import "log"
import "github.com/jbowtie/gokogiri/xml"
import "github.com/gosuri/uitable"

import "github.com/nlewo/contrail-introspect-cli/utils"
import "github.com/nlewo/contrail-introspect-cli/collection"

func CtrlIfmap() collection.DescCollection {
	return collection.DescCollection{
		BaseXpath: "IFMapTableShowResp/ifmap_db/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "node_name/text()",
			LongDetail:       collection.LongFormatFn(ifmapNeighbors),
		},
		PageArgs: []string{"controller-fqdn", "table-name", "search-string"},
		PageBuilder: func(args []string) collection.Sourcer {
			path := fmt.Sprintf("Snh_IFMapTableShowReq?table_name=%s&search_string=%s", args[1], args[2])
			return collection.Webui{Path: path, VrouterUrl: args[0], Port: 8083}
		},
		PrimaryField: "node_name",
	}
}

func ifmapNeighbors(t *uitable.Table, e collection.Element) {
	t.AddRow("Node name", "Neighbors")
	nodeName, _ := e.Node.Search("node_name/text()")
	t.AddRow(fmt.Sprintf("%s", nodeName[0]))
	neighbors, _ := e.Node.Search("neighbors/list/element/text()")
	for _, n := range neighbors {
		t.AddRow("", fmt.Sprintf("%s", n))
	}
	t.AddRow("", "")
}

func Interface() collection.DescCollection {
	return collection.DescCollection{
		BaseXpath: "__ItfResp_list/ItfResp/itf_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"uuid", "name", "vrf_name", "vm_uuid", "ip_addr", "mdata_ip_addr"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{Table: "db.interface.0", VrouterUrl: args[0], Port: 8085}
		},
		PrimaryField: "name",
	}
}

func Peering() collection.DescCollection {
	return collection.DescCollection{
		PageArgs:  []string{"vrouter-fqdn"},
		BaseXpath: "AgentXmppConnectionStatus/peer/list",
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Webui{Path: "Snh_AgentXmppConnectionStatusReq", VrouterUrl: args[0], Port: 8085}
		},
		DescElt: collection.DescElement{
			ShortDetailXpath: "controller_ip/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"controller_ip", "state", "flap_count", "cfg_controller"}),
		},
		PrimaryField: "name",
	}
}

func Si() collection.DescCollection {
	return collection.DescCollection{
		BaseXpath: "__ServiceInstanceResp_list/ServiceInstanceResp/service_instance_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "uuid/text()",
			// LongDetailHelp: []string{"Service instance uuid", "Type of service instance", "Virtual machine uuid"},
			LongDetail: collection.LongFormatXpaths([]string{"uuid", "service_type", "instance_id"}),
		},
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{Table: "db.service-instance.0", VrouterUrl: args[0], Port: 8085}
		},
		PrimaryField: "uuid",
	}
}

func Route() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"vrouter-fqdn", "vrf-name"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{VrouterUrl: args[0], Table: args[1] + ".uc.route.0,", Port: 8085}
		},
		BaseXpath: "__Inet4UcRouteResp_list/Inet4UcRouteResp/route_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "src_ip/text()",
			LongDetail:       collection.LongFormatFn(routeDetail)},
		PrimaryField: "src_ip",
	}
}
func Vrf() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{Table: "db.vrf.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__VrfListResp_list/VrfListResp/vrf_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"name", "ucindex"}),
		},
		PrimaryField: "name",
	}
}
func Vn() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{Table: "db.vn.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__VnListResp_list/VnListResp/vn_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"name", "vrf_name"}),
		},
		PrimaryField: "name",
	}
}

func RiSummary() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"controller-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Webui{Path: "Snh_ShowRoutingInstanceSummaryReq?search_string=", VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRoutingInstanceSummaryResp/instances/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"name", "virtual_network"}),
		},
		PrimaryField: "name",
	}
}

func CtrlRouteSummary() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"controller-fqdn", "search"},
		PageBuilder: func(args []string) collection.Sourcer {
			path := fmt.Sprintf("Snh_ShowRouteSummaryReq?search_string=%s", args[1])
			return collection.Webui{Path: path, VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRouteSummaryResp/tables/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "name/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"name", "prefixes", "paths", "primary_paths", "secondary_paths", "pending_updates"}),
		},
		PrimaryField: "name",
	}
}

func CtrlRoute() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"controller-fqdn", "routing-instance"},
		PageBuilder: func(args []string) collection.Sourcer {
			path := fmt.Sprintf("Snh_ShowRouteReq?x=%s.inet.0", args[1])
			return collection.Webui{Path: path, VrouterUrl: args[0], Port: 8083}
		},
		BaseXpath: "ShowRouteResp/tables/list/ShowRouteTable/routes/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "prefix/text()",
			LongDetail:       collection.LongFormatFn(controllerRoutePath),
		},
		PrimaryField: "prefix",
	}
}

func Mpls() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"vrouter-fqdn"},
		PageBuilder: func(args []string) collection.Sourcer {
			return collection.Remote{Table: "db.mpls.0", VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__MplsResp_list/MplsResp/mpls_list/list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "label/text()",
			LongDetail:       collection.LongFormatFn(mplsDetail),
		},
		PrimaryField: "label",
	}
}

func routeDetail(t *uitable.Table, e collection.Element) {
	srcIp, _ := e.Node.Search("src_ip/text()")
	srcPrefix, _ := e.Node.Search("src_plen/text()")
	t.AddRow(fmt.Sprintf("Src %s/%s", srcIp[0], srcPrefix[0]))
	paths, _ := e.Node.Search("path_list/list/PathSandeshData")

	t.AddRow("    Dst", "Peers", "MPLS label", "Local Pref", "Interface", "Dest VN")
	for _, path := range paths {
		nhs, _ := path.Search("nh/NhSandeshData//dip/text()")
		peers, _ := path.Search("peer/text()")
		label, _ := path.Search("label/text()")
		pref, _ := path.Search("path_preference_data/PathPreferenceSandeshData/preference/text()")
		destvn, _ := path.Search("dest_vn/text()")
		itf, _ := path.Search("nh/NhSandeshData/itf/text()")
		t.AddRow("    "+utils.Pretty(nhs), utils.Pretty(peers), utils.Pretty(label),
			 utils.Pretty(pref), utils.Pretty(itf), utils.Pretty(destvn))
	}
	t.AddRow("")
}

func mplsDetail(t *uitable.Table, e collection.Element) {
	label, err := e.GetField("label")
	if err != nil {
		log.Fatal(err)
	}
	t.AddRow(fmt.Sprintf("Label: %s", label))
	nexthopDetail(t, e.Node)
	t.AddRow("")
}

func nexthopDetail(t *uitable.Table, node xml.Node) {
	t.AddRow("    Type", "Interface", "Nexthop index")
	nhs, _ := node.Search("nh/NhSandeshData/type/text()")
	itf, _ := node.Search("nh/NhSandeshData/itf/text()")
	nhIdx, _ := node.Search("nh/NhSandeshData/nh_index/text()")
	t.AddRow("    "+utils.Pretty(nhs), utils.Pretty(itf), utils.Pretty(nhIdx))
}

func controllerRoutePath(t *uitable.Table, e collection.Element) {
	srcIp, _ := e.Node.Search("prefix/text()")
	t.AddRow(fmt.Sprintf("Prefix %s", srcIp[0]))
	paths, _ := e.Node.Search("paths/list/ShowRoutePath")

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

func AgentPing() collection.DescCollection {
	return collection.DescCollection{
		PageArgs: []string{"vrouter-fqdn", "vrf-name", "source-ip", "source-port", "dest-ip", "dest-port"},
		PageBuilder: func(args []string) collection.Sourcer {
			sourceIp := args[2]
			sourcePort := args[3]
			destIp := args[4]
			destPort := args[5]
			vrfName := args[1]

			path := fmt.Sprintf("Snh_PingReq?source_ip=%s&source_port=%s&dest_ip=%s&dest_port=%s&protocol=6&vrf_name=%s&packet_size=&count=1&interval=",
				sourceIp, sourcePort, destIp, destPort, vrfName)
			return collection.Webui{Path: path, VrouterUrl: args[0], Port: 8085}
		},
		BaseXpath: "__PingResp_list",
		DescElt: collection.DescElement{
			ShortDetailXpath: "resp/text()",
			LongDetail:       collection.LongFormatXpaths([]string{"resp", "rtt"}),
		},
	}
}
