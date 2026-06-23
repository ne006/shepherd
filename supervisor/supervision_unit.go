package supervisor

type SupervisionUnit interface {
	Start() error
	Stop(reason ProcessStateReason) error
	Restart(reason ProcessStateReason) error

	GetName() string
	SetName(string) error

	GetUIString() string
}
