package main

import (
	"clonr/cmd"
	"clonr/config"
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	cmd.Execute()
}

func init() {
	config.ConfigureLogger()
	title := figure.NewFigure("clonr", "rounded", true).String()
	fmt.Print(title)
}