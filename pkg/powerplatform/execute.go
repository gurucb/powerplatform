package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	// "get.porter.sh/porter/pkg/context"
	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

func (m *Mixin) loadAction(ctx context.Context) (*Action, error) {
	var action Action
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute(ctx context.Context) error {
	action, err := m.loadAction(ctx)
	if err != nil {
		return err
	}
	fmt.Print("Action Name: ")
	fmt.Println(action.Name)

	licenses := action.Steps[0].Licenses
	licenseString, err := json.Marshal(licenses)
	formattedlicenseString := strings.ReplaceAll(string(licenseString), "\"", "\\\"")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("License: ")
	fmt.Println(string(formattedlicenseString))

	fmt.Println("Dependencies: ")
	dependencies := action.Steps[0].Dependencies
	dependencyString, err := json.Marshal(dependencies)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(dependencyString))
	formatteddependencyString := strings.ReplaceAll(string(dependencyString), "\"", "\\\"")

	// fmt.Println(action.Steps[0].Flags.ToSlice(builder.Dashes(DefaultFlagDashes)))
	fmt.Println("Supported Regions: ")
	fmt.Println(action.Steps[0].SupportedRegions)
	supportedRegions := "[" + strings.Join(action.Steps[0].SupportedRegions, " ") + "]"

	fmt.Println("Target Environment: ")
	fmt.Println(action.Steps[0].TargetEnvironment)
	targetEnvironment := action.Steps[0].TargetEnvironment

	fmt.Println("PackageId: ")
	fmt.Println(action.Steps[0].PackageId)
	packageId := action.Steps[0].PackageId
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Licenses", formattedlicenseString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Dependencies", formatteddependencyString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("SupportedRegions", supportedRegions))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("TargetEnvironment", targetEnvironment))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("PackageId", packageId))
	_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	return err
}
