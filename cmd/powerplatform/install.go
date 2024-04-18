package main

import (
	"github.com/getporter/powerplatform/pkg/powerplatform"
	"github.com/spf13/cobra"
)

var (
	commandFile string
)

func buildInstallCommand(m *powerplatform.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Execute the install functionality of this mixin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//Do something here if needed
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute(cmd.Context())
		},
	}
	return cmd
}
