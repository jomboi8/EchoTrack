package cli

import (
	"fmt"
	"os"

	"echostrike/internal/generator"
	"echostrike/pkg/syslog"

	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview log templates without sending",
	Long:  `Dry-run mode to see what logs would look like.`,
	Run: func(cmd *cobra.Command, args []string) {
		gen := generator.New()

		fmt.Printf("Previewing template: %s\n", templateName)
		fmt.Println("----------------------------------------")

		for i := 0; i < 5; i++ {
			logMsg, err := gen.Generate(templateName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			msg := syslog.NewMessage(logMsg)
			msg.AppName = tag

			if format == "rfc5424" {
				msg.Format = syslog.RFC5424
			}

			fmt.Println(msg.String())
		}
		fmt.Println("----------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
	previewCmd.Flags().StringVarP(&templateName, "template", "T", "ssh-failed", "Log template name")
	previewCmd.Flags().StringVarP(&tag, "tag", "t", "echostrike", "Syslog tag")
	previewCmd.Flags().StringVar(&format, "format", "rfc3164", "Syslog format")
}
