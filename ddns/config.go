package ddns

type Zone struct {
	ZoneId  string `yaml:"zone_id"`
	Proxied bool   `yaml:"proxied"`
	Ipv6    bool   `yaml:"ipv6"`
}

type Config struct {
	ApiToken        string `yaml:"api_token"`
	RefreshInterval int    `yaml:"refresh_interval"`
	Zones           []Zone `yaml:"zones"`
}
