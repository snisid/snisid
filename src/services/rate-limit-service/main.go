package main

import (
	"context"
	"log"
	"net"

	pbrl "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"google.golang.org/grpc"
)

type rateLimitServer struct {
	pbrl.UnimplementedRateLimitServiceServer
}

func (s *rateLimitServer) ShouldRateLimit(ctx context.Context, req *pbrl.RateLimitRequest) (*pbrl.RateLimitResponse, error) {
	log.Printf("Received rate limit request from domain: %s", req.Domain)

	// Default response: OK (allow traffic)
	resp := &pbrl.RateLimitResponse{
		OverallCode: pbrl.RateLimitResponse_OK,
	}

	// Loop through the descriptors provided by the API Gateway
	for _, descriptor := range req.Descriptors {
		for _, entry := range descriptor.Entries {
			// Implement Custom Fraud Logic Here
			// Example: Block a specific IP acting suspiciously
			if entry.Key == "remote_address" && entry.Value == "203.0.113.50" {
				log.Printf("SECURITY ALERT: Blocked suspicious IP: %s", entry.Value)
				resp.OverallCode = pbrl.RateLimitResponse_OVER_LIMIT
				return resp, nil
			}

			// Example: Throttle excessive /v1/identity requests (Simulated Redis Check)
			if entry.Key == "generic_key" && entry.Value == "identity_api_abuse" {
				log.Printf("FRAUD ALERT: API abuse detected for path: /v1/identity")
				resp.OverallCode = pbrl.RateLimitResponse_OVER_LIMIT
				return resp, nil
			}
		}
	}

	return resp, nil
}

func main() {
	port := ":8081"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pbrl.RegisterRateLimitServiceServer(grpcServer, &rateLimitServer{})

	log.Printf("SNISID Rate Limit Service starting on port %s...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
