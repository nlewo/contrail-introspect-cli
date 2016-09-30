package requests

import "fmt"
import "log"

import "github.com/codegangsta/cli"

import "github.com/nlewo/contrail-introspect-cli/utils"

func Follow() cli.Command {
	return cli.Command{
		Name:      "agent-follow",
		Usage:     "From a compute node, a vrf and a route, follow the route to destination",
		ArgsUsage: "vrouter-fqdn vrf-name route-prefix",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "fqdn",
				Usage: "Suffix appended to introspect hostnames",
			}},
		Action: func(c *cli.Context) error {
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
			label, err := elt[0].GetField("path_list/list/PathSandeshData/label");
			if err != nil {
				log.Fatal(err)
			}
			nh, err := elt[0].GetField("path_list/list/PathSandeshData/nh/NhSandeshData/dip")
			if err != nil {
				log.Fatal(err)
			}

			nh_fqdn := utils.ResolveIp(nh) + suffix
			// elt.Long()
			fmt.Printf("2. Go with MPLS label %s to %s\n", label, nh_fqdn)

			args := make([]string, 1)
			args[0] = nh_fqdn
			page = DescMpls().PageBuilder(args)
			col = page.Load(DescMpls())
			elt = col.SearchStrict(label)
			itf, err := elt[0].GetField("nh/NhSandeshData/itf")
			if err != nil {
				log.Fatal(err)
			}

			// elt.Long()

			page = DescItf().PageBuilder(args)
			col = page.Load(DescItf())
			elt = col.SearchStrict(itf)
			vm_uuid, err := elt[0].GetField("vm_uuid")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("3. To interface %s of vm %s\n", itf, vm_uuid)
			// elt.Long()

			return nil
		},
	}
}
