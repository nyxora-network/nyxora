# typed: false
# frozen_string_literal: true

# Documentation: https://docs.brew.sh/Formula-Cookbook
#                https://rubydoc.brew.sh/Formula

class Nyxora < Formula
  desc "Adaptive multi-transport VPN/tunnel orchestrator"
  homepage "https://github.com/nyxorammd-lgtm/nyxora"
  url "https://github.com/nyxorammd-lgtm/nyxora/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "d455530328c00b6bb0dfc9dc97662c2ebc85e67e05e7ca2bfadc12f3a339a09a"
  license "MIT"
  head "https://github.com/nyxorammd-lgtm/nyxora.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X main.version=#{version}
    ]
    system "go", "build", *std_go_args(ldflags:), "./cmd/nyxora"
  end

  test do
    assert_match "nyxora v#{version}", shell_output("#{bin}/nyxora version")
  end
end
