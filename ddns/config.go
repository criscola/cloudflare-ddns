package ddns

import "time"

type Record struct {
	Name string `mapstructure:"name"`
}

type Zone struct {
	ZoneId          string   `mapstructure:"zone_id"`
	Proxied         bool     `mapstructure:"proxied"`
	Ipv6            bool     `mapstructure:"ipv6"`
	ExplicitRecords []Record `mapstructure:"explicit_records"`
}

type Config struct {
	ApiToken        string        `mapstructure:"api_token"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
	Zones           []Zone        `mapstructure:"zones"`
}

func (zone Zone) ContainsExplicitRecords() bool {
	return len(zone.ExplicitRecords) > 0
}
