package descriptions

import "fmt"
import "log"

import cli "gopkg.in/urfave/cli.v2"

import "github.com/nlewo/contrail-introspect-cli/utils"

func Follow() *cli.Command {
	return &cli.Command{
		Name:      "agent-follow",
		Usage:     "From a compute node, a vrf and a route, follow the route to destination",
		ArgsUsage: "vrouter-fqdn vrf-name route-prefix",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "fqdn",
				Usage: "Suffix appended to introspect hostnames",
			}},
		Action: func(c *cli.Context) error {
			// Get the interface
			suffix := ""
			if c.String("fqdn") != "" {
				suffix = "." + c.String("fqdn")
			}
			ip := c.Args().Get(2)

			fmt.Printf("1. Starting on %s for the route %s in the vrf %s\n", c.Args().Get(0), ip, c.Args().Get(1))
			page := Route().PageBuilder(c.Args().Slice())
			col, e := page.Load(Route())
			if e != nil {
				log.Fatal(e)
			}
			elt := col.SearchStrict(ip)
			label, err := elt[0].GetField("path_list/list/PathSandeshData/label")
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
			page = Mpls().PageBuilder(args)
			col, e = page.Load(Mpls())
			if e != nil {
				log.Fatal(e)
			}

			elt = col.SearchStrict(label)
			itf, err := elt[0].GetField("nh/NhSandeshData/itf")
			if err != nil {
				log.Fatal(err)
			}

			// elt.Long()

			page = Interface().PageBuilder(args)
			col, e = page.Load(Interface())
			if e != nil {
				log.Fatal(e)
			}
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
