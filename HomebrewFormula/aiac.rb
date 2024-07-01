# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Aiac < Formula
  desc "Artificial Intelligence Infrastructure-as-Code Generator"
  homepage "https://github.com/gofireflyio/aiac"
  version "5.1.1"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/gofireflyio/aiac/releases/download/v5.1.1/aiac_5.1.1_darwin-amd64.tar.gz"
      sha256 "c185225f5ab8334e8ada1c4b85616e015caf77f7ea7a8d921b56c1a4f51dd87a"

      def install
        bin.install "aiac"
      end
    end
    on_arm do
      url "https://github.com/gofireflyio/aiac/releases/download/v5.1.1/aiac_5.1.1_darwin-arm64.tar.gz"
      sha256 "4262f7ed048a11dbef9539a9a3584df561c880b9b6796ad287c84668bb7c5eea"

      def install
        bin.install "aiac"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/gofireflyio/aiac/releases/download/v5.1.1/aiac_5.1.1_linux-amd64.tar.gz"
        sha256 "b6b3f02781580a89ef1e80c1ec5ae66a91e0c158d34e6e27e41089ce66fb9aac"

        def install
          bin.install "aiac"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/gofireflyio/aiac/releases/download/v5.1.1/aiac_5.1.1_linux-arm64.tar.gz"
        sha256 "6372ae0b4617b99ad70fe3422d3299abccd003b0b8281544d309e25c61c16eb1"

        def install
          bin.install "aiac"
        end
      end
    end
  end

  test do
    system "#{bin}/aiac", "--help"
  end
end