test:
	go run gotest.tools/gotestsum@latest --format=testdox  -- -covermode=atomic -coverprofile=coverage.txt ./...
benchmark:
	go test -bench=. -benchmem ./...
format:
	go fmt ./...
format-check:
	@UNFORMATTED_FILES=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED_FILES" ]; then \
		echo "::error::The following files are not formatted correctly:"; \
		echo "$$UNFORMATTED_FILES"; \
		echo "--- Diff ---"; \
		gofmt -d .; \
		exit 1; \
	fi

.PHONY: test format benchmark format-check