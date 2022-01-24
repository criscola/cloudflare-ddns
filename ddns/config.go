package ddns

import "time"

type Zone struct {
	ZoneId  string `mapstructure:"zone_id"`
	Proxied bool   `mapstructure:"proxied"`
	Ipv6    bool   `mapstructure:"ipv6"`
}

type Config struct {
	ApiToken        string        `mapstructure:"api_token"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
	Zones           []Zone        `mapstructure:"zones"`
}
