package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/oledakotajoe/clonr/cmd"
	"github.com/oledakotajoe/clonr/config"
)

func main() {
	cmd.Execute()
}

func init() {
	config.ConfigureLogger()
	title := figure.NewFigure("clonr", "rounded", true).String()
	fmt.Print(title)
	fmt.Println("Clonr CLI -- Project Templating Engine")
}
