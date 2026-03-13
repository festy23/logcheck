Prepare logcheck for release. $ARGUMENTS is the version tag (e.g. v0.1.0).

## Steps

1. Delegate to **verifier** agent — full verification must pass
2. Verify plugin builds: `go build ./plugin/`
3. Ensure README.md has: installation instructions, rules table with ❌/✅ examples, configuration example
4. Ensure all exported types and functions have godoc comments
5. Run `go mod tidy`
6. Final `go test ./... -race -count=1`

Report readiness: what's done, what needs manual action (git tag, git push).
