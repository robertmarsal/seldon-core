package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"github.com/seldonio/seldon-core/executor/api"
	seldonclient "github.com/seldonio/seldon-core/executor/api/client"
	"github.com/seldonio/seldon-core/executor/api/grpc"
	"github.com/seldonio/seldon-core/executor/api/grpc/seldon"
	"github.com/seldonio/seldon-core/executor/api/grpc/seldon/proto"
	"github.com/seldonio/seldon-core/executor/api/grpc/tensorflow"
	"github.com/seldonio/seldon-core/executor/api/kafka"
	"github.com/seldonio/seldon-core/executor/api/rest"
	"github.com/seldonio/seldon-core/executor/api/tracing"
	"github.com/seldonio/seldon-core/executor/k8s"
	loghandler "github.com/seldonio/seldon-core/executor/logger"
	predictor2 "github.com/seldonio/seldon-core/executor/predictor"
	"github.com/seldonio/seldon-core/executor/proto/tensorflow/serving"
	"github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"net"
	"net/url"
	"os"
	"os/signal"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"syscall"
	"time"
)

var (
	configPath     = flag.String("config", "", "Path to kubconfig")
	sdepName       = flag.String("sdep", "", "Seldon deployment name")
	namespace      = flag.String("namespace", "", "Namespace")
	predictorName  = flag.String("predictor", "", "Name of the predictor inside the SeldonDeployment")
	httpPort       = flag.Int("http_port", 8080, "Executor port")
	grpcPort       = flag.Int("grpc_port", 8000, "Executor port")
	wait           = flag.Duration("graceful_timeout", time.Second*15, "Graceful shutdown secs")
	protocol       = flag.String("protocol", "seldon", "The payload protocol")
	transport      = flag.String("transport", "rest", "The network transport http or grpc")
	filename       = flag.String("file", "", "Load graph from file")
	hostname       = flag.String("hostname", "", "The hostname of the running server")
	logWorkers     = flag.Int("logger_workers", 5, "Number of workers handling payload logging")
	prometheusPath = flag.String("prometheus_path", "/metrics", "The prometheus metrics path")
	kafkaBroker    = flag.String("kafka_broker", "", "The kafka broker as host:port")
	kafkaTopicIn   = flag.String("kafka_topic_in", "", "The kafka input topic")
	kafkaTopicOut  = flag.String("kafka_topic_out", "", "The kafka output topic")
)

func getServerUrl(hostname string, port int) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://%s:%d/", hostname, port))
}

func runHttpServer(logger logr.Logger, predictor *v1.PredictorSpec, client seldonclient.SeldonApiClient, port int,
	probesOnly bool, serverUrl *url.URL, namespace string, protocol string, deploymentName string, prometheusPath string) {

	// Create REST API
	seldonRest := rest.NewServerRestApi(predictor, client, probesOnly, serverUrl, namespace, protocol, deploymentName, prometheusPath)
	seldonRest.Initialise()
	srv := seldonRest.CreateHttpServer(port)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err, "Server error")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) and SIGTERM
	// SIGKILL, SIGQUIT will not be caught.
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("shutting down")
	os.Exit(0)

}

func runGrpcServer(logger logr.Logger, predictor *v1.PredictorSpec, client seldonclient.SeldonApiClient, port int, serverUrl *url.URL, namespace string, protocol string, deploymentName string, annotations map[string]string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer, err := grpc.CreateGrpcServer(predictor, deploymentName, annotations, logger)
	if err != nil {
		log.Fatalf("Failed to create grpc server: %v", err)
	}
	if protocol == api.ProtocolSeldon {
		seldonGrpcServer := seldon.NewGrpcSeldonServer(predictor, client, serverUrl, namespace)
		proto.RegisterSeldonServer(grpcServer, seldonGrpcServer)
	} else {
		tensorflowGrpcServer := tensorflow.NewGrpcTensorflowServer(predictor, client, serverUrl, namespace)
		serving.RegisterPredictionServiceServer(grpcServer, tensorflowGrpcServer)
		serving.RegisterModelServiceServer(grpcServer, tensorflowGrpcServer)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Errorf("Grpc server error: %v", err)
	}
}

func main() {
	flag.Parse()

	if *sdepName == "" {
		log.Fatal("Required argument sdep missing")
	}

	if *namespace == "" {
		log.Fatal("Required argument namespace missing")
	}

	if *predictorName == "" {
		log.Fatal("Required argument predictor missing")
	}

	if !(*protocol == api.ProtocolSeldon || *protocol == api.ProtocolTensorflow) {
		log.Fatal("Invalid protocol: must be seldon or tensorflow")
	}

	if !(*transport == api.TransportRest || *transport == api.TransportGrpc || *transport == api.TransportKafka) {
		log.Fatal("Invalid transport: Only rest, grpc and kafka supported")
	}

	if *transport == api.TransportKafka {
		if *kafkaBroker == "" {
			log.Fatal("Required argument kafka_broker missing")
		}
		if *kafkaTopicIn == "" {
			log.Fatal("Required argument kafka_topic_in missing")
		}
	}

	logf.SetLogger(logf.ZapLogger(false))
	logger := logf.Log.WithName("entrypoint")

	// Set hostname
	if *hostname == "" {
		*hostname = os.Getenv("POD_NAME")
		if *hostname == "" {
			logger.Info("Hostname unset will use localhost")
			*hostname = "localhost"
		} else {
			logger.Info("Hostname found from env", "hostname", *hostname)
		}
	} else {
		logger.Info("Hostname provided on command line", "hostname", *hostname)
	}
	serverUrl, err := getServerUrl(*hostname, *httpPort)
	if err != nil {
		log.Fatal("Failed to create server url from", *hostname, *httpPort)
	}

	logger.Info("Flags", "transport", *transport)

	predictor, err := predictor2.GetPredictor(*predictorName, *filename, *sdepName, *namespace, configPath)
	if err != nil {
		log.Error(err, "Failed to get predictor")
		os.Exit(-1)
	}

	// Ensure standard OpenAPI seldon API file has this deployment's values
	err = rest.EmbedSeldonDeploymentValuesInSwaggerFile(*namespace, *sdepName)
	if err != nil {
		log.Error(err, "Failed to embed variables on OpenAPI template")
	}

	annotations, err := k8s.GetAnnotations()
	if err != nil {
		log.Error(err, "Failed to load annotations")
	}

	//Start Logger Dispacther
	loghandler.StartDispatcher(*logWorkers, logger, *sdepName, *namespace, *predictorName)

	//Init Tracing
	if *transport != api.TransportKafka {
		closer, err := tracing.InitTracing()
		if err != nil {
			log.Fatal("Could not initialize jaeger tracer", err.Error())
		}
		defer closer.Close()
	}

	switch *transport {
	case api.TransportGrpc:
		logger.Info("Running http probes only server ", "port", *httpPort)
		go runHttpServer(logger, predictor, nil, *httpPort, true, serverUrl, *namespace, *protocol, *sdepName, *prometheusPath)
		logger.Info("Running grpc server ", "port", *grpcPort)
		var clientGrpc seldonclient.SeldonApiClient
		if *protocol == "seldon" {
			clientGrpc = seldon.NewSeldonGrpcClient(predictor, *sdepName, annotations)
		} else {
			clientGrpc = tensorflow.NewTensorflowGrpcClient(predictor, *sdepName, annotations)
		}
		runGrpcServer(logger, predictor, clientGrpc, *grpcPort, serverUrl, *namespace, *protocol, *sdepName, annotations)
	case api.TransportRest:
		clientRest, err := rest.NewJSONRestClient(*protocol, *sdepName, predictor, annotations)
		if err != nil {
			log.Fatalf("Failed to create http client: %v", err)
		}
		logger.Info("Running http server ", "port", *httpPort)
		runHttpServer(logger, predictor, clientRest, *httpPort, false, serverUrl, *namespace, *protocol, *sdepName, *prometheusPath)
	case api.TransportKafka:
		kafkaServer := kafka.NewKafkaServer(*sdepName, *namespace, *protocol, annotations, serverUrl, predictor, *kafkaBroker, *kafkaTopicIn, *kafkaTopicOut, logger)
		err := kafkaServer.Serve()
		if err != nil {
			log.Fatal("Failed to serve kafka", err)
		}
	}

}
