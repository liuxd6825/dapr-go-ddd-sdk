package restapp

type AuthConfig struct {
	StoreDbKey string `json:"storeDbKey"`
	Enable     *bool  `json:"enable"`
}

func (a *AuthConfig) IsEnabled() bool {
	return a.Enable != nil && *a.Enable
}

func (a *AuthConfig) GetStoreDbKey() string {
	return a.StoreDbKey
}
