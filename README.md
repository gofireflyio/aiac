# ![AIAC](logo-header.svg#gh-light-mode-only) ![AIAC](logo-header-inverted.svg#gh-dark-mode-only)

Artificial Intelligence
Infrastructure-as-Code
Generator.

<kbd>[<img src="demo.gif" style="width: 100%; border: 1px solid silver;" border="1" alt="demo">](demo.gif)</kbd>

<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Use Cases and Example Prompts](#use-cases-and-example-prompts)
    * [Generate IaC](#generate-iac)
    * [Generate Configuration Files](#generate-configuration-files)
    * [Generate CI/CD Pipelines](#generate-cicd-pipelines)
    * [Generate Policy as Code](#generate-policy-as-code)
    * [Generate Utilities](#generate-utilities)
    * [Command Line Builder](#command-line-builder)
    * [Query Builder](#query-builder)
* [Instructions](#instructions)
    * [Installation](#installation)
    * [Configuration](#configuration)
    * [Usage](#usage)
        * [Command Line](#command-line)
            * [Listing Models](#listing-models)
            * [Generating Code](#generating-code)
        * [Via Docker](#via-docker)
        * [As a Library](#as-a-library)
    * [Upgrading from v4 to v5](#upgrading-from-v4-to-v5)
        * [Changes in Configuration](#changes-in-configuration)
        * [Changes in CLI Invokation](#changes-in-cli-invokation)
        * [Changes in Model Usage and Support](#changes-in-model-usage-and-support)
        * [Other Changes](#other-changes)
* [Example Output](#example-output)
* [Troubleshooting](#troubleshooting)
* [License](#license)

<!-- vim-markdown-toc -->

## Description

`aiac` is a library and command line tool to generate IaC (Infrastructure as Code)
templates, configurations, utilities, queries and more via [LLM](https://en.wikipedia.org/wiki/Large_language_model) providers such
as [OpenAI](https://openai.com/), [Amazon Bedrock](https://aws.amazon.com/bedrock/) and [Ollama](https://ollama.ai/).

The CLI allows you to ask a model to generate templates for different scenarios
(e.g. "get terraform for AWS EC2"). It composes an appropriate request to the
selected provider, and stores the resulting code to a file, and/or prints it to
standard output.

Users can define multiple "backends" targeting different LLM providers and
environments using a simple configuration file.

## Use Cases and Example Prompts

### Generate IaC

- `aiac get terraform for a highly available eks`
- `aiac get pulumi golang for an s3 with sns notification`
- `aiac get cloudformation for a neptundb`

### Generate Configuration Files

- `aiac get dockerfile for a secured nginx`
- `aiac get k8s manifest for a mongodb deployment`

### Generate CI/CD Pipelines

- `aiac get jenkins pipeline for building nodejs`
- `aiac get github action that plans and applies terraform and sends a slack notification`

### Generate Policy as Code

- `aiac get opa policy that enforces readiness probe at k8s deployments`

### Generate Utilities

- `aiac get python code that scans all open ports in my network`
- `aiac get bash script that kills all active terminal sessions`

### Command Line Builder

- `aiac get kubectl that gets ExternalIPs of all nodes`
- `aiac get awscli that lists instances with public IP address and Name`

### Query Builder

- `aiac get mongo query that aggregates all documents by created date`
- `aiac get elastic query that applies a condition on a value greater than some value in aggregation`
- `aiac get sql query that counts the appearances of each row in one table in another table based on an id column`

## Instructions

Before installing/running `aiac`, you may need to configure your LLM providers
or collect some information.

For **OpenAI**, you will need an API key in order for `aiac` to work. Refer to
[OpenAI's pricing model](https://openai.com/pricing?trk=public_post-text) for more information. If you're not using the API hosted
by OpenAI (for example, you may be using Azure OpenAI), you will also need to
provide the API URL endpoint.

For **Amazon Bedrock**, you will need an AWS account with Bedrock enabled, and
access to relevant models. Refer to the [Bedrock documentation](https://docs.aws.amazon.com/bedrock/latest/userguide/what-is-bedrock.html)
for more information.

For **Ollama**, you only need the URL to the local Ollama API server, including
the /api path prefix. This defaults to http://localhost:11434/api. Ollama does
not provide an authentication mechanism, but one may be in place in case of a
proxy server being used. This scenario is not currently supported by `aiac`.

### Installation

Via `brew`:

    brew tap gofireflyio/aiac https://github.com/gofireflyio/aiac
    brew install aiac

Using `docker`:

    docker pull ghcr.io/gofireflyio/aiac

Using `go install`:

    go install github.com/gofireflyio/aiac/v5@latest

Alternatively, clone the repository and build from source:

    git clone https://github.com/gofireflyio/aiac.git
    go build

### Configuration

`aiac` is configured via a TOML configuration file. Unless a specific path is
provided, `aiac` looks for a configuration file in the user's [XDG_CONFIG_HOME](https://en.wikipedia.org/wiki/Freedesktop.org#User_directories)
directory, specifically `${XDG_CONFIG_HOME}/aiac/aiac.toml`. On Unix-like
operating systems, this will default to "~/.config/aiac/aiac.toml". If you want
to use a different path, provide the `--config` or `-c` flag with the file's path.

The configuration file defines one or more named backends. Each backend has a
type identifying the LLM provider (e.g. "openai", "bedrock", "ollama"), and
various settings relevant to that provider. Multiple backends of the same LLM
provider can be configured, for example for "staging" and "production"
environments.

Here's an example configuration file:

```toml
default_backend = "openai"   # Default backend when one is not selected

[backends.official_openai]
type = "openai"
api_key = "API KEY"
default_model = "gpt-4o"     # Default model to use for this backend

[backends.azure_openai]
type = "openai"
url = "https://tenant.openai.azure.com/openai/deployments/test"
api_key = "API KEY"
api_version = "2023-05-15"   # Optional

[backends.aws_staging]
type = "bedrock"
aws_profile = "staging"
aws_region = "eu-west-2"

[backends.aws_prod]
type = "bedrock"
aws_profile = "production"
aws_region = "us-east-1"
default_model = "amazon.titan-text-express-v1"

[backends.localhost]
type = "ollama"
url = "http://localhost:11434/api"     # This is the default
```

### Usage

Once a configuration file is created, you can start generating code and you only
need to refer to the name of the backend. You can use `aiac` from the command
line, or as a Go library.

#### Command Line

##### Listing Models

Before starting to generate code, you can list all models available in a
backend:

    aiac -b aws_prof --list-models

This will return a list of all available models. Note that depending on the LLM
provider, this may list models that aren't accessible or enabled for the
specific account.

##### Generating Code

By default, aiac prints the extracted code to standard output and opens an
interactive shell that allows conversing with the model, retrying requests,
saving output to files, copying code to clipboard, and more:

    aiac terraform for AWS EC2

This will use the default backend in the configuration file and the default
model for that backend, assuming they are indeed defined. To use a specific
backend, provide the `--backend` or `-b` flag:

    aiac -b aws_prof terraform for AWS EC2

To use a specific model, provide the `--model` or `-m` flag:

    aiac -m gpt-4-turbo terraform for AWS EC2

You can ask `aiac` to save the resulting code to a specific file:

    aiac terraform for eks --output-file=eks.tf

You can use a flag to save the full Markdown output as well:

    aiac terraform for eks --output-file=eks.tf --readme-file=eks.md

If you prefer aiac to print the full Markdown output to standard output rather
than the extracted code, use the `-f` or `--full` flag:

    aiac terraform for eks -f

You can use aiac in non-interactive mode, simply printing the generated code
to standard output, and optionally saving it to files with the above flags,
by providing the `-q` or `--quiet` flag:

    aiac terraform for eks -q

In quiet mode, you can also send the resulting code to the clipboard by
providing the `--clipboard` flag:

    aiac terraform for eks -q --clipboard

Note that aiac will not exit in this case until the contents of the clipboard
changes. This is due to the mechanics of the clipboard.

#### Via Docker

All the same instructions apply, except you execute a `docker` image:

    docker run \
        -it \
        -v ~/.config/aiac/aiac.toml:~/.config/aiac/aiac.toml \
        ghcr.io/gofireflyio/aiac terraform for ec2

#### As a Library

You can use `aiac` as a Go library:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/gofireflyio/aiac/v5/libaiac"
)

func main() {
    aiac, err := libaiac.New() // Will load default configuration path.
                               // You can also do libaiac.New("/path/to/aiac.toml")
    if err != nil {
        log.Fatalf("Failed creating aiac object: %s", err)
    }

    ctx := context.TODO()

    models, err := aiac.ListModels(ctx, "backend name")
    if err != nil {
        log.Fatalf("Failed listing models: %s", err)
    }

    chat, err := aiac.Chat(ctx, "backend name", "model name")
    if err != nil {
        log.Fatalf("Failed starting chat: %s", err)
    }

    res, err = chat.Send(ctx, "generate terraform for eks")
    res, err = chat.Send(ctx, "region must be eu-central-1")
}
```

### Upgrading from v4 to v5

Version 5.0.0 introduced a significant change to the `aiac` API in both the
command line and library forms, as per feedback from the community.

#### Changes in Configuration

Before v5, there was no concept of a configuration file or named backends. Users
had to provide all the information necessary to contact a specific LLM provider
via command line flags or environment variables, and the library allowed
creating a "client" object that could only talk with one LLM provider.

Backends are now configured only via the configuration file. Refer to the
[Configuration](#configuration) section for instructions. Provider-specific flags such as
`--api-key`, `--aws-profile`, etc. (and their respective environment variables,
if any) are no longer accepted.

Since v5, backends are also named. Previously, the `--backend` and `-b` flags
referred to the name of the LLM provider (e.g. "openai", "bedrock", "ollama").
Now they refer to whatever name you've defined in the configuration file:

```toml
[backends.my_local_llm]
type = "ollama"
url = "http://localhost:11434/api"
```

Here we configure an Ollama backend named "my_local_llm". When you want to
generate code with this backend, you will use `-b my_local_llm` rather than
`-b ollama`, as multiple backends may exist for the same LLM provider.

#### Changes in CLI Invokation

Before v5, the command line was split into three subcommands: `get`,
`list-models` and `version`. Due to this hierarchical nature of the CLI, flags may
not have been accepted if they were provided in the "wrong location". For
example, the `--model` flag had to be provided after the word "get", otherwise
it would not be accepted. In v5, there are no subcommands, so the position of
the flags no longer matters.

The `list-models` subcommand is replaced with the flag `--list-models`, and the
`version` subcommand is replaced with the flag `--version`.

Before v5:

    aiac -b ollama list-models

Since v5:

    aiac -b my_local_llm --list-models

In earlier versions, the word "get" was actually a subcommand and not truly part
of the prompt sent to the LLM provider. Since v5, there is no "get" subcommand,
so you no longer need to add this word to your prompts.

Before v5:

    aiac get terraform for S3 bucket

Since v5:

    aiac terraform for S3 bucket

That said, adding either the word "get" or "generate" will not hurt, as v5 will
simply remove it if provided.

#### Changes in Model Usage and Support

Before v5, the models for each LLM provider were hardcoded in each backend
implementation, and each provider had a hardcoded default model. This
significantly limited the usability of the project, and required us to update
`aiac` whenever new models were added or deprecated. On the other hand, we could
provide extra information about each model, such as its context lengths and
type, as we manually extracted them from the provider documentation.

Since v5, `aiac` no longer hardcodes any models, including default ones. It
will not attempt to verify the model you select actually exists. The
`--list-models` flag will now directly contact the chosen backend API to get a
list of supported models. Setting a model when generating code simply sends its
name to the API as-is. Also, instead of hardcoding a default model for each
backend, users can define their own default models in the configuration file:

```toml
[backends.my_local_llm]
type = "ollama"
url = "http://localhost:11434/api"
default_model = "mistral:latest"
```

Before v5, `aiac` supported both completion models and chat models. Since v5,
it only supports chat models. Since none of the LLM provider APIs actually
note whether a model is a completion model or a chat model (or even an image
or video model), the `--list-models` flag may list models which are not actually
usable, and attempting to use them will result in an error being returned from
the provider API. The reason we've decided to drop support for completion models
was that they require setting a maximum amount of tokens for the API to
generate (at least in OpenAI), which we can no longer do without knowing the
context length. Chat models are not only a lot more useful, but they do not have
this limitation.

#### Other Changes

Most LLM provider APIs, when returning a response to a prompt, will include a
"reason" for why the response ended where it did. Generally, the response should
end because the model finished generating a response, but sometimes the response
may be truncated due to the model's context length or the user's token
utilization. When the response did not "stop" because it finished generation,
the response is said to be "truncated". Before v5, if the API returned that the
response was truncated, `aiac` returned an error. Since v5, an error is no longer
returned, as it seems that some providers do not return an accurate stop reason.
Instead, the library returns the stop reason as part of its output for users to
decide how to proceed.

## Example Output

Command line prompt:

    aiac dockerfile for nodejs with comments

Output:

```Dockerfile
FROM node:latest

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json ./

RUN npm install
# If you are building your code for production
# RUN npm ci --only=production

# Bundle app source
COPY . .

EXPOSE 8080
CMD [ "node", "index.js" ]
```

## Troubleshooting

Most errors that you are likely to encounter are coming from the LLM provider
API, e.g. OpenAI or Amazon Bedrock. Some common errors you may encounter are:

- "[insufficient_quota] You exceeded your current quota, please check your plan and billing details":
  As described in the [Instructions](#instructions) section, OpenAI is a paid API with a certain
  amount of free credits given. This error means you have exceeded your quota,
  whether free or paid. You will need to top up to continue usage.

- "[tokens] Rate limit reached...":
  The OpenAI API employs rate limiting as [described here](https://platform.openai.com/docs/guides/rate-limits/request-increase). `aiac` only performs
  individual requests and cannot workaround or prevent these rate limits. If
  you are using `aiac` in programmatically, you will have to implement throttling
  yourself. See [here](https://github.com/openai/openai-cookbook/blob/main/examples/How_to_handle_rate_limits.ipynb) for tips.

## License

This code is published under the terms of the [Apache License 2.0](/LICENSE).
