package matrix

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"rc-cli/filesystem"
)

type SynapseDef struct {
	ServerName string `yaml:"server_name"`
	PidFile    string `yaml:"pid_file"`
	Listeners  []struct {
		Port       int    `yaml:"port"`
		TLS        bool   `yaml:"tls"`
		Type       string `yaml:"type"`
		XForwarded bool   `yaml:"x_forwarded"`
		Resources  []struct {
			Names    []string `yaml:"names"`
			Compress bool     `yaml:"compress"`
		} `yaml:"resources"`
	} `yaml:"listeners"`
	LogConfig                string `yaml:"log_config"`
	MediaStorePath           string `yaml:"media_store_path"`
	RegistrationSharedSecret string `yaml:"registration_shared_secret"`
	ReportStats              bool   `yaml:"report_stats"`
	MacaroonSecretKey        string `yaml:"macaroon_secret_key"`
	FormSecret               string `yaml:"form_secret"`
	SigningKeyPath           string `yaml:"signing_key_path"`
	TrustedKeyServers        []struct {
		ServerName string `yaml:"server_name"`
	} `yaml:"trusted_key_servers"`
	AppServiceConfigFiles []string `yaml:"app_service_config_files"`
	Retention             struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"retention"`
	EnableRegistration                    bool     `yaml:"enable_registration"`
	EnableRegistrationWithoutVerification bool     `yaml:"enable_registration_without_verification"`
	SuppressKeyServerWarning              bool     `yaml:"suppress_key_server_warning"`
	FederationIPRangeBlacklist            []string `yaml:"federation_ip_range_blacklist"`
	Database                              struct {
		Name string `yaml:"name"`
		Args struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			Host     string `yaml:"host"`
			CpMin    int    `yaml:"cp_min"`
			CpMax    int    `yaml:"cp_max"`
		} `yaml:"args"`
	} `yaml:"database"`
	Redis struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
	} `yaml:"redis"`
	RcMessage struct {
		PerSecond  int `yaml:"per_second"`
		BurstCount int `yaml:"burst_count"`
	} `yaml:"rc_message"`
	RcRegistration struct {
		PerSecond  int `yaml:"per_second"`
		BurstCount int `yaml:"burst_count"`
	} `yaml:"rc_registration"`
	RcRegistrationTokenValidity struct {
		PerSecond  int `yaml:"per_second"`
		BurstCount int `yaml:"burst_count"`
	} `yaml:"rc_registration_token_validity"`
	RcLogin struct {
		Address struct {
			PerSecond  float64 `yaml:"per_second"`
			BurstCount int     `yaml:"burst_count"`
		} `yaml:"address"`
		Account struct {
			PerSecond  float64 `yaml:"per_second"`
			BurstCount int     `yaml:"burst_count"`
		} `yaml:"account"`
		FailedAttempts struct {
			PerSecond  float64 `yaml:"per_second"`
			BurstCount int     `yaml:"burst_count"`
		} `yaml:"failed_attempts"`
	} `yaml:"rc_login"`
	RcAdminRedaction struct {
		PerSecond  int `yaml:"per_second"`
		BurstCount int `yaml:"burst_count"`
	} `yaml:"rc_admin_redaction"`
	RcJoins struct {
		Local struct {
			PerSecond  int `yaml:"per_second"`
			BurstCount int `yaml:"burst_count"`
		} `yaml:"local"`
		Remote struct {
			PerSecond  int `yaml:"per_second"`
			BurstCount int `yaml:"burst_count"`
		} `yaml:"remote"`
	} `yaml:"rc_joins"`
	Rc3PidValidation struct {
		PerSecond  float64 `yaml:"per_second"`
		BurstCount int     `yaml:"burst_count"`
	} `yaml:"rc_3pid_validation"`
	RcInvites struct {
		PerRoom struct {
			PerSecond  int `yaml:"per_second"`
			BurstCount int `yaml:"burst_count"`
		} `yaml:"per_room"`
		PerUser struct {
			PerSecond  int `yaml:"per_second"`
			BurstCount int `yaml:"burst_count"`
		} `yaml:"per_user"`
	} `yaml:"rc_invites"`
	RcFederation struct {
		WindowSize  int `yaml:"window_size"`
		SleepLimit  int `yaml:"sleep_limit"`
		SleepDelay  int `yaml:"sleep_delay"`
		RejectLimit int `yaml:"reject_limit"`
		Concurrent  int `yaml:"concurrent"`
	} `yaml:"rc_federation"`
	FederationRrTransactionsPerRoomPerSecond int `yaml:"federation_rr_transactions_per_room_per_second"`
}

var Synapse SynapseDef

func readSynapseFile(basePath string) {
	_, filePath := getHomeserverFilePath(SynapseType, basePath)

	// Unmarshal defaults
	err := yaml.Unmarshal([]byte(homeserverSynapseDefaults), &Synapse)
	if err != nil {
		log.Fatal("[ReadSynapseFile] Could not unmarshal defaults: ", err)
	}

	// Load current file if exists
	if _, err := os.Stat(filePath); err == nil {
		currentFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("[ReadSynapseFile] Could not read file: ", err)
		}

		// Unmarshal current file
		err = yaml.Unmarshal(currentFile, &Synapse)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeSynapseFile(basePath string) {
	dirPath, filePath := getHomeserverFilePath(SynapseType, basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[WriteSynapseFile] Could not create directories: ", err)
	}

	// Marshal struct
	result, err := yaml.Marshal(Synapse)
	if err != nil {
		log.Fatal("[WriteSynapseFile] Could not marshal struct: ", err)
	}

	// Write file
	err = os.WriteFile(filePath, result, 0644)
	if err != nil {
		log.Fatal("[WriteSynapseFile] Could not write file: ", err)
	}
}
