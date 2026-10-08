MODULES := common echo funcs

.PHONY: deps test-all fmt vet lint clean $(addprefix test-,$(MODULES))

deps:
	@for m in $(MODULES); do \
		echo ">> deps $$m"; \
		(cd lambdas-go/$$m && GO111MODULE=on go mod download) || exit 1; \
	done

test-all: deps $(addprefix test-,$(MODULES))

$(addprefix test-,$(MODULES)): test-%:
	@echo ">> test $*"
	@cd lambdas-go/$* && GO111MODULE=on go test -v ./...

fmt:
	@for m in $(MODULES); do (cd lambdas-go/$$m && go fmt ./...) || exit 1; done

vet:
	@for m in $(MODULES); do (cd lambdas-go/$$m && go vet ./...) || exit 1; done

lint:
	@golangci-lint run ./...

clean:
	@find . -name bootstrap -type f -delete
	@rm -rf ../.aws-sam
