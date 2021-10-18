package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/oledakotajoe/clonr/cmd"
	"github.com/oledakotajoe/clonr/config"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd.Execute()
}

func init() {
	config.ConfigureLogger()
	title := figure.NewFigure("clonr", "rounded", true).String()
	fmt.Print(title)
	fmt.Println("Clonr CLI -- Project Templating Engine")
	err := doc.GenMarkdownTree(cmd.RootCmd, "./.resources")
	if err != nil {
		fmt.Print(err)
	}
}
