package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sweetfdtd/pkg/fdtd"
)

var (
	// Used for flags.
	workPath string

	rootCmd = &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			ts := &fdtd.FDTDTaskSet{}
			ts.BuildFromPath(workPath)
			ts.Run()
		},
	}
)

func init() {

	rootCmd.Flags().StringVarP(&workPath, "path", "d", "", "")
	rootCmd.MarkFlagRequired("path")

}

func Exec() {
	rootCmd.Execute()
}
