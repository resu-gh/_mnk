#!/usr/bin/env sh

find . -iname '*.go' -or -iname 'Makefile' | entr -r make
