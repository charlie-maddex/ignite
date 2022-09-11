/*
Copyright © 2022 Charlie Maddex (charlie@multi.sh)
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ICMP ping a host and return the results",
	Long:  `Provide a host or IP address to ping.`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		if len(args) == 0 {
			fmt.Fprintf(color.Output, "%s\n", red("Please provide a host or IP address to ping."))
			os.Exit(1)
		}

		pinger, err := ping.NewPinger(args[0])
		if err != nil {
			fmt.Fprintf(color.Output, "%s\n", red(err))
			os.Exit(1)
		}

		// Listen for Ctrl-C.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				pinger.Stop()
			}
		}()

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Fprintf(color.Output, "%s bytes from %s: time=%v\n", green(pkt.Nbytes), green(pkt.IPAddr), green(pkt.Rtt))
		}

		pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Fprintf(color.Output, "\n--- %s ping statistics ---\n", green(stats.Addr))
			fmt.Fprintf(color.Output, "%s packets transmitted, %s packets received, %s duplicates, %s%% packet loss\n", green(stats.PacketsSent), green(stats.PacketsRecv), green(stats.PacketsRecvDuplicates), green(stats.PacketLoss))
		}
		fmt.Fprintf(color.Output, "PING %s (%s):\n", green(pinger.Addr()), green(pinger.IPAddr()))

		err = pinger.Run()
		if err != nil {
			fmt.Fprintf(color.Output, "%s\n", red(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
