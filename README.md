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
    * [Usage](#usage)
        * [Command Line](#command-line)
        * [Via Docker](#via-docker)
        * [As a Library](#as-a-library)
* [Example Output](#example-output)
* [Troubleshooting](#troubleshooting)
* [Support Channels](#support-channels)
* [License](#license)

<!-- vim-markdown-toc -->

## Description

`aiac` is a command line tool to generate IaC (Infrastructure as Code) templates, configurations, utilities, queries and more
via [OpenAI](https://openai.com/)'s API. The CLI allows you to ask the model to generate templates
for different scenarios (e.g. "get terraform for AWS EC2"). It will make the
request, and store the resulting code to a file, or simply print it to standard
output. By default, `aiac` uses the same model used by ChatGPT, but allows using
different models.

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

You will need to provide an OpenAI API key in order for `aiac` to work. Refer to
[OpenAI's pricing model](https://openai.com/pricing?trk=public_post-text) for
more information. As of this writing, you get $5 in free credits upon signing up,
but generally speaking, this is a paid API.

### Installation

Via `brew`:

    brew install gofireflyio/aiac/aiac

Using `docker`:

    docker pull ghcr.io/gofireflyio/aiac

Using `go install`:

    go install github.com/gofireflyio/aiac/v3@latest

Alternatively, clone the repository and build from source:

    git clone https://github.com/gofireflyio/aiac.git
    go build

### Usage

1. Create your OpenAI API key [here](https://platform.openai.com/account/api-keys).
2. Click “Create new secret key” and copy it.
3. Provide the API key via the `OPENAI_API_KEY` environment variable or via the `--api-key` command line flag.

#### Command Line

By default, aiac prints the extracted code to standard output and opens an
interactive shell that allows retrying requests, enabling chat mode (for chat
models), saving output to files, copying code to clipboard, and more:

    aiac get terraform for AWS EC2

You can ask it to also store the code to a specific file with a flag:

    aiac get terraform for eks --output-file=eks.tf

You can use a flag to save the full Markdown output as well:

    aiac get terraform for eks --output-file=eks.tf --readme-file=eks.md

If you prefer aiac to print the full Markdown output to standard output rather
than the extracted code, use the `-f` or `--full` flag:

    aiac get terraform for eks -f

You can use aiac in non-interactive mode, simply printing the generated code
to standard output, and optionally saving it to files with the above flags,
by providing the `-q` or `--quiet` flag:

    aiac get terraform for eks -q

In quiet mode, you can also send the resulting code to the clipboard by
providing the `--clipboard` flag:

    aiac get terraform for eks -q --clipboard

Note that aiac will not exit in this case until the contents of the clipboard
changes. This is due to the mechanics of the clipboard.

By default, aiac uses the gpt-3.5-turbo chat model, but other models are
supported, including gpt-4. You can list all supported models:

    aiac list-models

To generate code with a different model, provide the `--model` flag:

    aiac get terraform for eks --model="text-davinci-003"

#### Via Docker

All the same instructions apply, except you execute a `docker` image:

    docker run \
        -it \
        -e OPENAI_API_KEY=[PUT YOUR KEY HERE] \
        ghcr.io/gofireflyio/aiac get terraform for ec2

#### As a Library

You can use aiac as a library:

```go
package main

import (
    "context"
    "os"

    "github.com/gofireflyio/aiac/v3/libaiac"
)

func main() {
    client := libaiac.NewClient(os.Getenv("OPENAI_API_KEY"))
    ctx    := context.TODO()

    // use the model-agnostic wrapper
    res, err := client.GenerateCode(
        ctx,
        libaiac.ModelTextDaVinci3,
        "generate terraform for ec2",
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed generating code: %s\n", err)
        os.Exit(1)
    }

    fmt.Fprintln(os.Stdout, res.Code)

    // use the completion API (for completion-only models)
    res, err = client.Complete(
        ctx,
        libaiac.ModelTextDaVinci3,
        "generate terraform for ec2",
    )

    // converse via a chat model
    chat := client.Chat(libaiac.ModelGPT35Turbo)
    res, err = chat.Send(ctx, "generate terraform for eks")
    res, err = chat.Send(ctx, "region must be eu-central-1")
}
```

## Example Output

Command line prompt:

    aiac get dockerfile for nodejs with comments

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

`aiac` is a command line client to OpenAI's API. Most errors that you are likely
to encounter are coming from this API. Some common errors you may encounter are:

- "[insufficient_quota] You exceeded your current quota, please check your plan and billing details":
  As described in the [Instructions](#instructions) section, OpenAI is a paid API with a certain
  amount of free credits given. This error means you have exceeded your quota,
  whether free or paid. You will need to top up to continue usage.

- "[tokens] Rate limit reached...":
  The OpenAI API employs rate limiting as [described here](https://platform.openai.com/docs/guides/rate-limits/request-increase). `aiac` only performs
  individual requests and cannot workaround or prevent these rate limits. If
  you are using `aiac` in programmatically, you will have to implement throttling
  yourself. See [here](https://github.com/openai/openai-cookbook/blob/main/examples/How_to_handle_rate_limits.ipynb) for tips.

## Support Channels

We have two main channels for supporting AIaC:

1. [Slack community](https://join.slack.com/t/firefly-community/shared_invite/zt-1m0d5c740-EhHAAFV5mhYBNXxcMWJp7g): general user support and engagement.
2. [GitHub Issues](https://github.com/gofireflyio/aiac/issues): bug reports and enhancement requests.

## License

This code is published under the terms of the [Apache License 2.0](/LICENSE).

# Terraform Template for VPC and EC2 Instance

This Terraform template creates a VPC and an EC2 instance in the VPC. The VPC includes an internet gateway and a route table that routes traffic to the internet gateway. The EC2 instance is launched in a specified subnet and uses a specified security group.

To use this template, modify the variables in `vpc/variables.tf` and `ec2/variables.tf` to suit your needs. Then, run the following commands:

```
terraform init
terraform apply
```

After the resources are created, the public IP address of the EC2 instance will be outputted.

Example output:

```
Outputs:

instance_ip = "X.X.X.X"
```
This code was generated by [Firefly](https://app.gofirefly.io)