package grpcserver

import (
	"fmt"
	"net"
	"net/http"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/patrick/test-grpc-docker-gitactions/proto/userpb"
	user_server "github.com/patrick/test-grpc-docker-gitactions/server/user"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func StartGrpc() (*grpc.Server, error) {
	grpc_port := "9010"

	maxMSGSize := 1000 * 1024 * 1024

	fmt.Println("initializing grpc server: port:", grpc_port)

	interceptors := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
	}

	options := []grpc.ServerOption{
		grpc.MaxSendMsgSize(maxMSGSize),
		grpc.MaxRecvMsgSize(maxMSGSize),
		grpc.ChainUnaryInterceptor(interceptors...),
	}

	grpcServer := grpc.NewServer(options...)

	userpb.RegisterUserServiceServer(grpcServer, &user_server.UserServer{})

	addr := ":" + grpc_port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Errorf("error creating listener: %v", err.Error())
	}

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			fmt.Println("failed to start server", " error ", err)

		}
	}()
	fmt.Println("grpc server is listening", " port ", grpc_port)

	return grpcServer, nil
}

func StartGrpcWeb(s *grpc.Server) (*http.Server, error) {
	grpc_webport := "9011"
	fmt.Println("initializing grpc server: web port:", grpc_webport)
	fmt.Println("initializing grpc web proxy server", " port ", grpc_webport)
	grpcWebServer := grpcweb.WrapServer(s)

	addr := ":" + grpc_webport
	httpServer := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println("failed to start proxy server", " error ", err)
		}
	}()

	// run goclearing side task
	//RunBackGroundTask()

	fmt.Println("grpc web proxy server listen", " port ", grpc_webport)

	return httpServer, nil
}
