package supervisor

type SupervisionUnit interface {
	Start() error
	Stop() error

	GetName() string
	SetName(string) error

	GetUIString() string
}
