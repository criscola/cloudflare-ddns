package ddns

type Zone struct {
	ZoneId  string `mapstructure:"zone_id"`
	Proxied bool   `mapstructure:"proxied"`
	Ipv6    bool   `mapstructure:"ipv6"`
}

type Config struct {
	ApiToken        string `mapstructure:"api_token"`
	RefreshInterval string `mapstructure:"refresh_interval"`
	Zones           []Zone `mapstructure:"zones"`
}
