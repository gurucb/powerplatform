package main

import (
	"github.com/getporter/powerplatform/pkg/powerplatform"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *powerplatform.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute(cmd.Context())
		},
	}
	return cmd
}
