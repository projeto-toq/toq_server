package globalmodel

type Environment struct {
	HTTP struct {
		Port           string `yaml:"port"`
		Network        string `yaml:"network"`
		ReadTimeout    string `yaml:"read_timeout"`
		WriteTimeout   string `yaml:"write_timeout"`
		MaxHeaderBytes int    `yaml:"max_header_bytes"`
		TLS            struct {
			Enabled  bool   `yaml:"enabled"`
			CertPath string `yaml:"cert_path"`
			KeyPath  string `yaml:"key_path"`
		} `yaml:"tls"`
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
	TELEMETRY struct {
		Enabled bool `yaml:"enabled"`
		OTLP    struct {
			Enabled  bool   `yaml:"enabled"`
			Endpoint string `yaml:"endpoint"`
			Insecure bool   `yaml:"insecure"`
		} `yaml:"otlp"`
		METRICS struct {
			Port string `yaml:"port"`
		} `yaml:"metrics"`
	}
	JWT struct {
		Secret string `yaml:"secret"`
	}
	AUTH struct {
		RefreshTTLDays                   int `yaml:"refresh_ttl_days"`
		AccessTTLMinutes                 int `yaml:"access_ttl_minutes"`
		MaxSessionRotations              int `yaml:"max_session_rotations"`
		SessionCleanerIntervalSeconds    int `yaml:"session_cleaner_interval_seconds"`
		ValidationCleanerIntervalSeconds int `yaml:"validation_cleaner_interval_seconds"`
		// StatusRulesPath                  string `yaml:"status_rules_path"`
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
		UseTLS       bool   `yaml:"use_tls"`
		UseSSL       bool   `yaml:"use_ssl"`
		SkipVerify   bool   `yaml:"skip_verify"`
		FromEmail    string `yaml:"from_email"`
		FromName     string `yaml:"from_name"`
		MaxRetries   int    `yaml:"max_retries"`
		TimeoutSecs  int    `yaml:"timeout_seconds"`
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
	S3 struct {
		Region     string `yaml:"region"`
		BucketName string `yaml:"bucket_name"`
		AdminRole  struct {
			AccessKeyID     string `yaml:"access_key_id"`
			SecretAccessKey string `yaml:"secret_access_key"`
		} `yaml:"admin"`
		ReaderRole struct {
			AccessKeyID     string `yaml:"access_key_id"`
			SecretAccessKey string `yaml:"secret_access_key"`
		} `yaml:"reader"`
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
	// Health endpoints are now integrated into the main HTTP server
	// No separate health configuration needed
}
