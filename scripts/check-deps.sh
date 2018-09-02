#!/usr/bin/env bash

if ! [ -x "$(command -v jq)" ]; then
	echo "Please install jq: https://stedolan.github.io/jq/"
	exit 1
fi
