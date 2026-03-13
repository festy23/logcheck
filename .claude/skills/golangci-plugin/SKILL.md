---
name: golangci-plugin
description: >
  How to build and configure golangci-lint module plugin for logcheck.
  Use this skill when working on plugin/plugin.go, .custom-gcl.yml,
  or when asked about golangci-lint integration, installation,
  or plugin configuration.
---

# golangci-lint Module Plugin

## plugin/plugin.go

```go
package main

import (
    "github.com/ivankonovalov/logcheck/pkg/analyzer"
    "github.com/golangci/plugin-module-register/register"
)

func init() {
    register.Plugin("logcheck", register.LoadModeTypesInfo, func(conf any) {
        // parse conf for custom settings if needed
        return analyzer.NewAnalyzer(), nil
    })
}
```

IMPORTANT: use `LoadModeTypesInfo`, NOT `LoadModeSyntax`.
Our analyzer needs type information for package resolution.

## .custom-gcl.yml

```yaml
version: v2.10.1
plugins:
  - module: 'github.com/ivankonovalov/logcheck'
    import: 'github.com/ivankonovalov/logcheck/plugin'
    path: .
```

## Build & run

```bash
golangci-lint custom        # builds ./custom-gcl binary
./custom-gcl run --enable logcheck ./...
```

## .golangci.yml example for users

```yaml
version: "2"
linters:
  default: none
  enable:
    - logcheck
  settings:
    custom:
      logcheck:
        type: "module"
        description: "Checks log messages for common issues"
        settings:
          disabled-rules: []
          sensitive-patterns:
            - api_secret
            - jwt_token
```
