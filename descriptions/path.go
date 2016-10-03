package descriptions

import "fmt"
import "github.com/codegangsta/cli"
import "log"

import "github.com/nlewo/contrail-introspect-cli/utils"

func Path() cli.Command {
	return cli.Command{
		Name:      "controller-path",
		Usage:     "Find the path to go from a prefix to another on in a routing instance",
		ArgsUsage: "controller-fqdn ri source-prefix dest-prefix",
		Action: func(c *cli.Context) error {

			controller := c.Args()[0]
			ri := c.Args()[1]
			srcIp := c.Args()[2]
			dstIp := c.Args()[3]

			page := CtrlRoute().PageBuilder([]string{controller, ri})
			col, e := page.Load(CtrlRoute())
			if e != nil {
				log.Fatal(e)
			}

			elt := col.SearchStrict(srcIp)
			if len(elt) < 1 {
				log.Fatal(fmt.Sprintf("Prefix %s not found in RI %s", srcIp, ri))
			}
			srcNode, err := elt[0].GetField("paths/list/ShowRoutePath/next_hop")
			if err != nil {
				log.Fatal(err)
			}

			page = CtrlRoute().PageBuilder([]string{controller, ri})
			col, e = page.Load(CtrlRoute())
			if e != nil {
				log.Fatal(e)
			}
			elt = col.SearchStrict(dstIp)
			if len(elt) < 1 {
				log.Fatal(fmt.Sprintf("Prefix %s not found in RI %s", dstIp, ri))
			}
			dstNode, err := elt[0].GetField("paths/list/ShowRoutePath/next_hop")
			if err != nil {
				log.Fatal(err)
			}

			label, err := elt[0].GetField("paths/list/ShowRoutePath/label")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("From prefix %s on %s to dst %s on %s with label %s\n", srcIp, utils.ResolveIp(srcNode), dstIp, utils.ResolveIp(dstNode), label)

			return nil
		},
	}
}
