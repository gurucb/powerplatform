package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/google/uuid"
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
	//TODO: Log Action Name

	licenses := action.Steps[0].Licenses
	licenseString, err := json.Marshal(licenses)
	if err != nil {
		fmt.Println("Failure parsing license JSON")
		fmt.Println("Error:", err)
	}
	formattedlicenseString := string(licenseString)
	fmt.Println("Unformatted String: \n" + formattedlicenseString)
	formattedlicenseString = "\\" + strings.ReplaceAll(formattedlicenseString, "\"", "\\") + "\\"
	fmt.Println("License String:\n" + formattedlicenseString)
	//TODO: Log formattedlicenseString

	//WORKING below here
	dependencies := action.Steps[0].Dependencies
	dependencyString, err := json.Marshal(dependencies)
	if err != nil {
		fmt.Println("Failure parsing Dependency JSON")
		fmt.Println("Error:", err)
	}
	formatteddependencyString := string(dependencyString)
	formatteddependencyString = "\"" + strings.ReplaceAll(formatteddependencyString, "\"", "\\\"") + "\""
	fmt.Println(formatteddependencyString)

	fmt.Println("CorrelationId: ", action.Steps[0].CorrelationId)
	fmt.Println("Token: ", action.Steps[0].Token)

	fmt.Println("Supported Regions: ")
	fmt.Println(action.Steps[0].SupportedRegions)
	supportedRegions := strings.Join(action.Steps[0].SupportedRegions, ",")

	fmt.Println("Target Environment: ")
	fmt.Println(action.Steps[0].TargetEnvironment)
	targetEnvironment := action.Steps[0].TargetEnvironment

	fmt.Println("PackageId: ")
	fmt.Println(action.Steps[0].PackageId)
	packageId := action.Steps[0].PackageId

	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("action", action.Name))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("correlationId", action.Steps[0].CorrelationId))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("token", action.Steps[0].Token))

	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Licenses", formattedlicenseString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Dependencies", formatteddependencyString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("SupportedRegions", supportedRegions))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("TargetEnvironment", targetEnvironment))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("PackageId", packageId))

	uuid := uuid.New()
	var outFilePath = "/cnab/app/" + uuid.String()
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("filePath", outFilePath))

	_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	return err
}
