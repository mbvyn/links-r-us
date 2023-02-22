.PHONY: deps lint lint-check-deps ci-check run-migrations test
include .env
export

ci-check: deps lint run-cdb-migrations test

deps:
	@if [ "$(go mod help | echo 'no-mod')" = "no-mod" ] || [ "${GO111MODULE}" = "off" ]; then \
		echo "[dep] fetching package dependencies";\
		go get -u github.com/golang/dep/cmd/dep;\
		dep ensure;\
	fi
lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		GO111MODULE=$(GO111MODULE) go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

lint:
	@if [ -z `which golangci-lint` ]; then \
		echo "golangci-lint not found in PATH"; \
	else \
		golangci-lint run \
			-E misspell \
			-E golint \
			-E gofmt \
			-E unconvert \
			--exclude-use-default=false \
			./...; \
	fi

migrate-check-deps:
	@if [ -z `which migrate` ]; then \
		echo "[go get] installing golang-migrate cmd with cockroachdb support";\
		if [ "${GO111MODULE}" = "off" ]; then \
			echo "[go get] installing github.com/golang-migrate/migrate/cmd/migrate"; \
			go get -tags 'cockroachdb postgres' -u github.com/golang-migrate/migrate/cmd/migrate;\
		else \
			echo "[go get] installing github.com/golang-migrate/migrate/v4/cmd/migrate"; \
			go get -tags 'cockroachdb postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate;\
		fi \
	fi

run-cdb-migrations: migrate-check-deps
	migrate -source file://linkgraph/store/cdb/migrations -database '$(subst postgresql,cockroach,${CDB_DSN})' up

run-db-migrations: run-cdb-migrations

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...
