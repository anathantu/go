module github/anathantu/go

go 1.16

replace github.com/anathantu/go => ../go

require (
	github.com/dmitriyGarden/consul-leader-election v1.2.0
	github.com/hashicorp/consul/api v1.11.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/protobuf v1.27.1
)
