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
