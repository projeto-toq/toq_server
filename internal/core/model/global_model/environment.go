package globalmodel

type Environment struct {
	GRPC struct {
		Port              string `yaml:"port"`
		Network           string `yaml:"network"`
		CertPath          string `yaml:"cert_path"`
		KeyPath           string `yaml:"key_path"`
		ClientCAPath      string `yaml:"client_ca_path"`
		RequireClientCert bool   `yaml:"require_client_cert"`
	}
	DB struct {
		URI string `yaml:"uri"`
	}
	DATABASE struct {
		Populate bool `yaml:"populate"`
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
	AUTH struct {
		RefreshTTLDays                int `yaml:"refresh_ttl_days"`
		AccessTTLMinutes              int `yaml:"access_ttl_minutes"`
		MaxSessionRotations           int `yaml:"max_session_rotations"`
		SessionCleanerIntervalSeconds int `yaml:"session_cleaner_interval_seconds"`
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
	FCM struct {
		CredentialsFile string `yaml:"credentials_file"`
		ProjectID       string `yaml:"project_id"`
	}
	GCS struct {
		ProjectID     string `yaml:"project_id"`
		AdminSAEmail  string `yaml:"admin_sa_email"`
		WriterSAEmail string `yaml:"writer_sa_email"`
		ReaderSAEmail string `yaml:"reader_sa_email"`
		AdminCreds    string `yaml:"admin_creds_path"`
		WriterCreds   string `yaml:"writer_creds_path"`
		ReaderCreds   string `yaml:"reader_creds_path"`
	}
	HEALTH struct {
		HTTPPort int    `yaml:"http_port"`
		UseTLS   bool   `yaml:"use_tls"`
		CertPath string `yaml:"cert_path"`
		KeyPath  string `yaml:"key_path"`
	}
}
