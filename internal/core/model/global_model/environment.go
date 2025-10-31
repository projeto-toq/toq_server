package globalmodel

type Environment struct {
	APP struct {
		Debug bool `yaml:"debug"`
	} `yaml:"app"`
	HTTP struct {
		Port           string `yaml:"port"`
		Network        string `yaml:"network"`
		ReadTimeout    string `yaml:"read_timeout"`
		WriteTimeout   string `yaml:"write_timeout"`
		IdleTimeout    string `yaml:"idle_timeout"`
		MaxHeaderBytes int    `yaml:"max_header_bytes"`
		GinMode        string `yaml:"gin_mode"`
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
		TRACES struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"traces"`
		METRICS struct {
			Enabled bool   `yaml:"enabled"`
			Port    string `yaml:"port"`
		} `yaml:"metrics"`
		LOGS struct {
			EXPORT struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"export"`
		} `yaml:"logs"`
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
	}
	SECURITY struct {
		HMAC struct {
			Secret      string `yaml:"secret"`
			Algorithm   string `yaml:"algorithm"`
			Encoding    string `yaml:"encoding"`
			SkewSeconds int    `yaml:"skew_seconds"`
		} `yaml:"hmac"`
	} `yaml:"security"`
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
	PhotoSession struct {
		SlotDurationMinutes                int    `yaml:"slot_duration_minutes"`
		SlotsPerPeriod                     int    `yaml:"slots_per_period"`
		MorningStartHour                   int    `yaml:"morning_start_hour"`
		AfternoonStartHour                 int    `yaml:"afternoon_start_hour"`
		BusinessStartHour                  int    `yaml:"business_start_hour"`
		BusinessEndHour                    int    `yaml:"business_end_hour"`
		PhotographerHorizonMonths          int    `yaml:"photographer_horizon_months"`
		PhotographerTimezone               string `yaml:"photographer_timezone"`
		PhotographerAgendaRefreshIntervalH int    `yaml:"photographer_agenda_refresh_interval_hours"`
	} `yaml:"photo_session"`
	FCM struct {
		CredentialsFile string `yaml:"credentials_file"`
		ProjectID       string `yaml:"project_id"`
	}
	SystemUser struct {
		ResetPasswordURL string `yaml:"reset_password_url"`
	} `yaml:"system_user"`
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
	Schedule struct {
		DefaultBlockRules []struct {
			Start string `yaml:"start"`
			End   string `yaml:"end"`
		} `yaml:"default_block_rules"`
	} `yaml:"schedule"`
	Profiles map[string]ProfileOverrides `yaml:"profiles"`
	// Health endpoints are now integrated into the main HTTP server
	// No separate health configuration needed
}

type ProfileOverrides struct {
	HTTP struct {
		Port string `yaml:"port"`
	} `yaml:"http"`
	Telemetry TelemetryOverrides `yaml:"telemetry"`
	Workers   struct {
		Enabled *bool `yaml:"enabled"`
	} `yaml:"workers"`
}

type TelemetryOverrides struct {
	Enabled *bool `yaml:"enabled"`
	OTLP    *struct {
		Enabled  *bool  `yaml:"enabled"`
		Endpoint string `yaml:"endpoint"`
		Insecure *bool  `yaml:"insecure"`
	} `yaml:"otlp"`
	Traces *struct {
		Enabled *bool `yaml:"enabled"`
	} `yaml:"traces"`
	Metrics *struct {
		Enabled *bool  `yaml:"enabled"`
		Port    string `yaml:"port"`
	} `yaml:"metrics"`
	Logs *struct {
		Export *struct {
			Enabled *bool `yaml:"enabled"`
		} `yaml:"export"`
	} `yaml:"logs"`
}

type HMACSecurityConfig struct {
	Secret      string
	Algorithm   string
	Encoding    string
	SkewSeconds int
}

func (e *Environment) GetHMACSecurityConfig() HMACSecurityConfig {
	return HMACSecurityConfig{
		Secret:      e.SECURITY.HMAC.Secret,
		Algorithm:   e.SECURITY.HMAC.Algorithm,
		Encoding:    e.SECURITY.HMAC.Encoding,
		SkewSeconds: e.SECURITY.HMAC.SkewSeconds,
	}
}
