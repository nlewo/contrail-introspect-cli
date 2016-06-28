CLI on ContraiL Introspects
===========================

## Installation

    $ go get github.com/nlewo/contrail-introspect-cli/contrail-introspect-cli

## Usage Examples

- List interfaces
```
    $ contrail-introspect-cli agent-itf vrouter_fqdn -l
    00000000-0000-0000-0000-000000000000 bond0.1002 default-domain:default-project:ip-fabric:__default__ 
    039b3555-e83d-480c-89d2-fb2cf767bf55 tap039b3555-e8 default-domain:default-project:network:network
    08893790-a8e2-4283-800000000-0000-0000-0000-000000000000 vhost0 default-domain:default-project:ip-fabric:__default__ 
```

- Get nexthops for `192.168.1.5` in the vrf `net1`
```
    $ contrail-introspect-cli --hosts hosts agent-route vrouter-fqdn domain:project:net1:net1 -s 192.168.1.5 -l
    Src 192.168.1.5
        Dst                        	Peers        	MPLS label	Interface	Dest VN                        
        vrouter-1                       10.12.128.10	30        	         	domain:project:net1
        vrtouer-1                       10.12.128.11	30        	         	domain:project:net1
```
The `--hosts` option takes a `hosts` file to translate introspect IPs to DNS names.

- Follow a route to the destination interface
```
    $ contrail-introspect-cli --hosts hosts follow vrouter-1.example.com vrf-name 10.210.3.5  
    1. Starting on vrouter-1.example.com for the route 10.210.3.5 in the vrf vrf-name
    2. Go with MPLS label 129 to vrouter-2.example.com
    3. To interface tap2a452941-0b of vm d1bd1a84-b479-4897-a6c4-4dce7c4c8f4d
```

- Get route details from a controller
```
	$ contrail-introspect-cli controller-route controller-1.example.com default-domain:openstack:public:public  -s 145 -l
	Prefix 8.8.8.145/32
	        Protocol	Nexthop      	Peers        	MPLS label
	        XMPP    	d-ocnclc-002w	d-ocnclc-002w	18        
	        BGP     	d-ocnclc-002w	d-octclc-0001	18        
	        XMPP    	d-ocnclc-000r	d-ocnclc-000r	21        

```