package view

type ActionView struct {
	Name     string
	Template string
	MimeType string
	Charset  string
	Debug    bool
	Data     any
}
