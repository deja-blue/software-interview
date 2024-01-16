package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"time"

	"log"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/deja-blue/software-interview/go/pkg/charge"
	"github.com/deja-blue/software-interview/go/pkg/gen/gql/resolver"
	"github.com/deja-blue/software-interview/go/pkg/gen/gql/schema"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	chargeListenPort = flag.String("charge-http-listen-port", "8080", "HTTP port to listen to for charge service")

	//go:embed playground.gohtml
	embeddedFiles embed.FS
)

func main() {
	flag.Parse()
	port, err := strconv.Atoi(*chargeListenPort)
	if err != nil {
		log.Fatalf("invalid charge-http-listen-port: %v", *chargeListenPort)
	}
	res := &resolver.Resolver{
		ChargeResolver: charge.NewResolver(
			charge.MockCharger("ABC", 7),
			charge.MockCharger("DEF", 7,
				charge.WithStateOfCharge(&charge.VehicleStateOfCharge{
					CurrentBatteryLevelKwH: 30,
					MaxBatteryLevelKwH:     70,
					RangeKmPerKwH:          5,
				}),
			),
			charge.MockCharger("GHI", 22,
				charge.WithStateOfCharge(&charge.VehicleStateOfCharge{
					CurrentBatteryLevelKwH: 30,
					MaxBatteryLevelKwH:     70,
					RangeKmPerKwH:          5,
				}),
				charge.WithDynamicStatus(time.Second*5)),
		),
	}

	srv := handler.New(schema.NewExecutableSchema(schema.Config{Resolvers: res}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), r))
}
