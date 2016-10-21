package descriptions

import "fmt"
import "log"
import "github.com/codegangsta/cli"

import "github.com/nlewo/contrail-introspect-cli/collection"

func Ping() cli.Command {
	return cli.Command{
		Name:      "ping",
		Usage:     "Generate one ping packet from a port id to dest ip and port number",
		ArgsUsage: "vrouter-fqdn port-uuid dest-ip dest-port",
		Action: func(c *cli.Context) error {
			agent := c.Args()[0]
			portUuid := c.Args()[1]
			destIp := c.Args()[2]
			destPort := c.Args()[3]
			col, e := collection.LoadCollection(Interface(), []string{agent})
			if e != nil {
				log.Fatal(fmt.Sprintf("Interface can not be loaded: %s", e))
			}

			elt, e := col.SearchStrictUniqueByKey("uuid", portUuid)
			if e != nil {
				fmt.Printf("%s\n", e)
				return e
			}
			vrfName, _ := elt.GetField("vrf_name")
			ipAddr, _ := elt.GetField("ip_addr")
			col, _ = collection.LoadCollection(AgentPing(), []string{agent, vrfName, ipAddr, "", destIp, destPort})
			col.Short()
			return nil
		},
	}
}
