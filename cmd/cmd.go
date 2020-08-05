package cmd

import (
	"log"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/spf13/cobra"
	"github.com/sweetfdtd/pkg/fdtd"
)

var (
	// Used for flags.
	workPath string

	rootCmd = &cobra.Command{
		Use:   "sweetfdtd",
		Short: "sweetfdtd is an custom FDTD runner.",
		Run: func(cmd *cobra.Command, args []string) {
			defer termui.Close()
			ts := &fdtd.FDTDTaskSet{}
			ts.BuildFromPath(workPath)
			if err := termui.Init(); err != nil {
				log.Fatalf("failed to initialize termui: %v", err)
			}
			tui := fdtd.NewTuiView(ts)
			termui.Render(tui.Grid)
			go func() {
				ticker := time.NewTicker(time.Millisecond * 100)
				for range ticker.C {
					tui.Update()
					termui.Render(tui.Grid)
				}
			}()
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
