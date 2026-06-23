package supervisor

import (
	"fmt"
	"strings"
)

type App struct {
	Name   string
	Config string

	Children []SupervisionUnit
}

func (app *App) GetName() string {
	return app.Name
}

func (app *App) SetName(name string) error {
	app.Name = name

	return nil
}

func (app *App) Start() error {
	for _, su := range app.Children {
		su.Start()
	}

	return nil
}

func (app *App) Stop(reason ProcessStateReason) error {
	for _, su := range app.Children {
		su.Stop(reason)
	}

	return nil
}

func (app *App) Restart(reason ProcessStateReason) error {
	for _, su := range app.Children {
		su.Restart(reason)
	}

	return nil
}

func (app *App) GetUIString() string {
	result := make([]string, 0)

	result = append(result, app.Name)

	for _, su := range app.Children {
		uiString := su.GetUIString()

		for _, s := range strings.Split(uiString, "\n") {
			result = append(result, fmt.Sprintf("  %s", s))
		}
	}

	return strings.Join(result, "\n")
}

func (app App) FindChild(name string) SupervisionUnit {
	for _, su := range app.Children {
		if su.GetName() == name {
			return su
		}
	}

	return nil
}

func (app App) AppendChild(child SupervisionUnit) error {
	app.Children = append(app.Children, child)

	return nil
}
