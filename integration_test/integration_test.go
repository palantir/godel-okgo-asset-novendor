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
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/okgo/okgotester"
	"github.com/stretchr/testify/require"
)

const (
	okgoPluginLocator  = "com.palantir.okgo:okgo-plugin:0.3.0"
	okgoPluginResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

	godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`
)

func TestNovendor(t *testing.T) {
	assetPath, err := products.Bin("novendor-asset")
	require.NoError(t, err)

	configFiles := map[string]string{
		"godel/config/godel.yml": godelYML,
		"godel/config/check.yml": "",
	}

	okgotester.RunAssetCheckTest(t,
		okgoPluginLocator, okgoPluginResolver,
		assetPath, "novendor",
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
`,
			},
		},
	)
}
