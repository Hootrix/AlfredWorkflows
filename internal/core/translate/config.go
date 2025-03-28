package translate

// ConfigItem 定义单个服务配置项
type ConfigItem struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url,omitempty"`
	Token     string `yaml:"token,omitempty"`
	AppKey    string `yaml:"app_key,omitempty"`
	AppSecret string `yaml:"app_secret,omitempty"`
}

// Config 定义整体配置结构体
type Config struct {
	Services []ConfigItem `yaml:"services"`
	Timeout  int          `yaml:"timeout"`
}

// GetConfigItemWithName 根据名称获取配置项
func (c *Config) GetConfigItemWithName(name string) *ConfigItem {
	for _, item := range c.Services {
		if item.Name == name {
			return &item
		}
	}
	return nil
}
