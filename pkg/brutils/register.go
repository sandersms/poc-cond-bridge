// SPDX-License-Identifier: Apache-2.0

// utility support functions for generic registration handler
package brutils

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// generic function for registering the endpoint to the gateway
type registerHandlerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

func RegisterGatewayHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption, registerFunc registerHandlerFunc, serviceName string) {
	err := registerFunc(ctx, mux, endpoint, opts)
	if err != nil {
		log.Panicf("cannot register %s handler server: %v", serviceName, err)
	}
}
