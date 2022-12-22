module github.com/drone/drone-gc

go 1.19

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-units v0.3.2
	github.com/drone/signal v0.0.0-20170915013802-ac5d07ef1315
	github.com/golang/mock v1.4.3
	github.com/google/go-cmp v0.4.0
	github.com/hashicorp/go-multierror v0.0.0-20171204182908-b7773ae21874
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/rs/zerolog v1.6.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.7 // indirect
	github.com/containerd/containerd v1.3.4 // indirect
	github.com/docker/go-connections v0.3.0 // indirect
	github.com/gogo/protobuf v0.0.0-20170307180453-100ba4e88506 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/hashicorp/errwrap v0.0.0-20141028054710-7554cd9344ce // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
