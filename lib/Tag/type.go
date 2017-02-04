package tag

type Tag string

func (t *Tag) New(s string) error {
	*t = Tag(s)
	return nil
}
