{
  nixpkgs ? builtins.fetchTarball {
    # From Release 19.03
    url = https://github.com/nixos/nixpkgs/archive/c8db7a8a16ee9d54103cade6e766509e1d1c8d7b.tar.gz;
    sha256 = "1b3h4mwpi10blzpvgsc0191k4shaw3nw0qd2p82hygbr8vv4g9dv";
  }
}:

with import nixpkgs {};

buildGoPackage {
  name = "contrail-introspect-cli-unstable";

  goPackagePath = "github.com/nlewo/contrail-introspect-cli";

  buildInputs = [ pkgconfig libxml2 ];

  src = ./.;
  goDeps = ./deps.nix;
}
