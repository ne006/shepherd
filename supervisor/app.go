package supervisor

type App struct {
	Name   string
	Config string

	Children []SupervisionUnit
}

func (app App) GetName() string {
	return app.Name
}

func (app App) SetName(name string) error {
	app.Name = name

	return nil
}

func (app App) Start() error {
	for _, su := range app.Children {
		su.Start()
	}

	return nil
}

func (app App) Stop() error {
	for _, su := range app.Children {
		su.Stop()
	}

	return nil
}

func (app App) FindChild(name string) *SupervisionUnit {
	for _, su := range app.Children {
		if su.GetName() == name {
			return &su
		}
	}

	return nil
}

func (app App) AppendChild(child *SupervisionUnit) error {
	app.Children = append(app.Children, *child)

	return nil
}
