#!/usr/bin/env bash

rm -rf "$HOME"/.bondsd
rm -rf "$HOME"/.bondscli

make install # assumes currently in project directory
