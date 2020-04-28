package image

type IGenerator interface {
	Generate(text, name string) (string, error)
}
