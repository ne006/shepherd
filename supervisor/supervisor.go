package supervisor

import (
	"fmt"
)

type Supervisor struct {
	svtree SupervisionTree
}

func (s *Supervisor) LoadApp(app App) error {
	if err := s.svtree.AppendChild(&app); err != nil {
		return err
	}

	return nil
}

func (s *Supervisor) Start(spath string) error {
	if child := s.FindChild(spath); child != nil {
		return (*child).Start()
	} else {
		return fmt.Errorf("%s is not defined", spath)
	}
}

func (s *Supervisor) Stop(spath string) error {
	if child := s.FindChild(spath); child != nil {
		return (*child).Stop(StateReasonUser)
	} else {
		return fmt.Errorf("%s is not defined", spath)
	}
}

func (s *Supervisor) Restart(spath string) error {
	if child := s.FindChild(spath); child != nil {
		return (*child).Restart(StateReasonUser)
	} else {
		return fmt.Errorf("%s is not defined", spath)
	}
}

func (s *Supervisor) GetUIString() string {
	return s.svtree.GetUIString()
}

func (s *Supervisor) FindChild(spath string) *SupervisionUnit {
	path := splitPath(spath)

	var current SupervisionUnit

	current = &s.svtree

	if spath == "" {
		return &current // Return root
	}

	for i, part := range path {
		if current == nil {
			break
		}

		if i == len(path) {
			break // Node matching last path part is the searched one
		} else {
			if currentComposite, isComposite := current.(Composite); isComposite {
				next := currentComposite.FindChild(part)

				if next != nil {
					current = next // Going down the tree along the path
				} else {
					return nil // Still need to go down the tree and the next node is nil
				}
			} else {
				return nil // Still need to go down the tree and the next node is a leaf
			}
		}
	}

	return &current
}
