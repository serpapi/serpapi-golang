# Highlights: Changes Since v1.0.0

## Major Features

### Context-aware requests
- Added `Context` variants for `Search`, `Html`, `Markdown`, `Location`, `Account`, and `SearchArchive`.
- Existing methods remain backward-compatible.
- Requests can now be cancelled or bounded by caller-provided deadlines.

### Asynchronous & Persistent Mode Support
- Added support for asynchronous search mode - allows non-blocking API calls
- Added support for persistent mode - enables connection reuse and better performance
- New demo files demonstrating async usage:
  - `demo/demo_async.go` - asynchronous search examples
  - `demo/demo_thread_pool.go` - thread pool implementation examples
  - `demo/demo.go` - general usage examples

### Improved API Key Handling
- Simplified API key handling for better memory safety in modern Go
- Enhanced security practices for API key management

### Client Configuration Improvements
- Enhanced client settings to support:
  - Persistent connections
  - Asynchronous operations
  - Customizable timeouts
  - Connection pooling (MaxIdleConnection, KeepAlive, TLSHandshakeTimeout)

## New Test Examples

### New Google Search Examples Added:
- Google AI Overview search
- Google Finance search
- Google Flights search
- Google Hotels search
- Google Images Light search
- Google Immersive Product search
- Google Lens search
- Google Light search
- Google News Light search
- Google News search
- Google Patents search
- Google Play Store search
- Google Shopping search
- Google Trends search
- Google Videos search

### Other Search Engine Examples:
- Bing search
- Yandex search
- Yelp search

## Documentation & Code Quality

### Documentation Enhancements:
- Major README.md improvements (+1,564 lines of documentation)
- Updated README documentation
- Improved code examples throughout
- Fixed wording to align with Ruby library
- Added comprehensive usage examples

### Code Quality Improvements:
- Added `go.sum` for module dependency checksums
- Fixed typos across the codebase
- Improved code examples
- Enhanced error handling
- Code coverage improved to **81%**

## Build System & CI/CD

### Makefile Enhancements:
- Added `make clean` target
- Added `make coverage` target (generates code coverage reports)
- Enhanced `make vet` and `make fmt` to support multiple standalone examples under `demo/*.go`
- Improved build targets and workflow

### CI/CD Improvements:
- Updated GitHub Actions workflow (renamed `go.yml` → `ci.yml`)
- CI now runs tests using `make test`
- Better integration with automated testing

## Statistics

- **14 commits** since v1.0.0
- **56 files changed**: 3,072 insertions(+), 798 deletions(-)
- **Code coverage**: Improved to 81%
