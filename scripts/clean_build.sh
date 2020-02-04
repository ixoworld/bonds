#!/usr/bin/env bash

rm -rf "$HOME"/.bondsd
rm -rf "$HOME"/.bondscli

cd "$HOME"/go/src/github.com/ixoworld/bonds/ || exit
make install
