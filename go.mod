module github.com/confluentinc/cli

require (
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/Shopify/sarama v1.20.1
	github.com/Shopify/toxiproxy v2.1.3+incompatible // indirect
	github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da
	github.com/codyaray/go-editor v0.3.0
	github.com/codyaray/go-printer v0.9.0
	github.com/codyaray/retag v0.0.0-20180529164156-4f3c7e6dfbe2 // indirect
	github.com/confluentinc/cc-structs v0.0.0-20190216225128-bc354c6bf010
	github.com/confluentinc/ccloud-sdk-go v0.0.6-0.20190226163025-48f4ae5f158f
	github.com/confluentinc/ccloudapis v0.0.0-20190227065628-cf9b2483670a
	github.com/confluentinc/go-printer v0.9.0
	github.com/confluentinc/protoc-gen-ccloud v0.0.1 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.2.1-0.20181128192352-1d3f30b51784
	github.com/golang/snappy v0.0.1 // indirect
	github.com/golangci/errcheck v0.0.0-20181003203344-ef45e06d44b6 // indirect
	github.com/golangci/golangci-lint v1.12.2
	github.com/google/uuid v1.1.0
	github.com/goreleaser/goreleaser v0.101.0
	github.com/hashicorp/go-hclog v0.0.0-20180910232447-e45cbeb79f04
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-plugin v0.0.0-20181030172320-54b6ff97d818
	github.com/hashicorp/yamux v0.0.0-20180826203732-cc6d2ea263b2 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.0.0
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.0-20180912035003-be2c049b30cc // indirect
	github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.2.0
	github.com/stretchr/testify v1.3.0
	github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43 // indirect
	golang.org/x/crypto v0.0.0-20190219172222-a4c6cb3142f2
	golang.org/x/oauth2 v0.0.0-20190220154721-9b3c75971fc9 // indirect
	golang.org/x/sys v0.0.0-20190221222158-ec7b60b042fd // indirect
	golang.org/x/tools v0.0.0-20190221204921-83362c3779f5 // indirect
	google.golang.org/genproto v0.0.0-20190219182410-082222b4a5c5 // indirect
	google.golang.org/grpc v1.18.0
	k8s.io/api v0.0.0-20181221193117-173ce66c1e39 // indirect
	k8s.io/apiextensions-apiserver v0.0.0-20190103235604-e7617803aceb // indirect
	k8s.io/apimachinery v0.0.0-20190109170643-c3a4c8673eae // indirect
	k8s.io/kube-openapi v0.0.0-20181114233023-0317810137be // indirect
)

replace (
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20171026124306-e509bb64fe11
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20170925234155-019ae5ada31d
)