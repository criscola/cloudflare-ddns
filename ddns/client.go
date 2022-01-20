package ddns

import (
	"github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

// Start starts the main application
func Start(cfClient *cloudflare.API, logger *zap.Logger, config Config) {
	// TODO: Main loop, get public IP using 1.1.1.1, read each zone and update A record for root domain
}
