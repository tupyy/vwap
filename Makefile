.PHONY: help tools build check run logs

help: help.all
build: build.local
check: check.imports check.fmt check.lint check.test
run: run.local

# Colors used in this Makefile
escape=$(shell printf '\033')
RESET_COLOR=$(escape)[0m
COLOR_YELLOW=$(escape)[38;5;220m
COLOR_RED=$(escape)[91m
COLOR_BLUE=$(escape)[94m

COLOR_LEVEL_TRACE=$(escape)[38;5;87m
COLOR_LEVEL_DEBUG=$(escape)[38;5;87m
COLOR_LEVEL_INFO=$(escape)[92m
COLOR_LEVEL_WARN=$(escape)[38;5;208m
COLOR_LEVEL_ERROR=$(escape)[91m
COLOR_LEVEL_FATAL=$(escape)[91m

define COLORIZE
sed -u -e "s/\\\\\"/'/g; \
s/method=\([^ ]*\)/method=$(COLOR_BLUE)\1$(RESET_COLOR)/g;        \
s/error=\"\([^\"]*\)\"/error=\"$(COLOR_RED)\1$(RESET_COLOR)\"/g;  \
s/msg=\"\([^\"]*\)\"/msg=\"$(COLOR_YELLOW)\1$(RESET_COLOR)\"/g;   \
s/level=trace/level=$(COLOR_LEVEL_TRACE)trace$(RESET_COLOR)/g;    \
s/level=debug/level=$(COLOR_LEVEL_DEBUG)debug$(RESET_COLOR)/g;    \
s/level=info/level=$(COLOR_LEVEL_INFO)info$(RESET_COLOR)/g;       \
s/level=warning/level=$(COLOR_LEVEL_WARN)warning$(RESET_COLOR)/g; \
s/level=error/level=$(COLOR_LEVEL_ERROR)error$(RESET_COLOR)/g;    \
s/level=fatal/level=$(COLOR_LEVEL_FATAL)fatal$(RESET_COLOR)/g"
endef


#####################
# Help targets      #
#####################

.PHONY: help.highlevel help.all

#help help.highlevel: show help for high level targets. Use 'make help.all' to display all help messages
help.highlevel:
	@grep -hE '^[a-z_-]+:' $(MAKEFILE_LIST) | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'

#help help.all: display all targets' help messages
help.all:
	@grep -hE '^#help|^[a-z_-]+:' $(MAKEFILE_LIST) | sed "s/#help //g" | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {if ($$1 ~ /\./) printf("    $(COLOR_BLUE)%-21s$(RESET_COLOR) %s\n", $$1, $$2); else printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'


#####################
# Tools targets     #
#####################

TOOLS_DIR=$(CURDIR)/tools/bin

.PHONY: tools.clean tools.get

#help tools.clean: remove everything in the tools/bin directory
tools.clean:
	rm -fr $(TOOLS_DIR)/*

#help tools.get: retrieve all the tools specified in gex
tools.get:
	cd $(CURDIR)/tools && go generate tools.go


#####################
# Build targets     #
#####################
NAME=vwat
VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD --abbrev-commit)

IMAGE_TAG=$(VERSION)-$(GIT_COMMIT)
IMAGE_NAME=$(NAME)

GOCACHE?=$(shell go env GOCACHE 2>/dev/null)

ifneq "$(strip $(GOCACHE))" ""
    GOCACHE_FLAGS=-v $(GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go
endif

.PHONY: build.prepare build.vendor build.vendor.full build.docker 

#help build.prepare: prepare target/ folder
build.prepare:
	@mkdir -p $(CURDIR)/target
	@rm -f $(CURDIR)/target/$(NAME)

#help build.vendor: retrieve all the dependencies used for the project
build.vendor:
	go mod vendor

#help build.vendor.full: retrieve all the dependencies after cleaning the go.sum
build.vendor.full:
	@rm -fr $(CURDIR)/vendor

#help build.docker: build a docker image
build.docker:
	DOCKER_BUILDKIT=1 docker build --ssh default --build-arg build_args="$(BUILD_ARGS)"  -t $(IMAGE_NAME):$(IMAGE_TAG) -f Dockerfile .

#help build.local: build locally a binary, in target/ folder
build.local: build.prepare
	go build -mod=vendor $(BUILD_ARGS) -ldflags "-X CommitID=$(GIT_COMMIT) -s -w" \
	-o $(CURDIR)/target/run $(CURDIR)/cmd/$(NAME)-server/main.go

#####################
# Check targets     #
#####################

LINT_COMMAND=golangci-lint run
FILES_LIST=$(shell ls -d */ | grep -v -E "vendor|tools|target|client|restapi|models")
TOOLS_DOCKER_IMAGE=go1.17:buster
MODULE_NAME=$(shell head -n 1 go.mod | cut -d '/' -f 3)

.PHONY: check.fmt check.imports check.lint check.test 

check.prepare:
	@docker pull $(TOOLS_DOCKER_IMAGE)

#help check.fmt: format go code
check.fmt: check.prepare
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)"  $(TOOLS_DOCKER_IMAGE) sh -c 'gofumpt -s -w $(FILES_LIST)'

#help check.imports: fix and format go imports
check.imports: check.prepare
	@# Removes blank lines within import block so that goimports does its magic in a deterministic way
	find $(FILES_LIST) -type f -name "*.go" | xargs -L 1 sed -i '/import (/,/)/{/import (/n;/)/!{/^$$/d}}'
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c 'goimports -w -local github.com/tupyy $(FILES_LIST)'
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c 'goimports -w -local github.com/tupyy/$(MODULE_NAME) $(FILES_LIST)'

#help check.lint: check if the go code is properly written, rules are in .golangci.yml
check.lint: check.prepare
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c '$(LINT_COMMAND)'

#help check.test: execute go tests
check.test: check.prepare
	go test -mod=vendor ./...


#####################
# Run               #
#####################

.PHONY: run.docker.stop run.docker.logs run.infra run.infra.stop 

#help run.docker: run the application on a container
run.docker:
	docker run -d --rm -v $(CURDIR)/resources/:/etc/$(NAME)/ --name $(NAME) $(IMAGE_NAME):$(IMAGE_TAG) --config /etc/$(NAME)/.$(NAME).yaml
	@docker logs -f $(NAME) | $(COLORIZE)

#help run.docker.stop: stop the container of the application
run.docker.stop:
	docker stop $(NAME)

#help run.docker.logs: display logs from the application in the container
run.docker.logs:
	docker logs -f $(NAME) | $(COLORIZE)
