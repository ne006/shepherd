package supervisor

import (
	"strings"
)

type SupervisionTree struct {
	Children    []SupervisionUnit
	OldChildren []SupervisionUnit
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

	stree.OldChildren = []SupervisionUnit{}

	return nil
}

func (stree *SupervisionTree) GetUIString() string {
	result := make([]string, 0)

	for _, su := range stree.Children {
		result = append(result, su.GetUIString())
	}

	return strings.Join(result, "\n")
}

func (stree *SupervisionTree) AppendChild(su SupervisionUnit) error {
	if _, isApp := su.(*App); !isApp {
		return fmt.Errorf("only an App can be a direct child of SupervisionTree")
	}

	var oldIdx int
	var oldApp *SupervisionUnit

	for i, app := range stree.Children {
		if app.GetName() == su.GetName() {
			oldIdx, oldApp = i, &app
		}
	}

	if oldApp != nil {
		stree.OldChildren = append(stree.OldChildren, *oldApp)

		stree.Children = append(
			stree.Children[:oldIdx],
			append(
				[]SupervisionUnit{su},
				stree.Children[oldIdx+1:]...,
			)...,
		)
	} else {
		stree.Children = append(stree.Children, su)
	}

	return nil
}

func (stree *SupervisionTree) FindChild(name string) SupervisionUnit {
	for _, app := range stree.Children {
		if app.GetName() == name {
			return app
		}
	}

	return nil
}
