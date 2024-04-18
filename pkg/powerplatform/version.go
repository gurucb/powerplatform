package powerplatform

import (
	"get.porter.sh/porter/pkg/mixin"
	"get.porter.sh/porter/pkg/pkgmgmt"
	"get.porter.sh/porter/pkg/porter/version"
	"github.com/getporter/powerplatform/pkg"
)

func (m *Mixin) PrintVersion(opts version.Options) error {
	metadata := mixin.Metadata{
		Name: "powerplatform",
		VersionInfo: pkgmgmt.VersionInfo{
			Version: pkg.Version,
			Commit:  pkg.Commit,
			Author:  "gurucb",
		},
	}
	return version.PrintVersion(m.Context, opts, metadata)
}
