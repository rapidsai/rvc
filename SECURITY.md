# Security Policy

`rvc` (RAPIDS Version Converter) is a small Go utility that converts between
RAPIDS CalVer versions (e.g. `25.10`) and the corresponding ucxx / former
UCX-Py SemVer versions (e.g. `0.46`). Unlike the other repositories in the
RAPIDS organization, `rvc` is published as both a CLI binary *and* a live
public HTTP service at <https://version.gpuci.io>, deployed as an AWS Lambda
behind API Gateway. The bulk of this SECURITY.md is therefore service-
shaped rather than library-shaped.

## Reporting a Vulnerability

Please report security vulnerabilities privately through one of the channels
below. **Do not open a public GitHub issue, PR, or discussion** for a
suspected vulnerability.

1. **NVIDIA Vulnerability Disclosure Program (preferred)**
   <https://www.nvidia.com/en-us/security/>
   Submit through the NVIDIA PSIRT web form. This is the fastest path to
   triage and tracking.

2. **Email NVIDIA PSIRT**
   psirt@nvidia.com — encrypt sensitive reports with the
   [NVIDIA PSIRT PGP key](https://www.nvidia.com/en-us/security/pgp-key).

3. **GitHub Private Vulnerability Reporting**
   Use the **Security** tab on this repository → *Report a vulnerability*.

Please include, where possible:

- Affected component (the live service at `version.gpuci.io`, the CLI
  binary, the Terraform deployment, or a specific Go package)
- Whether the issue concerns availability, integrity (incorrect
  conversions), confidentiality (information leak from logs / errors),
  authorization (bypassing the public-read intent), or supply chain
- For the live service: the URL, the inputs, and any observed response
- For the CLI: the binary version (`rvc --help` does not currently print
  one; provide the GitHub release tag or commit you built from)
- Any relevant CWE / CVE identifiers

NVIDIA PSIRT will acknowledge receipt and coordinate triage, fix
development, and coordinated disclosure. More on NVIDIA's response
process: <https://www.nvidia.com/en-us/security/psirt-policies/>.

## Security Architecture & Context

**Classification:** Service + CLI. A stateless Go function packaged two
ways:
- As an AWS Lambda fronted by API Gateway (the public `version.gpuci.io`
  endpoint).
- As a single-binary CLI distributed via GitHub Releases.

**Primary security responsibility:** Parse a small, regex-validated version
string from a URL path parameter or CLI flag, perform a closed-form
arithmetic conversion, and return the result. The function performs no
I/O, holds no state, and contacts no other services.

**Components and trust boundaries:**

- **`pkg/rvc/`** — the conversion logic (`GetUcxPyFromRapids`,
  `GetRapidsFromUcxPy`). Inputs are validated against anchored regular
  expressions:
  - RAPIDS: `^v?[0-9]{2}\.[0-9]{2}(\.[0-9]+)?$`
  - UCX-Py: `^v?0*\.[0-9]{1,2}(\.[0-9]+)?$`
- **`pkg/version/`** — build-time version metadata.
- **`cmd/rvc_cli/`** — CLI entrypoint (`-rapids` / `-ucx-py` flags).
- **`cmd/rvc_serverless/`** — AWS Lambda handler, dispatched from API
  Gateway path parameters.
- **`terraform/`** — AWS deployment topology:
  - `aws_api_gateway_rest_api "rvc"` — public REGIONAL REST API.
  - `GET /rapids/{version}` and `GET /ucx-py/{version}` with
    `authorization = "NONE"` (intentionally public).
  - `aws_iam_role "lambda_role"` granted only CloudWatch Logs writes
    (`logs:CreateLogGroup`, `logs:CreateLogStream`, `logs:PutLogEvents`).
- **Deploy pipeline.** The `deploy.yaml` GitHub Actions workflow deploys
  the Lambda via the Serverless framework on every merge to `main`.

**Out of scope for this policy:** vulnerabilities in AWS Lambda, API
Gateway, the Go runtime itself, the Serverless framework, or in the Go
modules `rvc` depends on (`aws-lambda-go`, `stretchr/testify`,
`gopkg.in/yaml.v3`). Report those to their respective projects.
Vulnerabilities in *how* `rvc` integrates with them — input validation,
error handling, IAM scope, deployment safety — are in scope.

## Threat Model

The threats below are concrete to `rvc`'s service and CLI shape. Two
relevant audit findings exist in the
[RAPIDS Security Audit](https://github.com/orgs/rapidsai/projects/207):
mutable workflow refs and `actions/checkout` token persistence (both
remediated).

1. **Unauthenticated public endpoint abuse / DoS.**
   API Gateway exposes `/rapids/{version}` and `/ucx-py/{version}` with
   `authorization = "NONE"`. This is intentional — the service is
   designed to be queried by build scripts and CI from anywhere. The
   consequences are scope-limited (no backend storage, minimal IAM,
   logs-only) but the service is exposed to general-internet request
   floods, log-volume amplification, and AWS-bill amplification by
   sustained traffic. There is no rate-limiting configured in the
   Terraform; mitigation depends on API Gateway throttling and Lambda
   reserved concurrency, both of which need to be set deliberately.

2. **Caller trust of response strings.**
   The README documents a pattern like
   `RAPIDS_VER=$(curl -sL https://version.gpuci.io/ucx-py/0.46)` —
   callers interpolate the response into shell variables and other
   commands. Because the input is regex-validated and the output is
   constructed via `fmt.Sprintf("%d.%02d", year, month)` from integer
   conversions, the response is constrained to a small alphabet
   (`[0-9.]`). Any future code path that returns un-sanitized error
   text or echoes input directly would expose downstream callers to
   shell-injection through their command interpolation. Preserving the
   "outputs are integers formatted by `%d`" property is a
   security-relevant invariant.

3. **Outdated Go module pins.**
   `go.mod` declares `go 1.17` and pins `aws-lambda-go v1.27.0`,
   `stretchr/testify v1.7.0`, and `yaml.v3 v3.0.0-20200615113413-...`
   These are all multiple years old and have had CVEs published in
   between. `gopkg.in/yaml.v3` specifically had CVE-2022-28948 against
   versions before `v3.0.0-20220521103104-...`. `rvc` does not parse
   YAML at request time, but the dependency travels with the binary
   and inflates the CVE-scan signal for consumers.

4. **CLI binary download without integrity verification.**
   The documented install is
   `wget https://github.com/rapidsai/rvc/releases/latest/download/rvc`
   — a *mutable* reference with no checksum or signature check. A
   compromised release asset or a hijacked maintainer account
   substitutes the binary on the next install.

5. **AWS deploy pipeline.**
   The Serverless-framework deploy runs from a GitHub Actions workflow
   on `main` merges. Compromise of the workflow's AWS deploy
   credentials, or successful `${{ }}` template injection in a deploy
   step, would let an attacker substitute the Lambda's code. The
   broader RAPIDS audit produced fixes for mutable workflow refs in
   this repository; preserving SHA-pinned action references on new
   contributions is ongoing.

6. **`actions/checkout` token persistence.**
   `actions/checkout` defaults to persisting `GITHUB_TOKEN` in the
   workspace's `.git/config`. The audit remediated this in `rvc`'s
   workflows; new workflow contributions should keep
   `persist-credentials: false` where appropriate.

7. **Log content leakage.**
   The Lambda role can write to CloudWatch Logs. If error paths ever
   log raw request URLs or unsanitized inputs, the logs accumulate
   data shaped by unauthenticated requesters. The current code path
   does not appear to do this, but it is worth checking on every
   change to the handler or to package error formatting.

## Critical Security Assumptions

The following are assumed of operators and consumers of `rvc`.

- **The service is public by design.**
  `authorization = "NONE"` is intentional. Operators relying on
  abuse-protection should configure API Gateway request throttling
  and Lambda reserved concurrency rather than expect access control
  inside the function.

- **The response alphabet stays narrow.**
  Today the handler returns either a `%d.%02d`-formatted integer
  string or an HTTP error from API Gateway. Callers that interpolate
  the response into shell commands rely on this. Future changes that
  expand the response shape (returning JSON with caller-controlled
  fields, echoing input in error text) should preserve the small,
  shell-safe output alphabet — or migrate callers to a parser that
  no longer trusts the wire format directly.

- **CLI binary integrity is verified out-of-band, or pinned.**
  Until the release flow signs binaries or publishes a checksum
  manifest, the `releases/latest/download/rvc` install rests on
  TLS to GitHub. Operators should either pin to a specific release
  tag and capture a checksum locally, or rebuild from source.

- **Go module pins are reviewed periodically.**
  `go.mod` will drift further from upstream over time; a periodic
  bump is required to pick up Go and dependency security fixes
  even when there are no functional changes.

- **The Lambda IAM role stays minimal.**
  The current role grants only CloudWatch Logs writes. Any future
  expansion (network access, S3 reads, parameter-store lookups)
  should be reviewed against this baseline and justified.

- **Deploy credentials are scoped to this stack.**
  The Serverless-framework deploy uses AWS credentials configured in
  GitHub Actions secrets. Those credentials should be scoped to the
  `rvc-*` resources only, not be broader account-level keys.

- **No persistent state.**
  `rvc` does not store user inputs, request bodies, or derived
  artifacts. Operators should preserve this property — adding any
  storage backend changes the threat model significantly.

## Supported Versions

The live service at `version.gpuci.io` runs whatever is deployed from
`main`. The CLI follows GitHub Releases; older release tags are not
back-ported with security fixes. Pull a recent tag or rebuild from
`main` to receive updates.

## Dependency Security

`rvc` depends on the standard Go runtime plus `aws-lambda-go`,
`stretchr/testify` (tests), and indirect modules including
`gopkg.in/yaml.v3`. Upstream CVEs in any of those — particularly in
the Go standard library and `aws-lambda-go` — should trigger a
module bump and a redeploy. The CLI binary's stdlib version is
fixed at build time; the Lambda's runtime version is provided by
AWS Lambda's managed Go runtime.
