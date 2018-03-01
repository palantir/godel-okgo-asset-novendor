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

package novendor

import (
	"io"
	"path"

	"github.com/go-yaml/yaml"
	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"
	"github.com/pkg/errors"
)

const (
	TypeName okgo.CheckerType     = "novendor"
	Priority okgo.CheckerPriority = 0
)

func Creator() checker.Creator {
	return checker.NewCreator(
		TypeName,
		Priority,
		func(cfgYML []byte) (okgo.Checker, error) {
			var cfg novendorCheckCfg
			if err := yaml.Unmarshal(cfgYML, &cfg); err != nil {
				return nil, errors.Wrapf(err, "failed to unmarshal configuration YAML %q", string(cfgYML))
			}
			return &novendorCheck{
				pkgRegexps:                cfg.PkgRegexps,
				includeVendorInImportPath: cfg.IncludeVendorInImportPath,
				ignorePkgs:                cfg.IgnorePkgs,
			}, nil
		},
	)
}

type novendorCheck struct {
	pkgRegexps                []string
	includeVendorInImportPath bool
	ignorePkgs                []string
}

type novendorCheckCfg struct {
	PkgRegexps                []string `yaml:"pkg-regexps"`
	IncludeVendorInImportPath bool     `yaml:"include-vendor-in-import-path"`
	IgnorePkgs                []string `yaml:"ignore-pkgs"`
}

func (c *novendorCheck) Type() (okgo.CheckerType, error) {
	return TypeName, nil
}

func (c *novendorCheck) Priority() (okgo.CheckerPriority, error) {
	return Priority, nil
}

func (c *novendorCheck) Check(pkgPaths []string, projectDir string, stdout io.Writer) {
	var args []string
	for _, regexp := range c.pkgRegexps {
		args = append(args, "--pkg-regexp", regexp)
	}
	if c.includeVendorInImportPath {
		args = append(args, "--full-import-path")
	}
	for _, pkg := range c.ignorePkgs {
		args = append(args, "--ignore-pkg", path.Join(projectDir, pkg))
	}
	cmd, wd := checker.AmalgomatedCheckCmd(string(TypeName), append(args, pkgPaths...), stdout)
	if cmd == nil {
		return
	}
	checker.RunCommandAndStreamOutput(cmd, func(line string) okgo.Issue {
		return okgo.NewIssueFromLine(line, wd)
	}, stdout)
}

func (c *novendorCheck) RunCheckCmd(args []string, stdout io.Writer) {
	checker.AmalgomatedRunRawCheck(string(TypeName), args, stdout)
}
