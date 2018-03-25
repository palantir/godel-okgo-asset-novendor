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

package legacy

import (
	"github.com/palantir/godel/pkg/versionedconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel-okgo-asset-novendor/novendor/config/internal/v0"
)

type Config struct {
	versionedconfig.ConfigWithLegacy `yaml:",inline"`
	Args                             []string `yaml:"args"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg Config
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal novendor-asset legacy configuration")
	}
	if len(legacyCfg.Args) == 0 {
		return nil, nil
	}

	var upgradedCfg v0.Config
	if legacyCfg.Args[0] != "--ignore" {
		return nil, errors.Errorf(`novendor-asset only supports legacy configuration if the first element in "args" is "--ignore"`)
	}
	for _, currIgnore := range legacyCfg.Args[1:] {
		upgradedCfg.IgnorePkgs = append(upgradedCfg.IgnorePkgs, currIgnore)
	}

	if len(upgradedCfg.IgnorePkgs) == 0 {
		// if "ignore" field is blank, same as having blank config
		return nil, nil
	}

	upgradedCfgBytes, err := yaml.Marshal(upgradedCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal novendor-asset v0 configuration")
	}
	return upgradedCfgBytes, nil
}
