package config_loader

import (
	"fmt"
	"os/exec"

	"github.com/ne006/shepherd/supervisor"

	"gopkg.in/yaml.v3"

	"os"
)

type config struct {
	App struct {
		Name     string        `yaml:"name"`
		Children []interface{} `yaml:"children"`
	}
}

func LoadConfig(path string) (*supervisor.App, error) {
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

	if appChildren, err := loadChildren(cfg.App.Children); err != nil {
		return nil, err
	} else {
		app.Children = appChildren
	}

	return &app, nil
}

func loadChildren(sourceList []interface{}) ([]supervisor.SupervisionUnit, error) {
	var children []supervisor.SupervisionUnit

	for _, source := range sourceList {
		if child, err := loadChild(source); err != nil {
			return children, err
		} else if child != supervisor.SupervisionUnit(nil) {
			children = append(children, child)
		}
	}

	return children, nil
}

func loadChild(source interface{}) (supervisor.SupervisionUnit, error) {
	if sourceMap, ok := source.(map[string]interface{}); ok {
		if sourceChildren, ok := sourceMap["children"].([]interface{}); ok {
			return loadGroup(sourceMap, sourceChildren)
		} else if sourceCommand, ok := sourceMap["command"].(string); ok {
			return loadProcess(sourceMap, sourceCommand)
		} else {
			return nil, fmt.Errorf("could not deduce type of %v", sourceMap)
		}
	} else {
		return nil, fmt.Errorf("%v should be a map", source)
	}
}

func loadGroup(sourceMap map[string]interface{}, sourceChildren []interface{}) (supervisor.SupervisionUnit, error) {
	var name string
	var ok bool

	if name, ok = sourceMap["name"].(string); !ok {
		return nil, fmt.Errorf("%v should be a string", sourceMap["name"])
	}

	unit := supervisor.Group{
		Name: name,
	}

	if unitChildren, err := loadChildren(sourceChildren); err != nil {
		return nil, err
	} else {
		unit.Children = unitChildren
	}

	return unit, nil
}

func loadProcess(sourceMap map[string]interface{}, sourceCommand string) (supervisor.SupervisionUnit, error) {
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
	}

	if unitArgs, err := loadArgs(sourceMap["args"]); err != nil {
		return nil, err
	} else {
		unit.Args = unitArgs
	}

	return unit, nil
}

func loadArgs(sourceArgs interface{}) ([]string, error) {
	var unitArgs []string

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
