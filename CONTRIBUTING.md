# Contributing to Golang-Vulkan-API

Thank you for your interest in contributing to `golang-vulkan-api`! This document outlines the process for contributing to the repository and provides instructions on how to build, test, and format the code.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally.
3.  **Create a branch** for your feature or bug fix:
    ```bash
    git checkout -b your-feature-branch
    ```

## Prerequisites

To build and test the project, you need the following installed:

*   **Go**: 1.22 or later.
*   **Vulkan SDK or development libraries**: (e.g., `libvulkan-dev` on Linux).
*   **pkg-config**
*   **CGO** must be enabled (`CGO_ENABLED=1`).

See the [main README](README.md#platform-specific-setup) for platform-specific setup instructions.

## Building the Project

We use `make` to simplify the build process. You can see all available commands by running `make help`.

To build the entire project:

```bash
make build
```

## Testing

All code changes should be thoroughly tested. We have a suite of tests including unit tests and tests with the race detector.

To run the standard tests:

```bash
make test
```

To run the tests with verbose output:

```bash
make test-verbose
```

You can also run tests with the Go race detector enabled using standard Go commands:

```bash
go test -race ./...
```

When running tests (`go test ./...`), expect to see NVML CGO deprecation warnings (e.g., related to `go-nvml`). These warnings are expected output and do not indicate a test failure.

## Code Quality and Formatting

We maintain a high standard of code quality. Before submitting a Pull Request, ensure your code passes linting and formatting checks.

To format your code using `gofmt`:

```bash
make format
```

To run the linters (requires `golangci-lint` to be installed):

```bash
make lint
```

The CI pipeline will automatically run `golangci-lint` on your pull requests. Please fix any issues reported by the linter. Note that our configuration includes strict function length (`funlen`) limits; test functions exceeding ~60 lines should be split into smaller sub-functions.

## Writing Tests

If you add a new feature or fix a bug, please include tests.

*   The project tests parameter validation errors using a custom `ValidationError` type. Correct test behavior should be verified by checking `errors.As(err, &valErr)` and asserting that `valErr.Field` matches the expected invalid parameter name.
*   For testing, the project uses a `test_helpers_test.go` convention containing mock generation functions (e.g., `fakeHandle()`, `fakeDevice()`) to create handles backed by real Go memory, which are safe for tests without invoking actual Vulkan C functions.
*   To test success paths of functions that internally call Vulkan C functions (without triggering SIGABRT in headless tests), use the pattern of assigning the target CGO-wrapping function to a package-level variable so it can be overridden and mocked in tests.

## Submitting a Pull Request

Once your changes are complete, tested, and formatted:

1.  **Commit your changes**. Ensure your commit messages are clear and descriptive.
2.  **Push to your fork** on GitHub.
3.  **Open a Pull Request** against the `main` branch.

### Pull Request Description

When creating your PR, please provide a clear description.
A good PR description should include:

*   🎯 **What:** What changes were made.
*   📊 **Why/Coverage:** Why the changes were necessary or what test coverage was added.
*   ✨ **Result:** The overall impact of the PR.

Example:

```markdown
🎯 **What:** Added a new test file `video_helpers_test.go`.
📊 **Coverage:** Added test coverage for `CreateVideoDeviceFunctions` and `GetVideoDeviceFunctions`.
✨ **Result:** Improved test coverage providing a baseline safety net for refactoring.
```

## Review Process

Once you submit your PR, maintainers will review your code. We may request changes or ask questions to ensure the code meets the project's standards.

Thank you for contributing!
