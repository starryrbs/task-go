package conf

type (
	Database struct {
		Source string
		Driver string
	}

	Server struct {
		Address string `yaml:"address"`
	}
	Conf struct {
		Database Database
		Server   Server
	}
)
