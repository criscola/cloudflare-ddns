package ddns

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const PublicIpPage = "https://1.1.1.1/cdn-cgi/trace"

var logger *zap.SugaredLogger

// Client is a wrapper around cloudflare.API.
type Client struct {
	*cloudflare.API
	Logger *zap.SugaredLogger
	Config Config
}

// Start starts the main application loop.
func (client Client) Start() {
	logger = client.Logger
	ctx := context.Background()

	// Main loop
	for {
		publicIp := getPublicIp()
		logger.Debugf("retrieved public IP: %s", publicIp)

		// TODO: Check and update zones concurrently
		for _, zoneCfg := range client.Config.Zones {
			logger = client.Logger.With("zoneId", zoneCfg.ZoneId)
			logger.Info("checking A record...")
			zoneDetails, err := client.ZoneDetails(ctx, zoneCfg.ZoneId)
			if err != nil {
				logger.Errorf("cannot get zone details: %s", err)
				continue
			}
			logger = logger.With("zoneName", zoneDetails.Name)
			logger.Debug("retrieved zone information successfully")
			ipv4RootRecord, err := client.getIpv4RootRecord(ctx, zoneCfg.ZoneId, zoneDetails.Name)
			if err != nil {
				logger.Errorf("cannot get IPv4 root record: %s", err)
				continue
			}
			logger.Debug("retrieved IPv4 root record successfully")
			oldIP := ipv4RootRecord.Content
			if oldIP == publicIp {
				logger.Info("record contains the same current public IP: skipping update...")
				continue
			}
			ipv4RootRecord.Content = publicIp
			err = client.UpdateDNSRecord(ctx, zoneCfg.ZoneId, ipv4RootRecord.ID, ipv4RootRecord)
			if err != nil {
				logger.Errorf("cannot update IPv4 root record: %s", err)
				continue
			}
			logger.Infof("IPv4 root record updated successfully, old IP %s was replaced with new IP %s", oldIP, publicIp)
		}
		refreshInterval, err := time.ParseDuration(client.Config.RefreshInterval)
		if err != nil {
			client.Logger.Fatal("cannot parse refresh interval, please specify an interval like this: 500ms, 10s, 5m etc.")
		}
		time.Sleep(refreshInterval)
	}
}

func (client Client) getIpv4RootRecord(ctx context.Context, zoneId string, domainName string) (cloudflare.DNSRecord, error) {
	dnsRecords, err := client.DNSRecords(ctx, zoneId, cloudflare.DNSRecord{})
	if err != nil {
		return cloudflare.DNSRecord{}, err
	}
	for _, record := range dnsRecords {
		if record.Type == "A" && record.Name == domainName {
			return record, nil
		}
	}
	return cloudflare.DNSRecord{}, nil
}

func (client Client) getIpv6RootRecord(ctx context.Context, zoneId string) (cloudflare.DNSRecord, error) {
	// TODO: ipv6 support
	return cloudflare.DNSRecord{}, nil
}

// getPublicIp gets the public IP by querying 1.1.1.1 and parsing its response.
func getPublicIp() string {
	response, err := http.Get(PublicIpPage)
	if err != nil {
		logger.Fatalf("cannot get public IP: %s", err)
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Fatalf("cannot read body: %s", err)
	}
	payload := string(resp)
	lines := strings.Split(payload, "\n")
	for _, line := range lines {
		if strings.Contains(line, "ip") {
			return strings.Split(line, "=")[1]
		}
	}
	return ""
}
