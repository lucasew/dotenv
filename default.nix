let
  pkgs = import <nixpkgs> {};
in pkgs.buildGoModule rec {
  name = "dotenv";
  version = "0.0.1";
  vendorSha256 = "0092fdfs73vxf3f9yllxg037i1mgap0x11s5xbs2b0x8s1dl59qm";
  src = ./.;
  meta = with pkgs.lib; {
    description = "Dotenv as a binary that loads the dotenv and calls the program";
    homepage = "https://github.com/lucasew/dotenv";
    platforms = platforms.linux;
  };
}
