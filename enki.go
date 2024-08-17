package enki

const version = "0.0.1"

type Enki struct {
	AppName string
}

func New(name string) Enki {
	app := Enki{}

	app.AppName = name

	return app
}

func (enki *Enki) Version() string {
	return version
}
