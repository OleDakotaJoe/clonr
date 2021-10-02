package cmd


import (
	"clonr/utils"
	"fmt"
	"github.com/go-git/go-git"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

var outputDir string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones a git project and initializes template engine",
	Long:  `All software has versions. This is Clonr'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Clonr Project Templating CLI --- v0.0.1")
		source := args[0]
		validateUrl(source)
		cloneProject(source, outputDir)

	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVar(&outputDir, "out", "./clonr-app", "The git URL to read from")
}


func validateUrl(source string) {
	fmt.Println("Validating src URL")
	_,err := url.ParseRequestURI(source)
	utils.CheckForError(err)
	fmt.Println("The src URL you provided is valid.")
}

func cloneProject(source string, outputDir string) {
	fmt.Println("Cloning git repo... Please Wait")
	_, err := git.PlainClone(outputDir, false,&git.CloneOptions{
		URL: source,
		Progress: os.Stdout,
	} )
	utils.CheckForError(err)
}