package dns

import (
	"nsproxy/cache"
	"time"
)

type DNSRepository struct {
	cache      cache.Cache
	defaultTTL time.Duration
}

func (dr *DNSRepository) SetDNSInCache(domain string, response []byte, count []byte) error {
	return dr.cache.Set(domain, append(count, response...), dr.defaultTTL)
}

func (dr *DNSRepository) GetDNSFromCache(domain string) ([]byte, []byte, error) {
	marshalledIps, err := dr.cache.Get(domain)
	if err != nil {
		return nil, nil, err
	}

	response := []byte(marshalledIps.(string))

	return response[2:], response[0:2], err
}

func NewDnsRepo(cache cache.Cache) *DNSRepository {
	return &DNSRepository{cache: cache}
}
