package config

type KV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type RequestService struct {
	DatabaseUrl     string                 `yaml:"database_url"`
	DatabaseType    string                 `yaml:"database_type"`
	ApiType         string                 `yaml:"api_type"`
	RestServiceHost string                 `yaml:"rest_service_host"`
	RestServicePort string                 `yaml:"rest_service_port"`
	WsServiceHost   string                 `yaml:"ws_service_host"`
	WsServicePort   string                 `yaml:"ws_service_port"`
	TestMode        string                 `yaml:"test_mode"`
	Games           []KV                   `yaml:"games"`
	Ext             map[string]interface{} `yaml:"ext"`
}

type OAuthProviderConfig struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	UserInfoUrl  string   `yaml:"user_info_url"`
	RedirectUrl  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`
}

type UserService struct {
	UserServiceHost  string              `yaml:"user_service_host"`
	UserServicePort  string              `yaml:"user_service_port"`
	AllowedHosts     []string            `yaml:"allowed_hosts"`
	GoogleProvider   OAuthProviderConfig `yaml:"google_provider"`
	FacebookProvider OAuthProviderConfig `yaml:"facebook_provider"`
	DBDriver         string              `yaml:"db_driver"`
	DBName           string              `yaml:"db_name"`
}
