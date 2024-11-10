package sentinel

import (
	context "context"
	jitter "github.com/go-toolbelt/jitter"
	grpc "google.golang.org/grpc"
	"log/slog"
	"time"
)

type sentinelServer struct {
	UnimplementedReplicatedLogSentinelServer
}

// if item number 3, 6, 9 etc - longer delay (up to 3 seconds)
var DELAYS = [3]int{1, 1, 3}

func (s sentinelServer) Replicate(_ context.Context, in *ReplicateRequest) (*ReplicateResponse, error) {

	delaySecondsBasedOnId := DELAYS[int(in.Id-1)%len(DELAYS)]

	initialDelay := time.Duration(delaySecondsBasedOnId) * time.Second

	delay := jitter.By(initialDelay, 500*time.Millisecond)
	////random delay 1-5 seconds
	//sleepSeconds := rand.Int()%5 + 1
	//
	//// if item number 3, 6, 9 etc - longer delay (6-10 seconds)
	//if (in.Id+1)%3 == 0 {
	//	sleepSeconds += 5
	//}

	time.Sleep(delay)

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
