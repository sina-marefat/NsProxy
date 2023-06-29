package dns

import (
	"encoding/json"
	"nsproxy/cache"
	"time"
)

type DNSRepository struct {
	cache      cache.Cache
	defaultTTL time.Duration
}

func (dr *DNSRepository) SetDNSInCache(domain string, IPs []string) error {
	marshaledIps, err := json.Marshal(IPs)
	if err != nil {
		return err
	}
	return dr.cache.Set(domain, marshaledIps, dr.defaultTTL)
}

func (dr *DNSRepository) GetDNSFromCache(domain string) ([]string, error) {
	marshalledIps, err := dr.cache.Get(domain)
	if err != nil {
		return nil, err
	}
	var IPs []string

	err = json.Unmarshal([]byte(marshalledIps.(string)), &IPs)

	return IPs, err
}

func NewDnsRepo(cache cache.Cache) *DNSRepository {
	return &DNSRepository{cache: cache}
}
