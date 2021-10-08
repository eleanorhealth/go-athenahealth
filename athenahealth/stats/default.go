package stats

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Request(method, path string) error {
	return nil
}

func (d *Default) ResponseSuccess() error {
	return nil
}

func (d *Default) ResponseError() error {
	return nil
}
