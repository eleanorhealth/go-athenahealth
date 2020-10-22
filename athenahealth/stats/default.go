package tokencacher

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) IncrRequests() error {
	return nil
}
