package supervisor

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type EnvironmentOwner interface {
	GetEnv() map[string]any
	SetEnv(map[string]any)
}

func loadEnv(eo EnvironmentOwner) []string {
	var env = make([]string, 0)

	funcMap := template.FuncMap{
		"env": os.Getenv,
	}

	tmplWithFuncs := template.New("env").Funcs(funcMap)

	for k, v := range eo.GetEnv() {
		envStr := fmt.Sprintf("%s=%v", k, v)
		var buffer strings.Builder

		if tmpl, err := tmplWithFuncs.Parse(envStr); err == nil {
			tmpl.Execute(&buffer, nil)

			env = append(env, buffer.String())
		}
	}

	return env
}
