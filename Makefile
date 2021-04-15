GO_PACKAGES = $$(go list ./... | grep -v vendor )
GO_FILES = $$(find . -name "*.go" | grep -v vendor | uniq)

myretail-mac:
	GODEBUG=tls13=0 GOOS=darwin GOARCH=amd64 go build -o myretail.darwin ./cmd/myRetail/main.go

myretail-linux:
   	GODEBUG=tls13=0 GOOS=linux GOARCH=amd64 go build -o myretail.linux ./cmd/myRetail/main.go

myretail-windows:
	GODEBUG=tls13=0 GOOS=windows GOARCH=amd64 go build -o myretail.windows ./cmd/myRetail/main.go

build: myretail-mac myretail-linux myretail-windows

generate:
	go generate ./...

unit-test:
	@go test ${GO_PACKAGES}

vet:
	@go vet ${GO_PACKAGES}

test: generate unit-test vet

cleandep:
	go mod tidy

bootstrap:
	go install "github.com/maxbrunsfeld/counterfeiter/v6"
	go install "github.com/onsi/ginkgo"
	go install "github.com/onsi/gomega"
	go install "golang.org/x/tools/cmd/goimports"
	go install "github.com/kelseyhightower/envconfig"

fmt:
	goimports -l -w $(GO_FILES)

all: fmt test build