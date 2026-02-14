package cli

import (
	"fmt"
	"os"
	"time"

	"echostrike/internal/generator"
	"echostrike/internal/sender"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var simulationType string

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Run attack simulations (brute-force, port-scan)",
	Long:  `Simulate specific attack patterns to test detection rules.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Sender
		var proto sender.Protocol = sender.UDP
		if protocol == "tcp" {
			proto = sender.TCP
		}
		if protocol == "tls" {
			proto = sender.TLS
		}

		s, err := sender.NewSender(proto, host, port)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		targetIP := generator.RandomIP()   // Victim IP
		attackerIP := generator.RandomIP() // Attacker IP

		fmt.Printf("Starting simulation: %s (Attacker: %s -> Victim: %s)\n", simulationType, attackerIP, targetIP)

		switch simulationType {
		case "brute-force":
			// SSH Brute Force: Many failures, then one success
			users := []string{"root", "admin", "service", "dbadmin"}
			for _, user := range users {
				for i := 0; i < 3; i++ {
					msg := fmt.Sprintf("Failed password for %s from %s port %d ssh2", user, attackerIP, generator.RandomPort())
					sendSyslog(s, msg)
					time.Sleep(100 * time.Millisecond)
				}
			}
			// Finally success
			msg := fmt.Sprintf("Accepted publickey for root from %s port %d ssh2", attackerIP, generator.RandomPort())
			sendSyslog(s, msg)

		case "port-scan":
			// Rapid connection attempts to different ports
			ports := []int{21, 22, 23, 25, 53, 80, 443, 3306, 8080}
			for _, p := range ports {
				msg := fmt.Sprintf("SRC=%s DST=%s PROTO=TCP DPT=%d ACTION=DROP", attackerIP, targetIP, p)
				sendSyslog(s, msg)
				time.Sleep(50 * time.Millisecond)
			}

		default:
			fmt.Println("Unknown simulation type. Available: brute-force, port-scan")
			os.Exit(1)
		}

		fmt.Println("Simulation completed.")
	},
}

func sendSyslog(s *sender.Sender, content string) {
	msg := syslog.NewMessage(content)
	msg.Facility = syslog.LOG_AUTH
	msg.Severity = syslog.LOG_WARNING
	msg.AppName = "sshd" // Simulate SSH daemon

	fmt.Printf("Sending: %s\n", content)
	s.Send(msg.String() + "\n")
}

func init() {
	rootCmd.AddCommand(simulateCmd)

	simulateCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Target IP")
	simulateCmd.Flags().IntVar(&port, "port", 514, "Target Port")
	simulateCmd.Flags().StringVar(&protocol, "proto", "udp", "Protocol")
	simulateCmd.Flags().StringVarP(&simulationType, "type", "T", "brute-force", "Simulation type (brute-force, port-scan)")
}
