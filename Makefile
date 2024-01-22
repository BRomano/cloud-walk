mocks:
	@which mockgen >/dev/null || (echo "Installing mockgen..." && go install go.uber.org/mock/mockgen@v0.3.0)
	go generate -x ./...


lint:
	@which golangci-lint >/dev/null || (echo "Installing lint...") && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
	golangci-lint run -v --deadline=200s --skip-dirs=vendor --disable-all --enable=revive --enable=unused --enable=goconst