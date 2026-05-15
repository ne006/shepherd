package supervisor

import (
	"strings"
)

type SupervisionTree struct {
	Children    []App
	OldChildren []App
}

func (stree *SupervisionTree) Start() error {
	if err := stree.StopOld(); err != nil {
		return err
	}

	for _, app := range stree.Children {
		app.Start()
	}

	return nil
}

func (stree *SupervisionTree) Stop() error {
	for _, app := range stree.Children {
		app.Stop(StateReasonUser)
	}

	return nil
}

func (stree *SupervisionTree) StopOld() error {
	for _, app := range stree.OldChildren {
		app.Stop(StateReasonUser)
	}

	stree.OldChildren = []App{}

	return nil
}

func (stree *SupervisionTree) GetUIString() string {
	result := make([]string, 0)

	for _, su := range stree.Children {
		result = append(result, su.GetUIString())
	}

	return strings.Join(result, "\n")
}

func (stree *SupervisionTree) Append(newApp App) error {
	var oldIdx int
	var oldApp *App

	for i, app := range stree.Children {
		if app.GetName() == newApp.GetName() {
			oldIdx, oldApp = i, &app
		}
	}

	if oldApp != nil {
		stree.OldChildren = append(stree.OldChildren, *oldApp)

		stree.Children = append(
			stree.Children[:oldIdx],
			append(
				[]App{newApp},
				stree.Children[oldIdx+1:]...,
			)...,
		)
	} else {
		stree.Children = append(stree.Children, newApp)
	}

	return nil
}

func (stree *SupervisionTree) FindChild(name string) *App {
	for _, app := range stree.Children {
		if app.Name == name {
			return &app
		}
	}

	return nil
}
