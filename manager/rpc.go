// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package manager

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ChrisMcKenzie/dropship/dropship"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ServiceServer struct {
	store store.Store
}

func serveRpc(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	dropship.RegisterRpcServiceServer(grpcServer, &RpcServiceServer{})
	// determine whether to use TLS
	log.Infof("RPC Server Listening on port %d", port)
	return grpcServer.Serve(lis)
}

func (s RpcServiceServer) RegisterService(ctx context.Context, svc *dropship.Service) (*dropship.RegisterResponse, error) {
	log.Infof("RegisterService Request Received for %s", svc.Name)
	fmt.Println(svc)

	// do stuff here
	// if service exists then do nothing.
	path := fmt.Sprintf(
		"%s/services/%s",
		DefaultKeyPrefix,
		svc.Name,
	)

	serviceExists, err := s.store.Exists(path)
	if err != nil {
		return &dropship.RegisterResponse{
			Success: false,
		}, err
	}

	fmt.Println(serviceExists)
	if !serviceExists {
		log.Infof("Registering %s service", svc.Name)

		data, err := json.Marshal(Service{
			Name: svc.Name,
		})

		err := s.store.Put(path, []byte(svc.Ip), nil)
		if err != nil {
			return &dropship.RegisterResponse{
				Success: false,
			}, err
		}
	}

	return &dropship.RegisterResponse{
		Success: true,
	}, nil
}
