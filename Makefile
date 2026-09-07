GO_TEST_FLAGS=-v -race
GO_COVERAGE_FILE=coverage.out
GOUTIL_ENABLED?=YES
GOUTIL_VERSION=1.26.6.2
GOUTIL_IMAGE=registry.gitlab.b2broker.tech/b2connect/b2connect/libs/docker/dockerfiles/goutil
ifeq ("$(GOUTIL_ENABLED)", "NO")
GOUTIL=
else
GOUTIL=docker run --rm \
	-e CI_SERVER_HOST -e CI_GO_PRIVATE_TOKEN -e CI_JOB_TOKEN \
	-e CGO_ENABLED=1 \
	-v "$(shell pwd):/src" \
	-u "$(shell id -u):$(shell id -g)" \
	$(GOUTIL_IMAGE):$(GOUTIL_VERSION)
endif

.PHONY: test_unit
test_unit:
	$(GOUTIL) go test $(GO_TEST_FLAGS) ./...

.PHONY: test_integration
test_integration:
	@echo "no integration tests"

.PHONY: coverage
coverage:
	$(GOUTIL) sh -c 'COVERFILE="$$(mktemp)" \
		&& go test -coverprofile "$${COVERFILE}" -coverpkg ./... ./... \
		&& (([ -f .coverignore ] && cat "$${COVERFILE}" | grep -vE "($$(printf "|%s" $$(cat .coverignore) | cut -c 2-))") || cat "$${COVERFILE}") > $(GO_COVERAGE_FILE)'

.PHONY: lint
lint:
	$(GOUTIL) golangci-lint run --config .golangci.yaml -v ./...

.PHONY: lint_fix
lint_fix:
	$(GOUTIL) golangci-lint run --config .golangci.yaml -v --fix ./...

.PHONY: gen_check
gen_check:
	@echo "no generated code"
