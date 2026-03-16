package supervisor

type Group struct {
	Name     string
	Children []SupervisionUnit
}

func (group *Group) GetName() string {
	return group.Name
}

func (group *Group) SetName(name string) error {
	group.Name = name

	return nil
}

func (group *Group) Start() error {
	for _, su := range group.Children {
		su.Start()
	}

	return nil
}

func (group *Group) Stop() error {
	for _, su := range group.Children {
		su.Stop()
	}

	return nil
}

func (group Group) FindChild(name string) *SupervisionUnit {
	for _, su := range group.Children {
		if su.GetName() == name {
			return &su
		}
	}

	return nil
}

func (group Group) AppendChild(child *SupervisionUnit) error {
	group.Children = append(group.Children, *child)

	return nil
}
