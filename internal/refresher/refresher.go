package refresher

import (
	"context"
	"time"

	"github.com/thomas-lepage/domain_exporter/internal/client"
	"github.com/thomas-lepage/domain_exporter/internal/safeconfig"
	"github.com/rs/zerolog/log"
)

type Refresher struct {
	ticker  *time.Ticker
	client  client.Client
	domains []safeconfig.Domain
	timeout time.Duration
}

func New(interval time.Duration, client client.Client, timeout time.Duration, domains ...safeconfig.Domain) Refresher {
	ticker := time.NewTicker(interval)
	return Refresher{
		ticker:  ticker,
		client:  client,
		domains: domains,
		timeout: timeout,
	}
}

func (r Refresher) Stop() {
	r.ticker.Stop()
}

func (r Refresher) Run(ctx context.Context) {
	log.Info().Msg("run refresher")
	r.Refresh(ctx)

	select {
	case <-r.ticker.C:
		r.Refresh(ctx)
	case <-ctx.Done():
		log.Info().Msg("refresher is finished")
		return
	}
}

func (r Refresher) Refresh(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	for _, domain := range r.domains {
		if _, err := r.client.ExpireTime(ctx, domain.Name, domain.Host); err != nil {
			log.Error().Err(err).Msgf("failed to get expire time for %s", domain)
		}
	}
	log.Debug().Msg("refresh is done")
}
