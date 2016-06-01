package main

import "fmt"
import "github.com/codegangsta/cli"

func Follow() cli.Command {
	return cli.Command{
		Name:      "follow",
		Usage:     "From a compute node, a vrf and a route, follow the route to destination",
		ArgsUsage: "vrouter-fqdn vrf-name route-prefix",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "fqdn",
				Usage:"Suffix appended to introspect hostnames",
			}},
		Action: func(c *cli.Context) {
			// Get the interface
			suffix := ""
			if c.String("fqdn") != "" {
				suffix = "." + c.String("fqdn")
			}
			ip := c.Args()[2]

			fmt.Printf("1. Starting on %s for the route %s in the vrf %s\n", c.Args()[0], ip, c.Args()[1])
			page := DescRoute().PageBuilder(c.Args())
			col := page.Load(DescRoute())
			elt := col.SearchStrict(ip)
			label := elt[0].GetField("path_list/list/PathSandeshData/label")
			nh := elt[0].GetField("path_list/list/PathSandeshData/nh/NhSandeshData/dip")
			nh_fqdn := ResolveIp(nh) + suffix
			// elt.Long()
			fmt.Printf("2. Go with MPLS label %s to %s\n", label, nh_fqdn)

			args := make([]string, 1)
			args[0] = nh_fqdn
			page = DescMpls().PageBuilder(args)
			col = page.Load(DescMpls())
			elt = col.SearchStrict(label)
			itf := elt[0].GetField("nh/NhSandeshData/itf")
			// elt.Long()

			page = DescItf().PageBuilder(args)
			col = page.Load(DescItf())
			elt = col.SearchStrict(itf)
			fmt.Printf("3. To interface %s of vm %s\n", itf, elt[0].GetField("vm_uuid"))
			// elt.Long()
		},
	}
}
