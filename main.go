package main

import (
	"clonr/cmd"
	"clonr/config"
)

func main() {
	cmd.Execute()
}

func init() {
	config.ConfigureLogger()
}