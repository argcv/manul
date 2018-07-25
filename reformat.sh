#!/usr/bin/env bash

find . -name '*.go' | xargs -n 1 -I{} -P 6 sh -c 'echo "reformat: {}" && gofmt -w {}'
