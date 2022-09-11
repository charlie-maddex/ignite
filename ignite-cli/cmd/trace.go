/*
Copyright © 2022 Charlie Maddex (charlie@multi.sh)
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

// traceCmd represents the trace command
var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Display the route packets take to reach a host",
	Long:  `Display the route packets take to reach a host. This is useful for debugging network issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		white := color.New(color.FgWhite).SprintFunc()

		if len(args) == 0 {
			fmt.Fprintf(color.Output, "%s\n", red("Please provide a host or IP address to trace."))
			os.Exit(1)
		}

		pinger, err := ping.NewPinger(args[0])
		if err != nil {
			fmt.Fprintf(color.Output, "%s\n", red(err))
			os.Exit(1)
		}

		fmt.Fprintf(color.Output, "TRACE %s (%s):\n", green(pinger.Addr()), green(pinger.IPAddr()))

		errch := make(chan error, 1)
		if runtime.GOOS == "windows" {
			tracecmd := exec.Command("traceroute", pinger.Addr())
			stdout, err := tracecmd.StdoutPipe()
			if err != nil {
				fmt.Fprintf(color.Output, "%s\n", red(err))
			}
			if err := tracecmd.Start(); err != nil {
				fmt.Fprintf(color.Output, "%s\n", red(err))
			}
			go func() {
				errch <- tracecmd.Wait()
			}()

			go func() {
				scanner := bufio.NewScanner(stdout)
				fmt.Println("")
				for scanner.Scan() {
					line := scanner.Text()
					fmt.Fprintf(color.Output, "%s\n", white(line))
				}
			}()

			select {
			case <-time.After(time.Second * 12):
				fmt.Fprintf(color.Output, "%s\n", red("traceroute timed out."))
				return
			case err := <-errch:
				if err != nil {
					fmt.Fprintf(color.Output, "%s\n", red(err))
				}
			}
		} else {
			tracecmd := exec.Command("traceroute", pinger.Addr())
			stdout, err := tracecmd.StdoutPipe()
			if err != nil {
				fmt.Fprintf(color.Output, "%s\n", red(err))
			}
			if err := tracecmd.Start(); err != nil {
				fmt.Fprintf(color.Output, "%s\n", red(err))
			}
			go func() {
				errch <- tracecmd.Wait()
			}()

			go func() {
				scanner := bufio.NewScanner(stdout)
				fmt.Println("")
				for scanner.Scan() {
					line := scanner.Text()
					fmt.Fprintf(color.Output, "%s\n", white(line))
				}
			}()

			select {
			case <-time.After(time.Second * 12):
				fmt.Fprintf(color.Output, "%s\n", red("traceroute timed out."))
				return
			case err := <-errch:
				if err != nil {
					fmt.Fprintf(color.Output, "%s\n", red(err))
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(traceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
