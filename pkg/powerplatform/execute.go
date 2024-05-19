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
	formattedlicenseString = "\\" + strings.ReplaceAll(formattedlicenseString, "\"", "\\") + "\\"
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
	supportedRegions := strings.Join(action.Steps[0].SupportedRegions, ",")

	targetEnvironment := action.Steps[0].TargetEnvironment

	packageId := action.Steps[0].PackageId

	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("action", action.Name))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("correlationId", action.Steps[0].CorrelationId))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("token", action.Steps[0].Token))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Licenses", formattedlicenseString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("Dependencies", formatteddependencyString))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("SupportedRegions", supportedRegions))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("TargetEnvironment", targetEnvironment))
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("PackageId", packageId))
	// var output string
	uuid := uuid.New()
	var outFilePath = "/cnab/app/" + uuid.String()
	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("filePath", outFilePath))

	_, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	// Comment by Guru (Action.go: 115 Step), make below lines active once we work on outputs from PowerPlatform provisioning API

	// if _, err := os.Stat(outFilePath); os.IsNotExist(err) {
	// 	fmt.Println("Output file does not exist")
	// 	return err
	// }

	// fmt.Println("Output:")
	// fmt.Println(output)

	// executedStep := action.Steps[0]
	// outputData, err := os.ReadFile(outFilePath)
	// if len(executedStep.Instruction.Outputs) > 0 {
	// 	var instructionOutput = InstructionOutput{Name: executedStep.Instruction.Name, Outputs: executedStep.Instruction.Outputs}
	// 	builder.ProcessJsonPathOutputs(ctx, m.RuntimeConfig, instructionOutput, string(outputData))
	// }

	return err
}

type InstructionOutput struct {
	Name    string   `yaml:"name"`
	Outputs []Output `yaml:"outputs,omitempty"`
}

func (s InstructionOutput) GetOutputs() []builder.Output {
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}
