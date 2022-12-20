package element

import (
	"encoding/json"
	"log"
	"os"
	"rc-cli/filesystem"
)

type ConfigDef struct {
	DefaultServerConfig struct {
		Homeserver struct {
			BaseUrl    string `json:"base_url"`
			ServerName string `json:"server_name"`
		} `json:"m.homeserver"`
		IdentityServer struct {
			BaseUrl string `json:"base_url"`
		} `json:"m.identity_server"`
	} `json:"default_server_config"`
	Brand                   string   `json:"brand"`
	IntegrationsUiUrl       string   `json:"integrations_ui_url"`
	IntegrationsRestUrl     string   `json:"integrations_rest_url"`
	IntegrationsWidgetsUrls []string `json:"integrations_widgets_urls"`
	HostingSignupLink       string   `json:"hosting_signup_link"`
	BugReportEndpointUrl    string   `json:"bug_report_endpoint_url"`
	UISIAutorageshakeApp    string   `json:"uisi_autorageshake_app"`
	ShowLabsSettings        bool     `json:"showLabsSettings"`
	Piwik                   struct {
		Url       string `json:"url"`
		SiteId    int    `json:"siteId"`
		PolicyUrl string `json:"policyUrl"`
	} `json:"piwik"`
	RoomDirectory struct {
		Servers []string `json:"servers"`
	} `json:"roomDirectory"`
	EnablePresenceByHsUrl struct {
		HttpsMatrixOrg             bool `json:"https://matrix.org"`
		HttpsMatrixClientMatrixOrg bool `json:"https://matrix-client.matrix.org"`
	} `json:"enable_presence_by_hs_url"`
	TermsAndConditionsLinks []struct {
		Url  string `json:"url"`
		Text string `json:"text"`
	} `json:"terms_and_conditions_links"`
	HostSignup struct {
		Brand             string   `json:"brand"`
		CookiePolicyUrl   string   `json:"cookiePolicyUrl"`
		Domains           []string `json:"domains"`
		PrivacyPolicyUrl  string   `json:"privacyPolicyUrl"`
		TermsOfServiceUrl string   `json:"termsOfServiceUrl"`
		Url               string   `json:"url"`
	} `json:"hostSignup"`
	Sentry struct {
		Dsn         string `json:"dsn"`
		Environment string `json:"environment"`
	} `json:"sentry"`
	Posthog struct {
		ProjectApiKey string `json:"projectApiKey"`
		ApiHost       string `json:"apiHost"`
	} `json:"posthog"`
	Features struct {
		FeatureSpotlight bool `json:"feature_spotlight"`
	} `json:"features"`
	MapStyleUrl string `json:"map_style_url"`
}

var Config ConfigDef

func ReadConfigFile(basePath string) {
	_, filePath := getConfigFilePath(basePath)

	// Unmarshal defaults
	err := json.Unmarshal([]byte(elementDefaults), &Config)
	if err != nil {
		log.Fatal("[Element][ReadConfigFile] Could not unmarshal defaults: ", err)
	}

	// Load current file if exists
	if _, err := os.Stat(filePath); err == nil {
		currentFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("[Element][ReadConfigFile] Could not read file: ", err)
		}

		// Unmarshal current file
		err = json.Unmarshal(currentFile, &Config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteConfigFile(basePath string) {
	dirPath, filePath := getConfigFilePath(basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[Element][WriteConfigFile] Could not create directories: ", err)
	}

	// Marshal struct
	result, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		log.Fatal("[Element][WriteConfigFile] Could not marshal struct: ", err)
	}

	// Write file
	err = os.WriteFile(filePath, result, 0644)
	if err != nil {
		log.Fatal("[Element][WriteConfigFile] Could not write file: ", err)
	}
}
