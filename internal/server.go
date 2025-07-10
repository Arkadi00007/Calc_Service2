package internal

import (
	"Calc_Service2/pkg/calculation"
	pb "Calc_Service2/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"strconv"
)

type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
}

func (s *OrchestratorServer) GetTask(ctx context.Context, _ *emptypb.Empty) (*pb.TaskReply, error) {

	for _, expr := range expressions {
		if expr.Status == "processing" && expr.SubStatus != "waiting" {
			expr.SubStatus = "waiting"
			task := createTasks(expr.ID, expr.Expression)

			if task == nil {
				return &pb.TaskReply{HasTask: false}, nil
			}

			return &pb.TaskReply{
				HasTask: true,
				Task: &pb.TaskRequest{
					Id:            int64(task.ID),
					Operand1:      task.Arg1,
					Operand2:      task.Arg2,
					Operator:      task.Operation,
					OperationTime: int32(task.Operation_time),
				},
			}, nil
		}
	}
	return &pb.TaskReply{HasTask: false}, nil

}
func (s *OrchestratorServer) SendResult(ctx context.Context, res *pb.TaskResponse) (*emptypb.Empty, error) {
	expr, exists := expressions[int(res.Id)]
	if !exists {
		return &emptypb.Empty{}, nil
	}

	// Вставка результата
	for i := 2; i < len(*expr.Expression); i++ {
		if calculation.IsOperator((*expr.Expression)[i][0]) {
			*expr.Expression = append((*expr.Expression)[:i-2], append([]string{strconv.FormatFloat(res.Result, 'f', -1, 64)}, (*expr.Expression)[i+1:]...)...)
			expr.SubStatus = ""
			break
		}
	}

	if len(*expr.Expression) == 1 {
		expr.Status = "completed"
		expr.Result = res.Result
		db.Exec("UPDATE expressions SET status=?, result=? WHERE id=?", "completed", expr.Result, expr.ID)
	}
	return &emptypb.Empty{}, nil
}

func runGRPC() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, &OrchestratorServer{})
	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//service Orchestrator {
//rpc GetTask (google.protobuf.Empty) returns (TaskReply);
//rpc SendResult (TaskResponse) returns (google.protobuf.Empty);
//}
