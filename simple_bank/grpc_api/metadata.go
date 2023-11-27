package grpcapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	gatewayUserAgentHeaderKey     = "grpcgateway-user-agent"
	grpcAgentHeaderKey            = "user-agent"
	gatewayXForwardedForHeaderKey = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(gatewayUserAgentHeaderKey); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(grpcAgentHeaderKey); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		if clientIPs := md.Get(gatewayXForwardedForHeaderKey); len(clientIPs) > 0 {
			meta.ClientIp = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIp = p.Addr.String()
	}
	return meta
}
