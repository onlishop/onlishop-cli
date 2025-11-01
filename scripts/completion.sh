#!/usr/bin/env bash

rm -rf completions
mkdir completions
go run . completion bash > completions/onlishop-cli.bash
go run . completion zsh > completions/onlishop-cli.zsh
go run . completion fish > completions/onlishop-cli.fish