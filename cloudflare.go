package main

import (
	"context"
	"myip/logger"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"go.uber.org/zap"
)

var (
	zone         string
	domain       string
	client       *cloudflare.API               = nil
	zoneResource *cloudflare.ResourceContainer = nil
)

func init() {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	zone = os.Getenv("CLOUDFLARE_ZONE_ID")
	domain = os.Getenv("DDNS_DOMAIN")
	if len(token) == 0 || len(zone) == 0 || len(domain) == 0 {
		return
	}
	var err error
	client, err = cloudflare.NewWithAPIToken(token)
	if err != nil {
		logger.Panic("cloudflare.NewWithAPIToken", zap.Error(err))
	}
	zoneResource = cloudflare.ZoneIdentifier(zone)
}

func GetRecord(ctx context.Context) (cloudflare.DNSRecord, error) {

	rs, _, err := client.ListDNSRecords(ctx, zoneResource, cloudflare.ListDNSRecordsParams{
		Name: domain,
		Type: "A",
	})
	if err != nil {
		logger.Error("client.ListDNSRecords", zap.Error(err))
		return cloudflare.DNSRecord{}, err
	}
	if len(rs) == 0 {
		logger.Error("client.ListDNSRecords", zap.String("zone", zone), zap.String("domain", domain))
		return cloudflare.DNSRecord{}, nil
	}
	return rs[0], nil
}

func PutRecord(ctx context.Context, ip string) error {

	var (
		err          error
		record       cloudflare.DNSRecord
		createParams cloudflare.CreateDNSRecordParams = cloudflare.CreateDNSRecordParams{
			Type:    "A",
			Name:    domain,
			Content: ip,
			TTL:     1,
		}
		updateParams cloudflare.UpdateDNSRecordParams = cloudflare.UpdateDNSRecordParams{
			Type:    "A",
			Name:    domain,
			Content: ip,
		}
	)

	r, err := GetRecord(ctx)
	if err != nil {
		logger.Error("getRecord", zap.Error(err))
		return err
	}

	if r.ID == "" {
		record, err = client.CreateDNSRecord(ctx, zoneResource, createParams)
	} else {
		updateParams.ID = r.ID
		record, err = client.UpdateDNSRecord(ctx, zoneResource, updateParams)
	}
	if err != nil {
		logger.Error("client.CreateDNSRecord", zap.Error(err), zap.Any("record", r), zap.Any("createParams", createParams), zap.Any("updateParams", updateParams))
		return err
	}
	logger.Info("PutRecord", zap.Any("record", record))
	return nil
}
