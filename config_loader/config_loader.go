package config_loader

import (
	"fmt"
	"maps"
	"os/exec"

	"github.com/ne006/shepherd/supervisor"
	"go.uber.org/zap"

	"gopkg.in/yaml.v3"

	"os"
)

type config struct {
	App struct {
		Name     string         `yaml:"name"`
		Children []interface{}  `yaml:"children"`
		Env      map[string]any `yaml:"env"`
	}
}

func LoadConfig(path string, logger *zap.SugaredLogger) (*supervisor.App, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	var app supervisor.App

	app.Name = cfg.App.Name
	app.Config = path

	app.SetEnv(cfg.App.Env)

	if appChildren, err := loadChildren(cfg.App.Children, &app, logger); err != nil {
		return nil, err
	} else {
		app.Children = appChildren
	}

	return &app, nil
}

func loadChildren(sourceList []interface{}, parentUnit supervisor.Composite, logger *zap.SugaredLogger) ([]supervisor.SupervisionUnit, error) {
	var children []supervisor.SupervisionUnit

	for _, source := range sourceList {
		if child, err := loadChild(source, parentUnit, logger); err != nil {
			return children, err
		} else if child != supervisor.SupervisionUnit(nil) {
			children = append(children, child)
		}
	}

	return children, nil
}

func loadChild(source interface{}, parentUnit supervisor.Composite, logger *zap.SugaredLogger) (supervisor.SupervisionUnit, error) {
	if sourceMap, ok := source.(map[string]interface{}); ok {
		if sourceChildren, ok := sourceMap["children"].([]interface{}); ok {
			return loadGroup(sourceMap, sourceChildren, parentUnit, logger)
		} else if sourceCommand, ok := sourceMap["command"].(string); ok {
			return loadProcess(sourceMap, sourceCommand, parentUnit, logger)
		} else {
			return nil, fmt.Errorf("could not deduce type of %v", sourceMap)
		}
	} else {
		return nil, fmt.Errorf("%v should be a map", source)
	}
}

func loadGroup(sourceMap map[string]interface{}, sourceChildren []interface{}, parentUnit supervisor.Composite, logger *zap.SugaredLogger) (supervisor.SupervisionUnit, error) {
	var name string
	var ok bool

	if name, ok = sourceMap["name"].(string); !ok {
		return nil, fmt.Errorf("%v should be a string", sourceMap["name"])
	}

	unit := supervisor.Group{
		Name: name,
	}

	loadEnv(&unit, sourceMap, parentUnit)

	if unitChildren, err := loadChildren(sourceChildren, &unit, logger); err != nil {
		return nil, err
	} else {
		unit.Children = unitChildren
	}

	return &unit, nil
}

func loadProcess(sourceMap map[string]interface{}, sourceCommand string, parentUnit supervisor.Composite, logger *zap.SugaredLogger) (supervisor.SupervisionUnit, error) {
	var name string
	var ok bool

	if name, ok = sourceMap["name"].(string); !ok {
		return nil, fmt.Errorf("%v should be a string", sourceMap["name"])
	}

	unit := supervisor.Process{
		Name: name,

		Cmd: exec.Cmd{
			Path: sourceCommand,
		},

		Logger: logger,
	}

	if unitArgs, err := loadArgs(sourceMap["args"]); err != nil {
		return nil, err
	} else {
		unit.Args = unitArgs
	}

	loadEnv(&unit, sourceMap, parentUnit)

	return &unit, nil
}

func loadArgs(sourceArgs interface{}) ([]string, error) {
	var unitArgs []string

	if sourceArgs == nil {
		return []string{}, nil
	}

	if rawArgs, ok := sourceArgs.([]interface{}); ok {
		for _, arg := range rawArgs {
			if stringArg, ok := arg.(string); ok {
				unitArgs = append(unitArgs, stringArg)
			} else {
				return nil, fmt.Errorf("%v should be a string, got %T", arg, arg)
			}
		}

		return unitArgs, nil
	} else {
		return nil, fmt.Errorf("%v should be a list, got %T", sourceArgs, sourceArgs)
	}
}

func loadEnv(unit supervisor.SupervisionUnit, sourceMap map[string]interface{}, parentUnit supervisor.Composite) {
	if eo, ok := unit.(supervisor.EnvironmentOwner); ok {
		currentEnv := make(map[string]any)

		if peo, ok := parentUnit.(supervisor.EnvironmentOwner); ok {
			maps.Copy(currentEnv, peo.GetEnv())
		}

		if unitEnv, ok := sourceMap["env"].(map[string]any); ok {
			maps.Copy(currentEnv, unitEnv)
		}

		eo.SetEnv(currentEnv)
	}
}
