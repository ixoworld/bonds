#!/usr/bin/env bash

rm -rf "$HOME"/.bondsd
rm -rf "$HOME"/.bondscli

cd ../ # assumes currently in bonds/scripts/ folder
make install
