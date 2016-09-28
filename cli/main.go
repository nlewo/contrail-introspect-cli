package main

import "fmt"
import "os"
import "log"

import "github.com/moovweb/gokogiri/xpath"
import "github.com/codegangsta/cli"

import "github.com/nlewo/contrail-introspect-cli/requests"
import "github.com/nlewo/contrail-introspect-cli/utils"

func multiple(vrouter string, vrf_name string, count bool) {
	url := "http://" + vrouter + ":8085" + "/Snh_PageReq?x=begin:-1,end:-1,table:" + vrf_name + ".uc.route.0,"

	var doc = requests.Load(url, false)
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
			utils.HostMap, err = utils.LoadHostsFile(c.GlobalString("hosts"))
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
		GenCommand(requests.DescRoute(), "agent-route", "Show routes on agent"),
		GenCommand(requests.DescItf(), "agent-itf", "Show interfaces on agent"),
		GenCommand(requests.DescSi(), "agent-si", "Show service instances on agent"),
		GenCommand(requests.DescVrf(), "agent-vrf", "Show vrfs on agent "),
		GenCommand(requests.DescPeering(), "agent-peering", "Peering with controller on agent"),
		GenCommand(requests.DescVn(), "agent-vn", "Show virtual networks on agent"),
		GenCommand(requests.DescMpls(), "agent-mpls", "Show mpls on agent"),
		requests.Follow(),
		requests.Path(),
		GenCommand(requests.DescRiSummary(), "controller-ri", "Show routing instances on controller"),
		GenCommand(requests.DescCtrlRoute(), "controller-route", "Show routes on controller"),
		GenCommand(requests.DescCtrlRouteSummary(), "controller-route-summary", "Show routes summary on controller"),
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
