# Changelog: Changes Since v1.0.0

## Overview
This document summarizes all changes made to the SerpApi Go library since tag `1.0.0` was created.

## v1.3.0 - 2026-09-24

- Added `Context` variants (`SearchContext`, `HtmlContext`, `MarkdownContext`, `LocationContext`, `AccountContext`, `SearchArchiveContext`) for cancellation and deadlines. Existing methods are unchanged wrappers.
- Added `*serpapi.HTTPError` (status code, URL, body) for non-2xx responses.
- Security: the API key is redacted from URLs in `HTTPError` and network errors.
- HTTP transport now starts from Go defaults: honors `HTTP_PROXY`/`HTTPS_PROXY`, enables HTTP/2, defaults `KeepAlive` to 60s, applies `MaxIdleConnection` per host.
- Search archive IDs are escaped in request paths; response bodies are always closed.
- Pagination example in README "Simple Usage" and `demo/demo.go` (follows `serpapi_pagination.next`).
- Examples under `test/example/` aligned with serpapi-ruby: added `google_ai_mode` and `google_light_search`, removed `google_ai_overview`. Added a markdown test and offline request-construction tests; CI runs `go vet` and library tests.

**Behavior changes:**
- Non-2xx error messages are now `serpapi request failed: <status>: <body>` instead of the JSON `error` value.
- `Html` returns an error on non-2xx responses instead of the error page body.

## v1.2.0 - 2026-08-16

- Added `Markdown()` returning `output=md` results optimized for LLMs and AI agents.

## v1.1.0 - 2026-01-26

**Statistics:**
- **14 commits** since v1.0.0
- **56 files changed**: 3,072 insertions(+), 798 deletions(-)
- **Code coverage**: Improved to 81%

---

## Major Features & Enhancements

### 1. Asynchronous & Persistent Mode Support
- **Added support for asynchronous search mode** - allows non-blocking API calls
- **Added support for persistent mode** - enables connection reuse and better performance
- New demo files demonstrating async usage:
  - `demo/demo_async.go` - asynchronous search examples
  - `demo/demo_thread_pool.go` - thread pool implementation examples
  - `demo/demo.go` - general usage examples

### 2. Improved API Key Handling
- **Simplified API key handling** for better memory safety in modern Go
- Enhanced security practices for API key management

### 3. Client Configuration Improvements
- Enhanced client settings to support:
  - Persistent connections
  - Asynchronous operations
  - Customizable timeouts
  - Connection pooling (MaxIdleConnection, KeepAlive, TLSHandshakeTimeout)

---

## New Test Examples & Coverage

### New Google Search Examples Added:
- `example_search_google_ai_overview_test.go` - AI Overview search
- `example_search_google_finance_test.go` - Finance search
- `example_search_google_flights_test.go` - Flights search
- `example_search_google_hotels_test.go` - Hotels search
- `example_search_google_images_light_test.go` - Light Images search
- `example_search_google_immersive_product_test.go` - Immersive Product search
- `example_search_google_lens_test.go` - Google Lens search
- `example_search_google_light_test.go` - Light search
- `example_search_google_news_light_test.go` - Light News search
- `example_search_google_news_test.go` - News search
- `example_search_google_patents_test.go` - Patents search
- `example_search_google_play_test.go` - Play Store search
- `example_search_google_shopping_test.go` - Shopping search
- `example_search_google_trends_test.go` - Trends search
- `example_search_google_videos_test.go` - Videos search

### Other Search Engine Examples:
- `example_search_bing_test.go` - Bing search
- `example_search_yandex_test.go` - Yandex search
- `example_search_yelp_test.go` - Yelp search

### Test Infrastructure Improvements:
- Improved error handling in test account
- Fixed broken tests
- Better test organization (moved examples to `test/example/` directory)
- Code coverage improved to **81%**

---

## Documentation & Code Quality

### Documentation Enhancements:
- **Major README.md improvements** (+1,564 lines of documentation)
- Updated README documentation
- Improved code examples throughout
- Fixed wording to align with Ruby library
- Added comprehensive usage examples

### Code Quality Improvements:
- Added `go.sum` for module dependency checksums
- Fixed typos across the codebase
- Improved code examples
- Enhanced error handling

---

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

---

## File Structure Changes

### New Files Added:
- `.gitignore` - Git ignore rules
- `go.sum` - Go module checksums
- `demo/demo.go` - Demo examples
- `demo/demo_async.go` - Async demo examples
- `demo/demo_thread_pool.go` - Thread pool demo
- 20+ new test example files

### Files Removed:
- `oobt/demo.go` - Removed (replaced by new demo structure)

### Files Modified:
- `serpapi.go` - Core library improvements (+137 lines)
- All test files - Improved and reorganized
- `Makefile` - Enhanced build targets
- `README.md` - Major documentation update

---

## Commit History

1. **a523a21** (9 weeks ago) - add go.sum for module dependency checksums
2. **078937d** (9 weeks ago) - fix a few typo
3. **f1bfb86** (7 months ago) - gitaction run test using `make test`
4. **59f728c** (7 months ago) - add `make clean` target
5. **32fbe75** (7 months ago) - improve code coverage
6. **f7139b8** (7 months ago) - bulk improvements to documentation, examples, test generate code coverage data `make coverage` 81% allow `make vet` and `make fmt` to support multiple standalone example under demo/*.go
7. **d469f1f** (8 months ago) - fix wording to align with Ruby library
8. **a8d9abc** (8 months ago) - simplify API key handling for better memory safe in modern go
9. **1328aca** (8 months ago) - [test] improve error handling in test account
10. **8b48832** (8 months ago) - fix broken tests
11. **7ff502f** (8 months ago) - improve client setting to support persistent and asynchronous mode
12. **9c344c2** (8 months ago) - update examples and documentation. add async demo
13. **c29ca33** (2 years, 1 month ago) - improve code example
14. **fc37030** (2 years, 1 month ago) - minor tweaks

---

## Summary

Since v1.0.0, the SerpApi Go library has seen significant improvements in:

1. **Functionality**: Added async and persistent mode support
2. **Documentation**: Major documentation overhaul with comprehensive examples
3. **Testing**: Expanded test coverage to 81% with many new example tests
4. **Code Quality**: Improved error handling, memory safety, and code organization
5. **Developer Experience**: Better build system, CI/CD, and demo examples

The library is now more robust, well-documented, and feature-rich, making it easier for developers to integrate SerpApi into their Go applications.
