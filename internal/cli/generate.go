package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"echostrike/internal/generator"
	"echostrike/internal/sender"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var (
	rate         int
	duration     time.Duration
	templateName string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate high-volume syslog traffic from templates",
	Long: `Generate realistic syslog traffic using built-in templates.
Support rate limiting and duration control.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate protocol
		var proto sender.Protocol
		switch strings.ToLower(protocol) {
		case "tcp":
			proto = sender.TCP
		case "udp":
			proto = sender.UDP
		case "tls":
			proto = sender.TLS
		default:
			fmt.Printf("Invalid protocol: %s\n", protocol)
			os.Exit(1)
		}

		// Initialize Generator
		gen := generator.New()

		// Initialize Sender
		s, err := sender.NewSender(proto, host, port)
		if err != nil {
			fmt.Printf("Error connecting to target: %v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		fmt.Printf("Starting generation: Template=%s Target=%s:%d Rate=%d/s Duration=%s\n",
			templateName, host, port, rate, duration)

		// Rate Limiter
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()

		timeout := time.After(duration)
		count := 0

		for {
			select {
			case <-timeout:
				fmt.Printf("\nCompleted. Sent %d messages.\n", count)
				return
			case <-ticker.C:
				// Generate log
				logMsg, err := gen.Generate(templateName)
				if err != nil {
					fmt.Printf("Error generating log: %v\n", err)
					os.Exit(1)
				}

				// Create syslog message
				msg := syslog.NewMessage(logMsg)
				msg.AppName = tag

				// Set Format
				if strings.ToLower(format) == "rfc5424" {
					msg.Format = syslog.RFC5424
				}

				// Set Facility/Severity (Reuse global flags)
				f, _ := syslog.ParseFacility(strings.ToLower(facility))
				msg.Facility = f
				sev, _ := syslog.ParseSeverity(strings.ToLower(severity))
				msg.Severity = sev

				// Send
				if err := s.Send(msg.String() + "\n"); err != nil {
					fmt.Printf("Error sending: %v\n", err)
				}
				count++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Target IP/Hostname")
	generateCmd.Flags().IntVar(&port, "port", 514, "Target Port")
	generateCmd.Flags().StringVar(&protocol, "proto", "udp", "Protocol (tcp, udp, tls)")
	generateCmd.Flags().StringVarP(&templateName, "template", "T", "ssh-failed", "Log template name")
	generateCmd.Flags().IntVarP(&rate, "rate", "r", 1, "Logs per second")
	generateCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Duration to run")

	generateCmd.Flags().StringVarP(&tag, "tag", "t", "echostrike", "Syslog tag/app-name")
	generateCmd.Flags().StringVar(&format, "format", "rfc3164", "Syslog format")
	generateCmd.Flags().StringVar(&facility, "facility", "local0", "Facility")
	generateCmd.Flags().StringVar(&severity, "severity", "info", "Severity")
}
