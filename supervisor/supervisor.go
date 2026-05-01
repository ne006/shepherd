package supervisor

type Supervisor struct {
	svtree SupervisionTree
}

func (s *Supervisor) LoadApp(app App) error {
	if err := s.svtree.Append(app); err != nil {
		return err
	}

	return nil
}

func (s *Supervisor) Start() error {
	return s.svtree.Start()
}

func (s *Supervisor) Stop() error {
	return s.svtree.Stop()
}

func (s *Supervisor) GetUIString() string {
	return s.svtree.GetUIString()
}
