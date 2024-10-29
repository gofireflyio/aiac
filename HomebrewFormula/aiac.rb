# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Aiac < Formula
  desc "Artificial Intelligence Infrastructure-as-Code Generator"
  homepage "https://github.com/gofireflyio/aiac"
  version "5.3.0"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/gofireflyio/aiac/releases/download/v5.3.0/aiac_5.3.0_darwin-amd64.tar.gz"
      sha256 "30c3c0026fb9a45c7579fad9fb92e978a2f81e1eb7c7b36aea50dd2aa7c2aae8"

      def install
        bin.install "aiac"
      end
    end
    on_arm do
      url "https://github.com/gofireflyio/aiac/releases/download/v5.3.0/aiac_5.3.0_darwin-arm64.tar.gz"
      sha256 "d7de1ee6cabdeae40334a97536cefdf9db424ce1edc35ec19b6ce8361fd0da35"

      def install
        bin.install "aiac"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/gofireflyio/aiac/releases/download/v5.3.0/aiac_5.3.0_linux-amd64.tar.gz"
        sha256 "c11af7053dcbf946375670e022612ab8cec03973e3deb828cfe57d090a8ba606"

        def install
          bin.install "aiac"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/gofireflyio/aiac/releases/download/v5.3.0/aiac_5.3.0_linux-arm64.tar.gz"
        sha256 "5067ccefcffe9726bed23a8f68dc4bed67fc19d767a2b67b1d10c0babb8002df"

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
