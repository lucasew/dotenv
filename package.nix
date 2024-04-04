{ buildGoModule
, lib
}:

buildGoModule {
  name = "dotenv";
  version = "0.0.1";

  vendorHash = "sha256-FadCW9CogyX06kWH0MFVr4Z4BnidUp/ccH2Po11zIgE=";

  src = ./.;

  meta = with lib; {
    description = "Dotenv as a binary that loads the dotenv and calls the program";
    homepage = "https://github.com/lucasew/dotenv";
    platforms = platforms.linux;
  };
}
