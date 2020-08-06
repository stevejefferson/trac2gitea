package markdown

type Converter struct {
}

func CreateConverter() *Converter {
	return nil
}

func (converter *Converter) Convert(tracText string) string {
	return ""
}

func (converter *Converter) ConvertLink(tracText string, linkPrefix string) string {
	return ""
}
