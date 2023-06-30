package dns

import (
	"nsproxy/cache"
	"nsproxy/config"
	"time"
)

type DNSRepository struct {
	cache      cache.Cache
	defaultTTL time.Duration
}

func (dr *DNSRepository) SetDNSInCache(domain string, response []byte) error {
	return dr.cache.Set(domain, response, dr.defaultTTL)
}

func (dr *DNSRepository) GetDNSFromCache(domain string) ([]byte, error) {
	marshalledIps, err := dr.cache.Get(domain)
	if err != nil {
		return nil, err
	}

	response := []byte(marshalledIps.(string))

	return response, err
}

func NewDnsRepo(cache cache.Cache) *DNSRepository {
	return &DNSRepository{cache: cache,defaultTTL: config.GetConfig().CacheExpireTTL}
}
