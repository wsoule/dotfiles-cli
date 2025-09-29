# Documentation: https://docs.brew.sh/Formula-Cookbook
#                https://rubydoc.brew.sh/Formula
class Dotfiles < Formula
  desc "Modern dotfiles manager with Homebrew and GNU Stow integration"
  homepage "https://github.com/wyatsoule/go-dotfiles"
  url "https://github.com/wyatsoule/go-dotfiles/archive/v1.0.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"
  version "1.0.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", "#{bin}/dotfiles"
  end

  def post_install
    puts ""
    puts "🎉 Dotfiles Manager installed successfully!"
    puts ""
    puts "Get started with:"
    puts "  dotfiles onboard          # Complete developer setup"
    puts "  dotfiles init             # Initialize configuration"
    puts "  dotfiles github setup     # Set up GitHub SSH"
    puts ""
    puts "For help: dotfiles --help"
  end

  test do
    system "#{bin}/dotfiles", "--version"
    system "#{bin}/dotfiles", "--help"
  end
end