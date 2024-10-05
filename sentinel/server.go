package sentinel

import (
	context "context"
	grpc "google.golang.org/grpc"
	"log/slog"
)

type sentinelServer struct {
	UnimplementedReplicatedLogSentinelServer
}

func (s sentinelServer) Replicate(_ context.Context, in *ReplicateRequest) (*ReplicateResponse, error) {
	isDuplicate := syncReplica(int(in.Id), in.Message, in.Timestamp)

	if isDuplicate {
		slog.Warn("Duplicate item replication,", "message", in)
	} else {
		slog.Info("Replicated item successfully,", "message", in)
	}

	return &ReplicateResponse{Ack: true}, nil
}

func InitServer(s *grpc.Server) {
	RegisterReplicatedLogSentinelServer(s, &sentinelServer{})
}
