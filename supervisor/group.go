package supervisor

import (
	"fmt"
	"strings"
)

type Group struct {
	Name     string
	Children []SupervisionUnit

	envHash map[string]any
}

func (group *Group) GetName() string {
	return group.Name
}

func (group *Group) SetName(name string) error {
	group.Name = name

	return nil
}

func (group *Group) GetEnv() map[string]any {
	return group.envHash
}

func (group *Group) SetEnv(envHash map[string]any) {
	group.envHash = envHash
}

func (group *Group) Start() error {
	for _, su := range group.Children {
		su.Start()
	}

	return nil
}

func (group *Group) Stop(reason ProcessStateReason) error {
	for _, su := range group.Children {
		su.Stop(reason)
	}

	return nil
}

func (group *Group) Restart(reason ProcessStateReason) error {
	for _, su := range group.Children {
		su.Restart(reason)
	}

	return nil
}

func (group *Group) GetUIString() string {
	result := make([]string, 0)

	result = append(result, group.Name)

	for _, su := range group.Children {
		uiString := su.GetUIString()

		for _, s := range strings.Split(uiString, "\n") {
			result = append(result, fmt.Sprintf("  %s", s))
		}
	}

	return strings.Join(result, "\n")
}

func (group Group) FindChild(name string) SupervisionUnit {
	for _, su := range group.Children {
		if su.GetName() == name {
			return su
		}
	}

	return nil
}

func (group Group) AppendChild(child SupervisionUnit) error {
	group.Children = append(group.Children, child)

	return nil
}
