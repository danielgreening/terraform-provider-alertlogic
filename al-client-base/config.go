package al_client_base

type Config struct {
	AccessKey           string
	YarpEndpoint        string
	CredsFilename       string
	DebugLogging        bool
	Insecure            bool
	MaxRetries          int
	Profile             string
	SecretKey           string
	SkipCredsValidation bool
	Token               string
	UserAgentProducts   []*UserAgentProduct
}

type UserAgentProduct struct {
	Extra   []string
	Name    string
	Version string
}
