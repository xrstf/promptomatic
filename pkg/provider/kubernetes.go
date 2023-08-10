// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import (
	"os"
	"path/filepath"
	"strings"

	"go.xrstf.de/promptomatic/pkg/sys"
)

type KubernetesEnv string

const (
	KubernetesEnvDev   KubernetesEnv = "dev"
	KubernetesEnvProd  KubernetesEnv = "prod"
	KubernetesEnvCI    KubernetesEnv = "ci"
	KubernetesEnvOther KubernetesEnv = "other"
)

type KubernetesStatus struct {
	Kubeconfig string
	Identifier string
	Symbol     string
	Context    string
	Env        KubernetesEnv
}

type kubeconfigStyle struct {
	Symbol     string
	Identifier string
	Env        KubernetesEnv
}

var kubeconfigs = map[string]kubeconfigStyle{
	"~/loodse/clusters/dev/kubeconfig":          {Symbol: "☸", Identifier: "dev", Env: KubernetesEnvDev},
	"~/loodse/clusters/captain/kubeconfig":      {Symbol: "☢", Identifier: "captain", Env: KubernetesEnvProd},
	"~/loodse/clusters/ci/kubeconfig":           {Symbol: "", Identifier: "ci", Env: KubernetesEnvCI},
	"~/loodse/clusters/ci-worker/kubeconfig":    {Symbol: "", Identifier: "ci-worker", Env: KubernetesEnvCI},
	"~/loodse/clusters/dc-ci-worker/kubeconfig": {Symbol: "", Identifier: "dc-ci-worker", Env: KubernetesEnvCI},
}

func NewKubernetesStatus(e *sys.Environment) (*KubernetesStatus, error) {
	if e.Kubeconfig == "" {
		return nil, nil
	}

	if !sys.CommandExists("kubectl") {
		return nil, nil
	}

	kubecontext := os.Getenv("KUBECTX") // KUBECTX is only valid in my own kubectl wrapper script

	for path, example := range kubeconfigs {
		path = strings.Replace(path, "~", e.HomeDirectory, 1)

		if e.Kubeconfig == path {
			return &KubernetesStatus{
				Kubeconfig: e.Kubeconfig,
				Symbol:     example.Symbol,
				Identifier: example.Identifier,
				Env:        example.Env,
				Context:    kubecontext,
			}, nil
		}
	}

	s := &KubernetesStatus{
		Env:        KubernetesEnvOther,
		Kubeconfig: e.Kubeconfig,
		Context:    kubecontext,
	}

	// does KUBECONFIG point to an existing file?
	if _, err := os.Stat(e.Kubeconfig); err == nil {
		// shorten the path as much as possible
		abs, err := filepath.Abs(e.Kubeconfig)
		if err == nil {
			// determine the path relative to the home directory
			homeRelative := abs
			if strings.HasPrefix(abs, e.HomeDirectory) {
				homeRelative = strings.Replace(abs, e.HomeDirectory, "~", 1)
			}

			// determine the path relative to the current directory
			pwdRelative := abs
			if strings.HasPrefix(abs, e.WorkingDirectory) {
				pwdRelative = strings.Replace(abs, e.WorkingDirectory, ".", 1)
			}

			shortest := abs
			if len(homeRelative) < len(shortest) {
				shortest = homeRelative
			}
			if len(pwdRelative) < len(shortest) {
				shortest = pwdRelative
			}

			s.Kubeconfig = shortest
		}
	}

	return s, nil
}
