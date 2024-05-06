package powerplatform

import (
	"get.porter.sh/porter/pkg/exec/builder"
	"get.porter.sh/porter/pkg/runtime"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MarshalYAML converts the action back to a YAML representation
// install:
//
//	powerplatform:
//	  ...
func (a Action) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{a.Name: a.Steps}, nil
}

// MakeSteps builds a slice of Step for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - powerplatform: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Step struct {
	Instruction `yaml:"powerplatform"`
}

// Actions is a set of actions, and the steps, passed from Porter.
type Actions []Action

// UnmarshalYAML takes chunks of a porter.yaml file associated with this mixin
// and populates it on the current action set.
// install:
//
//	powerplatform:
//	  ...
//	powerplatform:
//	  ...
//
// upgrade:
//
//	powerplatform:
//	  ...
func (a *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, Action{})
	if err != nil {
		return err
	}

	for actionName, action := range results {
		for _, result := range action {
			s := result.(*[]Step)
			*a = append(*a, Action{
				Name:  actionName,
				Steps: *s,
			})
		}
	}
	return nil
}

var _ builder.HasOrderedArguments = Instruction{}
var _ builder.ExecutableStep = Instruction{}
var _ builder.StepWithOutputs = Instruction{}

type Instruction struct {
	Description                string                   `yaml:"description"`
	CorrelationId              string                   `yaml:"correlationId"`
	Token                      string                   `yaml:"token"`
	Licenses                   []map[string]interface{} `yaml:"license"`
	Dependencies               []map[string]interface{} `yaml:"dependencies"`
	SupportedRegions           []string                 `yaml:"supportedRegions"`
	TargetEnvironment          string                   `yaml:"targetEnvironment"`
	PackageId                  string                   `yaml:"packageId"`
	Arguments                  []string                 `yaml:"arguments,omitempty"`
	Flags                      builder.Flags            `yaml:"flags,omitempty"`
	builder.IgnoreErrorHandler `yaml:"ignoreError,omitempty"`
	RuntimeConfig              runtime.RuntimeConfig
}

func (s Instruction) GetCommand() string {
	return "PowerPlatformClient"
}

func (s Instruction) GetWorkingDir() string {
	// return s.WorkingDir
	return ""
}

func (s Instruction) GetArguments() []string {
	// return s.Arguments
	return nil
}

func (s Instruction) GetSuffixArguments() []string {
	// return s.SuffixArguments
	return nil
}

func (s Instruction) GetFlags() builder.Flags {
	return s.Flags
}

func (s Instruction) SuppressesOutput() bool {
	// return s.SuppressOutput
	return false
}

func (s Instruction) GetOutputs() []builder.Output {
	// Go doesn't have generics, nothing to see here...
	// outputs := make([]builder.Output, len(s.Outputs))
	// for i := range s.Outputs {
	// 	outputs[i] = s.Outputs[i]
	// }
	// return outputs
	return nil
}

var _ builder.OutputJsonPath = Output{}
var _ builder.OutputFile = Output{}
var _ builder.OutputRegex = Output{}

type Output struct {
	Name string `yaml:"name"`

	// See https://porter.sh/mixins/exec/#outputs
	// TODO: If your mixin doesn't support these output types, you can remove these and the interface assertions above, and from #/definitions/outputs in schema.json
	JsonPath string `yaml:"jsonPath,omitempty"`
	FilePath string `yaml:"path,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}

func (o Output) GetFilePath() string {
	return o.FilePath
}

func (o Output) GetRegex() string {
	return o.Regex
}
