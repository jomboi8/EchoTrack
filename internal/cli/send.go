package cli

import (
	"fmt"
	"os"
	"strings"

	"echostrike/internal/sender"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var (
	host     string
	port     int
	protocol string
	message  string
	tag      string
	format   string
	facility string
	severity string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a single syslog message",
	Long:  `Send a single syslog message to a remote server via TCP, UDP, or TLS.`,
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
			fmt.Printf("Invalid protocol: %s. Use tcp, udp, or tls\n", protocol)
			os.Exit(1)
		}

		// Create message
		msg := syslog.NewMessage(message)
		msg.AppName = tag

		// Set format
		if strings.ToLower(format) == "rfc5424" {
			msg.Format = syslog.RFC5424
		} else {
			msg.Format = syslog.RFC3164
		}

		// Set Facility/Severity
		f, err := syslog.ParseFacility(strings.ToLower(facility))
		if err != nil {
			fmt.Printf("Error parsing facility: %v\n", err)
			os.Exit(1)
		}
		msg.Facility = f

		sev, err := syslog.ParseSeverity(strings.ToLower(severity))
		if err != nil {
			fmt.Printf("Error parsing severity: %v\n", err)
			os.Exit(1)
		}
		msg.Severity = sev

		// Initialize sender
		s, err := sender.NewSender(proto, host, port)
		if err != nil {
			fmt.Printf("Error creating sender: %v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		// Send
		logLine := msg.String()
		fmt.Printf("Sending to %s:%d (%s): %s\n", host, port, proto, logLine)

		if err := s.Send(logLine + "\n"); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Message sent successfully.")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVar(&host, "host", "127.0.0.1", "Target IP/Hostname")
	sendCmd.Flags().IntVar(&port, "port", 514, "Target Port")
	sendCmd.Flags().StringVar(&protocol, "proto", "udp", "Protocol (tcp, udp, tls)")
	sendCmd.Flags().StringVarP(&message, "message", "m", "EchoStrike Test Message", "Log message content")
	sendCmd.Flags().StringVarP(&tag, "tag", "t", "echostrike", "Syslog tag/app-name")
	sendCmd.Flags().StringVar(&format, "format", "rfc3164", "Syslog format (rfc3164, rfc5424)")
	sendCmd.Flags().StringVar(&facility, "facility", "local0", "Syslog facility (e.g. auth, cron, local0)")
	sendCmd.Flags().StringVar(&severity, "severity", "info", "Syslog severity (e.g. info, debug, err)")
}
