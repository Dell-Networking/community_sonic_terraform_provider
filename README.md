# Terraform Provider for Community SONiC


[![Contributions welcome](https://img.shields.io/badge/contributions-welcome-orange.svg)](#-developing-the-provider)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Dell-Networking/community_sonic_terraform_provider/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/Dell-Networking/community_sonic_terraform_provider)](https://github.com/Dell-Networking/community_sonic_terraform_provider/issues)

Built and maintained by [Ben Goldstone](https://github.com/benjamingoldstone/) and [Contributors](https://github.com/Dell-Networking/community_sonic_terraform_provider/graphs/contributors)

------------------


- [Requirements](#-requirements)
- [Building the Provider](#-building-the-provider)
- [How to Contribute](#-how-to-contribute)
- [Adding Dependencies](#-adding-dependencies)
- [Using the Provider](#-using-the-provider)
- [Developing the Provider](#-developing-the-provider)

Welcome to the Terraform Provider for Community SONiC!

You can find the provider on the [Hashicorp Terraform Registry](https://registry.terraform.io/dell/community-sonic)

## 📋 Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## 🚀 Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install .
```

Make sure you have overridden the community-sonic provider with a dev-override block in your `~/.terraformrc` file:

```shell
  dev_overrides {
      "registry.terraform.io/dell/community-sonic" = "/home/[user]/go/bin"
  }
```

Ensure the above referenced directory matches your GOBIN go environment setup (recent go versions will default to `$GOPATH/bin`). You can view your go environment with `go env`

You can check the terraform provider by entering the `examples/provider-install-verification` directly and running `terraform plan`. 

## 📋 Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## 📋 Using the Provider

This section is under construction..

## 👏 Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install .`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources and may cost money to run.

```shell
make testacc
```
