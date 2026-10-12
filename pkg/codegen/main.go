package main

import (
	"os"

	controllergen "github.com/rancher/wrangler/v3/pkg/controller-gen"
	"github.com/rancher/wrangler/v3/pkg/controller-gen/args"
	v1 "k8s.io/api/core/v1"
)

func main() {
	_ = os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/harvester/rancherd/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			v1.GroupName: {
				Types: []interface{}{
					v1.Node{},
				},
				GenerateClients: true,
			},
		},
	})
}
