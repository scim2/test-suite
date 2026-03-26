# SCIM 2.0 Compliance Test Suite

> **Status: Initial Draft** - One pass through the SCIM RFCs has been completed.
> Coverage is limited and many requirements are still untested. Expect significant
> changes as additional iterations are done.

Black-box compliance tests for SCIM 2.0 servers based on
[RFC 7642](https://tools.ietf.org/html/rfc7642) (Definitions, Overview, Concepts, and Requirements),
[RFC 7643](https://tools.ietf.org/html/rfc7643) (Core Schema), and
[RFC 7644](https://tools.ietf.org/html/rfc7644) (Protocol).

Each testable requirement is extracted from the RFC text with its RFC 2119 keyword
(MUST, SHOULD, MAY, etc.) and mapped to one or more Go test functions. After a run
a compliance report is generated showing pass/fail/warn per requirement and overall
coverage.

## Project layout

| Package | Purpose |
|---|---|
| `spec/` | Parsed RFC requirements with IDs, compliance levels, and source locations |
| `compliance/` | Go test files that exercise a SCIM server against the requirements |
| `scim/` | Minimal SCIM HTTP client used by the tests |
| `testserver/` | In-memory SCIM server for running the suite without an external target |

## Running the tests

### Against the built-in test server

```sh
nix develop --command go test ./compliance/ -v -count=1
```

### Against an external SCIM server

```sh
nix develop --command go test ./compliance/ -v -count=1 \
  -scim.url=https://scim.example.com \
  -scim.token=YOUR_BEARER_TOKEN
```

Basic auth is also supported:

```sh
nix develop --command go test ./compliance/ -v -count=1 \
  -scim.url=https://scim.example.com \
  -scim.user=admin -scim.pass=secret
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-scim.url` | *(empty, uses built-in server)* | SCIM base URL to test against |
| `-scim.token` | | Bearer token for authentication |
| `-scim.user` | | Basic auth username |
| `-scim.pass` | | Basic auth password |
| `-scim.force` | | Comma-separated features to force-enable (e.g. `filter,patch`) |
| `-scim.report` | `compliance-report.txt` | Path for the compliance report |

### Features

The suite discovers supported features via `/ServiceProviderConfig` and
`/ResourceTypes`. Tests for optional features (filter, patch, bulk, sort, etag,
changePassword) are run in soft mode when the feature is not advertised: failures
are recorded as warnings rather than errors. Use `-scim.force` to treat a feature
as required.

## Compliance report

After each run a plain-text report is written (default: `compliance-report.txt`)
containing:

- Summary counts per RFC 2119 level (MUST, SHOULD, MAY, etc.)
- List of failures, warnings, passed, and untested requirements
- Coverage percentage of testable requirements

## License

[Apache License 2.0](LICENSE)
