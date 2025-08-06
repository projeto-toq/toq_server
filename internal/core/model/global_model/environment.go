package globalmodel

type Environment struct {
	GRPC struct {
		Port    string `yaml:"port"`
		Network string `yaml:"network"`
	}
	DB struct {
		URI string `yaml:"uri"`
	}
	REDIS struct {
		URL string `yaml:"url"`
	}
	LOG struct {
		Level     string `yaml:"level"`
		AddSource bool   `yaml:"addsource"`
		ToFile    bool   `yaml:"tofile"`
		Path      string `yaml:"path"`
		Filename  string `yaml:"filename"`
	}
	JWT struct {
		Secret string `yaml:"secret"`
	}
	GMAIL struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	SMS struct {
		AccountSid string `yaml:"account_sid"`
		AuthToken  string `yaml:"auth_token"`
		MyNumber   string `yaml:"my_number"`
	}
	EMAIL struct {
		SMTPServer   string `yaml:"smtp_server"`
		SMTPPort     int    `yaml:"smtp_port"`
		SMTPUser     string `yaml:"smtp_user"`
		SMTPPassword string `yaml:"smtp_password"`
	}
	CNPJ struct {
		Token   string `yaml:"token"`
		URLBase string `yaml:"url_base"`
	}
	CEP struct {
		Token   string `yaml:"token"`
		URLBase string `yaml:"url_base"`
	}
	CPF struct {
		Token   string `yaml:"token"`
		URLBase string `yaml:"url_base"`
	}
}
