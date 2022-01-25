package ddns

import (
	"context"
	cf "github.com/cloudflare/cloudflare-go"
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
	*cf.API
	Logger *zap.SugaredLogger
	Config Config
}

// Start starts the main application loop.
func (client Client) Start() {
	logger = client.Logger
	ctx := context.Background()

	// Main loop
	for {
		// TODO: Check and update zones concurrently
		// Divide zones array per number of cores available
		// Spawn goroutine to handle each zone concurrently
		// Return done in loop when finished, then wait refresh interval before cycling

		publicIp, err := getPublicIp()
		if err != nil {
			logger.Errorf("cannot get public IP: %s", err)
			time.Sleep(client.Config.RefreshInterval)
			continue
		}
		logger.Debugf("retrieved public IP: %s", publicIp)

		// For each zone in config
		for _, zoneCfg := range client.Config.Zones {
			// Get zone id, zone name and attach them to the logging trace
			logger = client.Logger.With("zoneId", zoneCfg.ZoneId)
			logger.Info("getting zone details...")
			zoneDetails, err := client.ZoneDetails(ctx, zoneCfg.ZoneId)
			if err != nil {
				logger.Errorf("cannot get zone details: %s", err)
				continue
			}
			logger = logger.With("zoneName", zoneDetails.Name)
			logger.Debug("retrieved zone information successfully")

			// If there is any explicit record set to update, do it
			if zoneCfg.ContainsExplicitRecords() {
				logger.Debug("zoneCfg contains explicit records")
				// For each explicit record in config
				for _, recordCfg := range zoneCfg.ExplicitRecords {
					// TODO: Maybe extract function?
					subLogger := logger.With("recordName", recordCfg.Name)
					subLogger.Info("updating explicit record")
					explicitIpv4Record, err := client.getIpv4Record(ctx, zoneCfg.ZoneId, recordCfg.Name+"."+zoneDetails.Name)
					if err != nil {
						subLogger.Errorf("cannot get explicit record: %s", err)
						continue
					}
					// Update explicit record with public IP
					subLogger.Debug("retrieved explicit record successfully")
					if isRecordIpUpToDate(explicitIpv4Record, publicIp) {
						subLogger.Info("record is already up to date with current public IP: skipping update...")
						continue
					}
					oldIp := explicitIpv4Record.Content
					explicitIpv4Record.Content = publicIp
					err = client.updateRecord(ctx, explicitIpv4Record, publicIp, zoneCfg.ZoneId)
					if err != nil {
						subLogger.Errorf("cannot update explicit record: %s", err)
						continue
					}
					subLogger.Infof("explicit record updated successfully, old IP %s was replaced with new IP %s", oldIp, publicIp)
				}
				continue
			}
			// Else, update root record from zone name
			ipv4RootRecord, err := client.getIpv4Record(ctx, zoneCfg.ZoneId, zoneDetails.Name)
			if err != nil {
				logger.Errorf("cannot get IPv4 root record: %s", err)
				continue
			}
			logger.Debug("retrieved IPv4 root record successfully")
			if isRecordIpUpToDate(ipv4RootRecord, publicIp) {
				logger.Info("record is already up to date with current public IP: skipping update...")
				continue
			}
			// Else, update
			oldIp := ipv4RootRecord.Content
			err = client.updateRecord(ctx, ipv4RootRecord, publicIp, zoneCfg.ZoneId)
			if err != nil {
				logger.Errorf("cannot update IPv4 root record: %s", err)
				continue
			}
			logger.Infof("IPv4 root record updated successfully, old IP %s was replaced with new IP %s", oldIp, publicIp)
		}
		time.Sleep(client.Config.RefreshInterval)
	}
}

func (client Client) updateRecord(ctx context.Context, record cf.DNSRecord, publicIp, zoneId string) error {
	record.Content = publicIp
	return client.UpdateDNSRecord(ctx, zoneId, record.ID, record)
}

func (client Client) getIpv4Record(ctx context.Context, zoneId, name string) (cf.DNSRecord, error) {
	return client.getRecordFromZone(ctx, zoneId, func(record cf.DNSRecord) bool {
		if record.Type == "A" && record.Name == name {
			return true
		}
		return false
	})
}

func (client Client) getIpv6Record(ctx context.Context, zoneId, name string) (cf.DNSRecord, error) {
	return client.getRecordFromZone(ctx, zoneId, func(record cf.DNSRecord) bool {
		if record.Type == "AAAA" && record.Name == name {
			return true
		}
		return false
	})
}

func (client Client) getRecordFromZone(ctx context.Context, zoneId string, filter func(cf.DNSRecord) bool) (cf.DNSRecord, error) {
	dnsRecords, err := client.DNSRecords(ctx, zoneId, cf.DNSRecord{})
	if err != nil {
		return cf.DNSRecord{}, err
	}
	for _, record := range dnsRecords {
		if filter(record) {
			return record, nil
		}
	}
	return cf.DNSRecord{}, nil
}

// getPublicIp gets the public IP by querying 1.1.1.1 and parsing its response.
func getPublicIp() (string, error) {
	response, err := http.Get(PublicIpPage)
	if err != nil {
		//logger.Fatalf("cannot get public IP: %s", err)
		return "", err
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//logger.Fatalf("cannot read body: %s", err)
		return "", err
	}
	payload := string(resp)
	lines := strings.Split(payload, "\n")
	for _, line := range lines {
		if strings.Contains(line, "ip") {
			return strings.Split(line, "=")[1], nil
		}
	}
	return "", nil
}

func isRecordIpUpToDate(record cf.DNSRecord, publicIp string) bool {
	return record.Content == publicIp
}
