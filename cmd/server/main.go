package main

import (
	"os"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/pkg/env"
	"github.com/go-kratos/kratos-layout/pkg/registry"

	zaplog "github.com/go-kratos/kratos-layout/pkg/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// id is the service instance id.
	id string
)

func init() {
	json.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	var err error
	id, err = os.Hostname()
	if err != nil {
		id = "unknown"
	}

	if Name == "" {
		Name = env.GetOrDefault("SERVICE_NAME", "kratos_layout")
	}

	if Version == "" {
		Version = env.GetOrDefault("SERVICE_VERSION", "0.0.1")
	}
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, r *nacos.Registry) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	logger := zaplog.InitDefaultLogger(zapcore.DebugLevel)

	c := config.New(
		config.WithSource(
			apollo.NewSource(
				apollo.WithAppID(env.GetOrDefault("APOLLO_APP_ID", "kratos_layout")),
				apollo.WithCluster(env.GetOrDefault("APOLLO_CLUSTER", "dev")),
				apollo.WithEndpoint(env.GetOrDefault("APOLLO_ENDPOINT", "http://localhost:8080")),
				apollo.WithNamespace(env.GetOrDefault("APOLLO_NAMESPACE", "application,bootstrap.yaml")),
				apollo.WithSecret(env.GetOrDefault("APOLLO_SECRET", "fc4cacadc4cb486b91419d67f6d7918b")),
			),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Value("bootstrap").Scan(&bc); err != nil {
		panic(err)
	}

	r, err := registry.NewNacosRegistryFromEnv()
	if err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, r, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
