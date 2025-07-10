package agent

import (
	pb "Calc_Service2/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"log"
	"time"
)

var grpcConn pb.OrchestratorClient

func InitGRPC() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	grpcConn = pb.NewOrchestratorClient(conn)
}

func getTask() (*pb.TaskRequest, error) {
	resp, err := grpcConn.GetTask(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	if !resp.HasTask {
		return nil, nil
	}
	return resp.Task, nil
}

func sendResult(taskID int64, result float64) error {
	_, err := grpcConn.SendResult(context.Background(), &pb.TaskResponse{
		Id:     taskID,
		Result: result,
	})
	return err
}

func compute(task *pb.TaskRequest) (float64, error) {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operator {
	case "+":
		return task.Operand1 + task.Operand2, nil
	case "-":
		return task.Operand1 - task.Operand2, nil
	case "*":
		return task.Operand1 * task.Operand2, nil
	case "/":
		if task.Operand2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Operand1 / task.Operand2, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", task.Operator)
	}
}

func worker() {
	for {
		task, err := getTask()
		if err != nil || task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		result, err := compute(task)
		if err != nil {
			log.Printf("Compute error: %v", err)
			continue
		}

		if err := sendResult(task.Id, result); err != nil {
			log.Printf("Send result failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func Agents() {
	InitGRPC()
	go worker()
	select {}
}
