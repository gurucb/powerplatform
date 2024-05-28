package powerplatform

import (
	"context"
	"fmt"
	"html/template"

	"get.porter.sh/porter/pkg/exec/builder"
	"gopkg.in/yaml.v2"
)

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config MixinConfig
}

// MixinConfig represents configuration that can be set on the powerplatform mixin in porter.yaml
// mixins:
// - powerplatform:
//	  clientVersion: "v0.0.0"

type MixinConfig struct {
	ClientVersion string `yaml:"clientVersion,omitempty"`
}

type buildConfig struct {
	MixinConfig
}

const dockerfileLines = `RUN apt-get update && apt-get install wget -y
RUN apt-get update && apt-get install -y gpg
RUN wget -O - https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor -o microsoft.asc.gpg
RUN mv microsoft.asc.gpg /etc/apt/trusted.gpg.d/
RUN wget https://packages.microsoft.com/config/debian/11/prod.list
RUN mv prod.list /etc/apt/sources.list.d/microsoft-prod.list
RUN chown root:root /etc/apt/trusted.gpg.d/microsoft.asc.gpg
RUN chown root:root /etc/apt/sources.list.d/microsoft-prod.list
RUN apt-get update && apt-get install -y libicu-dev && rm -rf /var/lib/apt/lists/*
 `

// RUN apt-get update &&
// 	apt-get install -y libicu-dev && rm -rf /var/lib/apt/lists/*

// Build will generate the necessary Dockerfile lines
// for an invocation image using this mixin
func (m *Mixin) Build(ctx context.Context) error {

	// Create new Builder.
	var input BuildInput

	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	suppliedClientVersion := input.Config.ClientVersion

	if suppliedClientVersion != "" {
		m.ClientVersion = suppliedClientVersion
	}

	fmt.Fprintln(m.Out, dockerfileLines)
	fmt.Fprintln(m.Out, `ARG GITHUB_TOKEN`)
	fmt.Fprintln(m.Out, `RUN apt-get update && apt-get install -y --no-install-recommends curl unzip`)
	fmt.Fprintln(m.Out, `RUN -H 'Accept: application/vnd.github.v3.raw' curl "https://${GITHUB_TOKEN}@raw.githubusercontent.com/hemantkathuria/privatemixintest/main/mixins/power/v0.0.1/cli/PowerPlatformClient" -o "/cnab/app/PowerPlatformClient"`)
	fmt.Fprintln(m.Out, `RUN chmod 0777 /cnab/app`)
	fmt.Fprintln(m.Out, `RUN echo $PATH`)
	fmt.Fprintln(m.Out, `ENV PATH="$PATH:/cnab/app"`)
	fmt.Fprintln(m.Out, `RUN echo $PATH`)
	fmt.Fprintln(m.Out, `RUN mkdir -p /cnab/app/logs`)
	fmt.Fprintln(m.Out, `RUN chmod 0777 /cnab/app/logs`)
	tmpl, err := template.New("dockerfile").Parse(dockerfileLines)
	if err != nil {
		return fmt.Errorf("error parsing Dockerfile template for the Fabric mixin: %w", err)
	}

	cfg := buildConfig{MixinConfig: input.Config}

	if err = tmpl.Execute(m.Out, cfg); err != nil {
		return fmt.Errorf("error generating Dockerfile lines for the Fabric mixin: %w", err)
	}

	// Example of pulling and defining a client version for your mixin
	// fmt.Fprintf(m.Out, "\nRUN curl https://get.helm.sh/helm-%s-linux-amd64.tar.gz --output helm3.tar.gz", m.ClientVersion)

	return nil
}
