test:
	CGO_ENABLED=0 go test -mod=mod -short -cover -v -coverprofile=coverage.out -covermode=atomic ./...

# test-all runs all tests including integration tests
test-all:
	CGO_ENABLED=0 go test -mod=mod -cover -v -coverprofile=coverage.out -covermode=atomic ./...

# Makes the test runner stop after the first failing test
# (however, tests running in parallel to the test in question will finish)
test-failfast:
	CGO_ENABLED=0 go test -failfast -mod=mod -short -cover -v -coverprofile=coverage.out -covermode=atomic ./...
