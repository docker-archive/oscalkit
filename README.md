# oscalkit

[![CircleCI](https://circleci.com/gh/docker/oscalkit.svg?style=svg)](https://circleci.com/gh/docker/oscalkit) [![codecov](https://codecov.io/gh/docker/oscalkit/branch/master/graph/badge.svg)](https://codecov.io/gh/docker/oscalkit) [![GoDoc](https://godoc.org/github.com/docker/oscalkit?status.svg)](https://godoc.org/github.com/docker/oscalkit)

> In development. Since the OSCAL standard is still under active development, parsing errors may occur if running the included CLI tool against OSCAL documents that are developed against iterations of the schemas that aren't supported. Individual [Releases](https://github.com/docker/oscalkit/releases) of `oscalkit` will indicate in the notes which commits in the usnistgov/OSCAL repo against which the tool has been tested.

Barebones Go SDK for the [Open Security Controls Assessment Language (OSCAL)](https://csrc.nist.gov/Projects/Open-Security-Controls-Assessment-Language) which is in development by the [National Institute of Standards and Technology (NIST)](https://www.nist.gov/). A CLI tool is also included for processing OSCAL documents, converting between OSCAL-formatted XML, JSON and YAML and for converting from [OpenControl](http://opencontrol.cfapps.io/) projects in to OSCAL. The tool also supports Go source code generation from OSCAL formatted catalog and profile artifacts.

Documentation for the OSCAL standard can be found at https://pages.nist.gov/OSCAL.

## Supported OSCAL Components

The following OSCAL components are currently supported:

|Component|Schemas|
|---------|-------|
|[Catalog](https://pages.nist.gov/OSCAL/concepts/#oscal-catalogs)|[XSD](https://github.com/usnistgov/OSCAL/blob/master/schema/xml/oscal-catalog-schema.xsd) \| [JSON schema](https://github.com/usnistgov/OSCAL/blob/master/schema/json/oscal-catalog-schema.json) \| [metaschema](https://github.com/usnistgov/OSCAL/blob/master/schema/metaschema/oscal-catalog-metaschema.xml)|
|[Profile](https://pages.nist.gov/OSCAL/concepts/#oscal-profiles)|[XSD](https://github.com/usnistgov/OSCAL/blob/master/schema/xml/oscal-profile-schema.xsd) \| [JSON schema](https://github.com/usnistgov/OSCAL/blob/master/schema/json/oscal-profile-schema.json) \| [metaschema](https://github.com/usnistgov/OSCAL/blob/master/schema/metaschema/oscal-profile-metaschema.xml)|
|Implementation (WIP)|Currently based on a combination of the model being developed in [usnistgov/OSCAL#216](https://github.com/usnistgov/OSCAL/issues/216) and the component definition prototype in [this Gist](https://gist.github.com/anweiss/8afd321b6bf2a9d4e1679657a1b8f2fe)|

## Installing

You can download the appropriate `oscalkit` command-line utility for your system from the [GitHub Releases](https://github.com/docker/oscalkit/releases) page. You can move it to an appropriate directory listed in your `$PATH` environment variable. A Homebrew recipe is also available for macOS along with a [Docker image](https://hub.docker.com/r/docker/oscalkit/) which has been published to Docker Hub.

### Homebrew

    $ brew tap docker/homebrew-oscalkit
    $ brew install oscalkit

### Docker

> Running the `oscalkit` Docker container requires either bind-mounting the directory containing your source files or passing file contents in to the command via stdin.

    $ docker pull docker/oscalkit:0.2.0
    $ docker run -it --rm -v $PWD:/data -w /data docker/oscalkit:0.2.0 convert oscal-core.xml

via stdin:

    $ docker run -it --rm docker/oscalkit:0.2.0 convert < oscal-core.xml

## Usage

```
NAME:
   oscalkit - OSCAL toolkit

USAGE:
   oscalkit [global options] command [command options] [arguments...]

VERSION:
   0.2.0


COMMANDS:
     convert         convert between one or more OSCAL file formats and from OpenControl format
     validate        validate files against OSCAL XML and JSON schemas
     sign            sign OSCAL JSON artifacts
     generate        generates go code against provided profile
     implementation  generates go code for implementation against provided profile and excel sheet
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d    enable debug command output
   --help, -h     show help
   --version, -v  print the version
```

### Convert between XML and JSON

`oscalkit` can be used to convert one or more source files between OSCAL-formatted XML and JSON.

```
NAME:
   oscalkit convert oscal - convert between one or more OSCAL file formats

USAGE:
   oscalkit convert oscal [command options] [source-files...]

DESCRIPTION:
   Convert between OSCAL-formatted XML and JSON files. The command accepts
   one or more source file paths and can also be used with source file contents
   piped/redirected from STDIN.

OPTIONS:
   --output-path value, -o value  Output path for converted file(s). Defaults to current working directory
   --output-file value, -f value  File name for converted output from STDIN. Defaults to "stdin.<json|xml|yaml>"
   --yaml                         If source file format is XML or JSON, also generate equivalent YAML output
```

#### Examples

Convert OSCAL-formatted NIST 800-53 declarations from XML to JSON:

    $ oscalkit convert oscal SP800-53-declarations.xml

Convert OSCAL-formatted NIST 800-53 declarations from XML to JSON via STDIN (note the use of "-"):

    $ cat SP800-53-declarations.xml | oscalkit convert oscal -

### Signing OSCAL JSON with JWS

`oscalkit` can be used to sign OSCAL-formatted JSON artifacts using JSON Web Signature (JWS)

```
NAME:
   oscalkit sign - sign OSCAL JSON artifacts

USAGE:
   oscalkit sign [command options] [files...]

OPTIONS:
   --key value, -k value  private key file for signing. Must be in PEM or DER formats. Supports RSA/EC keys and X.509 certificats with embedded RSA/EC keys
   --alg value, -a value  algorithm for signing. Supports RSASSA-PKCS#1v1.5, RSASSA-PSS, HMAC, ECDSA and Ed25519
```

The following signing algorithms are supported:

 Signing / MAC              | Algorithm identifier(s)
 :------------------------- | :------------------------------
 RSASSA-PKCS#1v1.5          | RS256, RS384, RS512
 RSASSA-PSS                 | PS256, PS384, PS512
 HMAC                       | HS256, HS384, HS512
 ECDSA                      | ES256, ES384, ES512
 Ed25519                    | EdDSA

#### Examples

Sign OSCAL-formatted JSON using a PEM-encoded private key file and the PS256 signing algorithm:

    $ oscalkit sign --key jws-example-key.pem --alg PS256 NIST_SP-800-53_rev4_catalog.json

### Convert from OpenControl project to OSCAL [Experimental]

> This feature has been temporarily disabled pending https://github.com/usnistgov/OSCAL/issues/216 and https://github.com/usnistgov/OSCAL/issues/215

`oscalkit` also supports converting OpenControl projects to OSCAL-formatted JSON. You will need both the path to the `opencontrol.yaml` file and the `opencontrols/` directory which is created when you run a `compliance-masonry get` command.

```
NAME:
   oscalkit convert opencontrol - convert from OpenControl format to OSCAL "implementation" format

USAGE:
   oscalkit convert opencontrol [command options] [opencontrol.yaml-filepath] [opencontrols-dir-path]

DESCRIPTION:
   Convert OpenControl-formatted "component" and "OpenControl" YAML into
   OSCAL-formatted "implementation" layer JSON

OPTIONS:
   --yaml, -y  Generate YAML in addition to JSON
   --xml, -x   Generate XML in addition to JSON
```

### Examples

Convert OpenControl project to OSCAL-formatted JSON:

    $ oscalkit convert opencontrol ./opencontrol.yaml ./opencontrols/

### Validate against XML and JSON schemas

The tool supports validation of OSCAL-formatted XML and JSON files against the corresponding OSCAL XML schemas (.xsd) and JSON schemas. XML schema validation requires the `xmllint` tool on the local machine (included with macOS and Linux. Windows installation instructions [here](https://stackoverflow.com/a/21227833))

```
NAME:
   oscalkit validate - validate files against OSCAL XML and JSON schemas

USAGE:
   oscalkit validate [command options] [files...]

DESCRIPTION:
   Validate OSCAL-formatted XML files against a specific XML schema (.xsd)
   or OSCAL-formatted JSON files against a specific JSON schema

OPTIONS:
   --schema value, -s value  schema file to validate against
```

#### Examples

Validate FedRAMP profile in OSCAL-formatted JSON against the corresponding JSON schema

    $ oscalkit validate -s oscal-core.json fedramp-annotated-wrt-SP800-53catalog.json

## Developing

`oscalkit` is developed with [Go](https://golang.org/) (1.11+). If you have Docker installed, the included `Makefile` can be used to run unit tests and compile the application for Linux, macOS and Windows. Otherwise, the native Go toolchain can be used.

### Dependency management

Dependencies are managed with [Go 1.11 Modules](https://github.com/golang/go/wiki/Modules). The `vendor/` folder containing the dependencies is checked in with the source for backwards compatibility with previous versions of Go. When using Go 1.11 with `GO111MODULE=on`, you can verify the dependencies as follows:

    $ go mod verify

### Compile

You can use the included `Makefile` to generate binaries for your OS as follows (requires [Docker](https://docs.docker.com/engine/installation/)):

Compile for Linux:

    $ GOOS=linux GOARCH=amd64 make

Compile for macOS:

    $ GOOS=darwin GOARCH=amd64 make

Compile for Windows:

    $ GOOS=windows GOARCH=amd64 make

### Website and documentation

Both the website and corresponding documentation are being developed in `docs/`. The content is developed using the [Hugo](https://gohugo.io/) framework. The static content is generated and published in `docs/public`, which is a separate Git worktree that is tied to the [`gh-pages`](https://github.com/docker/oscalkit/tree/gh-pages) branch and publicly accessible via https://docker.github.io/oscalkit.

The GoDoc for the SDK can be found [here](https://godoc.org/github.com/docker/oscalkit).

### Releasing

The [GoReleaser](https://goreleaser.com/) tool is used to publish `oscalkit` to GitHub Releases. The following release artifacts are currently supported:

- OSX binary
- Linux binary
- Windows binary
- Docker Image
- Homebrew recipe
