package balance

import (
	"errors"
	"fmt"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

// resolver is the implementaion of grpc.naming.Resolver
type resolver struct {
	serviceName string // service name to resolve
}

// NewResolver return resolver with service name
func NewResolver(serviceName string) *resolver {
	fmt.Println("serviceName: ", serviceName)
	return &resolver{serviceName: serviceName}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"

func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("grpclb: no service name provided")
	}

	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}

	// Return watcher
	return &watcher{re: re, client: *client}, nil
}
