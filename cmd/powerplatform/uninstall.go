package main

import (
	"github.com/getporter/powerplatform/pkg/powerplatform"
	"github.com/spf13/cobra"
)

func buildUninstallCommand(m *powerplatform.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute(cmd.Context())
		},
	}
	return cmd
}
