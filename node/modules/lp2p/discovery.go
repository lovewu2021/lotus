package lp2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"go.uber.org/fx"

	"github.com/filecoin-project/go-lotus/node/modules/helpers"
)

const discoveryConnTimeout = time.Second * 30

type discoveryHandler struct {
	ctx  context.Context
	host host.Host
}

func (dh *discoveryHandler) HandlePeerFound(p peer.AddrInfo) {
	log.Warnw("discovred peer", "peer", p)
	ctx, cancel := context.WithTimeout(dh.ctx, discoveryConnTimeout)
	defer cancel()
	if err := dh.host.Connect(ctx, p); err != nil {
		log.Warnw("failed to connect to peer found by discovery", "error", err)
	}
}

func DiscoveryHandler(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host) *discoveryHandler {
	return &discoveryHandler{
		ctx:  helpers.LifecycleCtx(mctx, lc),
		host: host,
	}
}

func SetupDiscovery(mdns bool, mdnsInterval int) func(helpers.MetricsCtx, fx.Lifecycle, host.Host, *discoveryHandler) error {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host, handler *discoveryHandler) error {
		if mdns {
			if mdnsInterval == 0 {
				mdnsInterval = 5
			}
			service, err := discovery.NewMdnsService(helpers.LifecycleCtx(mctx, lc), host, time.Duration(mdnsInterval)*time.Second, discovery.ServiceTag)
			if err != nil {
				log.Errorw("mdns error", "error", err)
				return nil
			}
			service.RegisterNotifee(handler)
		}
		return nil
	}
}
