# Release

TF.tf releases are prepared for the Terraform Registry with GoReleaser.

## Requirements

- GitHub repository name must remain `terraform-provider-tensorflow`.
- Release tags must be valid SemVer tags with a leading `v`, such as `v0.1.0`.
- Release artifacts must include zip archives, SHA256 checksums, checksum
  signature, and the Terraform Registry manifest.
- The Terraform Registry signing key must match the GPG key used by the release
  workflow.

HashiCorp's publishing docs describe the public registry requirements:

- Provider repositories must match `terraform-provider-{NAME}`.
- Plugin Framework providers should set protocol version `6.0` in the registry
  manifest.
- GitHub release tags must use SemVer with a leading `v`.
- Avoid replacing an already released version; release a new version instead.

## GitHub Secrets

The release workflow expects:

- `GPG_PRIVATE_KEY`
- `PASSPHRASE`
- `GPG_FINGERPRINT`

`GITHUB_TOKEN` is provided by GitHub Actions.

## Dry Run

Run a local build without publishing:

```sh
goreleaser release --snapshot --clean
```

## Release

Create and push a tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

The GitHub Actions release workflow creates a draft GitHub release. Review the
assets before publishing the draft.

## Registry Docs

Registry-facing docs live in `docs/index.md` and `docs/data-sources`. The
broader project usage guide stays in `docs/usage`.
