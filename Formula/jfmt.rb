# Homebrew formula for jfmt
# To use: brew tap josharsh/tap && brew install jfmt
# Or: brew install josharsh/tap/jfmt

class Jfmt < Formula
  desc "Fast JSON formatter with colors, clipboard support, and auto-fix"
  homepage "https://github.com/josharsh/jfmt"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/josharsh/jfmt/releases/download/v#{version}/jfmt-darwin-arm64"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
    else
      url "https://github.com/josharsh/jfmt/releases/download/v#{version}/jfmt-darwin-amd64"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/josharsh/jfmt/releases/download/v#{version}/jfmt-linux-arm64"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/josharsh/jfmt/releases/download/v#{version}/jfmt-linux-amd64"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "jfmt"
  end

  test do
    assert_equal "{\n  \"a\": 1\n}", pipe_output("#{bin}/jfmt", '{"a":1}').strip
  end
end
