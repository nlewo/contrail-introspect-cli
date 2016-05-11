Installation
============

    $ export GOPATH=$PWD
    $ go get github.com/codegangsta/cli
    $ go get github.com/moovweb/gokogiri
    $ go get github.com/gosuri/uitable
    $ go install github.com/nlewo/contrail-introspect-cli
    $ $GOPATH/bin/contrail-introspect-cli help

Usage
=====

To list interfaces:

    $ $GOPATH/bin/contrail-introspect-cli itf vrouter_fqdn

To get nexthops of a routes in the vrf vrf_public:

    $ $GOPATH/bin/contrail-introspect-cli route --long vrouter_fqdn vrf_public 
