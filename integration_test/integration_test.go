// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/okgo/okgotester"
	"github.com/stretchr/testify/require"
)

const (
	okgoPluginLocator  = "com.palantir.okgo:check-plugin:1.0.0-rc7"
	okgoPluginResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
)

func TestCheck(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	assetPath, err := products.Bin("novendor-asset")
	require.NoError(t, err)

	configFiles := map[string]string{
		"godel/config/godel.yml":        godelYML,
		"godel/config/check-plugin.yml": "",
	}

	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(okgoPluginLocator, okgoPluginResolver)
	require.NoError(t, err)

	okgotester.RunAssetCheckTest(t,
		pluginProvider,
		pluginapitester.NewAssetProvider(assetPath),
		"novendor",
		".",
		[]okgotester.AssetTestCase{
			{
				Name: "novendor failures",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src:     `package foo`,
					},
					{
						RelPath: "vendor/github.com/org/repo/bar/bar.go",
						Src:     `package bar`,
					},
				},
				ConfigFiles: configFiles,
				WantError:   true,
				WantOutput: `Running novendor...
github.com/org/repo
Finished novendor
Check(s) produced output: [novendor]
`,
			},
			{
				Name: "novendor failures from inner directory",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src:     `package foo`,
					},
					{
						RelPath: "vendor/github.com/org/repo/bar/bar.go",
						Src:     `package bar`,
					},
					{
						RelPath: "inner/bar",
					},
				},
				ConfigFiles: configFiles,
				Wd:          "inner",
				WantError:   true,
				WantOutput: `Running novendor...
github.com/org/repo
Finished novendor
Check(s) produced output: [novendor]
`,
			},
			{
				Name: "check filters work for novendor output",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src:     `package foo`,
					},
					{
						RelPath: "vendor/github.com/org/repo/bar/bar.go",
						Src:     `package bar`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/check-plugin.yml": `
checks:
  novendor:
    filters:
    - value: github.com/org/repo
`,
				},
				WantError: false,
				WantOutput: `Running novendor...
Finished novendor
`,
			},
		},
	)
}

func TestUpgradeConfig(t *testing.T) {
	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(okgoPluginLocator, okgoPluginResolver)
	require.NoError(t, err)

	assetPath, err := products.Bin("novendor-asset")
	require.NoError(t, err)
	assetProvider := pluginapitester.NewAssetProvider(assetPath)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginProvider,
		[]pluginapitester.AssetProvider{assetProvider},
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: `legacy configuration with empty "args" field is updated`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  novendor:
    filters:
      - value: "should have comment or be unexported"
      - type: name
        value: ".*.pb.go"
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for check-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `checks:
  novendor:
    filters:
    - value: should have comment or be unexported
    exclude:
      names:
      - .*.pb.go
`,
				},
			},
			{
				Name: `legacy configuration with "ignore" args is upgraded`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  novendor:
    args:
      - "--ignore"
      - "./vendor/github.com/palantir/go-novendor"
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for check-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `checks:
  novendor:
    config:
      ignore-pkgs:
      - ./vendor/github.com/palantir/go-novendor
`,
				},
			},
			{
				Name: `legacy configuration with args other than "ignore" fails`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  novendor:
    args:
      - "-help"
`,
				},
				Legacy:    true,
				WantError: true,
				WantOutput: `Failed to upgrade configuration:
	godel/config/check-plugin.yml: failed to upgrade configuration: failed to upgrade check "novendor" legacy configuration: failed to upgrade asset configuration: novendor-asset only supports legacy configuration if the first element in "args" is "--ignore"
`,
				WantFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  novendor:
    args:
      - "-help"
`,
				},
			},
			{
				Name: `valid v0 config works`,
				ConfigFiles: map[string]string{
					"godel/config/check-plugin.yml": `
checks:
  novendor:
    config:
      # comment
      ignore-pkgs:
        - "./vendor/github.com/palantir/go-novendor"
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `
checks:
  novendor:
    config:
      # comment
      ignore-pkgs:
        - "./vendor/github.com/palantir/go-novendor"
`,
				},
			},
		},
	)
}
