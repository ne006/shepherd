package supervisor

type Composite interface {
	FindChild(string) SupervisionUnit
	AppendChild(SupervisionUnit) error
}
