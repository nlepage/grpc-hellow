package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/nlepage/grpc-hellow/go/hellow/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	hellowCmd = &cobra.Command{
		Use:   "hellow",
		Short: "hellow is a gRPC hello world",
	}
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "invokes the hellow client",
		RunE:  sayHellow,
		Args:  cobra.ExactArgs(2),
	}
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "starts the hellow server",
		RunE:  serve,
		Args:  cobra.NoArgs,
	}
	host     string
	port     = 9090
	isStream bool
)

func init() {
	clientCmd.Flags().StringVar(&host, "host", "localhost", "Host to call")
	clientCmd.Flags().BoolVarP(&isStream, "stream", "s", false, "Ask for a streamed response")
	hellowCmd.AddCommand(clientCmd)
	hellowCmd.AddCommand(serverCmd)
}

func main() {
	hellowCmd.Execute()
}

func sayHellow(cmd *cobra.Command, args []string) error {
	name := args[0]
	count, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := service.NewHellowClient(conn)

	if isStream {
		stream, err := client.StreamHellow(context.Background(), &service.SayHellowRequest{
			Name:  name,
			Count: int64(count),
		})
		if err != nil {
			return err
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Println(res.Message)
		}
	} else {
		res, err := client.SayHellow(context.Background(), &service.SayHellowRequest{
			Name:  name,
			Count: int64(count),
		})
		if err != nil {
			return err
		}
		fmt.Println(res.Message)
	}

	return nil
}

func serve(cmd *cobra.Command, args []string) error {
	return (&server{}).serve()
}

type server struct{}

func (s *server) serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	service.RegisterHellowServer(grpcSrv, s)
	if err := grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("Failed to serve: %v", err)
	}

	return nil
}

func (s *server) SayHellow(ctx context.Context, req *service.SayHellowRequest) (*service.SayHellowResponse, error) {
	message := fmt.Sprintf("Hellow %s, this is go !", req.Name)
	split := make([]string, req.Count)
	for i := range split {
		split[i] = message
	}
	return &service.SayHellowResponse{
		Message: strings.Join(split, "\n"),
	}, nil
}

func (s *server) StreamHellow(req *service.SayHellowRequest, stream service.Hellow_StreamHellowServer) error {
	message := fmt.Sprintf("Hellow %s, this is go !", req.Name)
	for i := int64(0); i < req.Count; i++ {
		time.Sleep(time.Second)
		if err := stream.Send(&service.SayHellowResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
	return nil
}

var _ service.HellowServer = (*server)(nil)
