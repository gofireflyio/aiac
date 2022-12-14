# ![AIAC](logo-header.svg#gh-light-mode-only) ![AIAC](logo-header-inverted.svg#gh-dark-mode-only)

Artificial Intelligence
Infrastructure-as-Code
Generator.

<kbd>[<img src="demo.gif" style="width: 100%; border: 1px solid silver;" border="1" alt="demo">](demo.gif)</kbd>

<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Quick Start](#quick-start)
    * [Usage via OpenAI API](#usage-via-openai-api)
    * [Usage via ChatGPT](#usage-via-chatgpt)
* [Example Prompts](#example-prompts-and-usecases)
* [Example Output](#example-output)
* [Acknowledgements](#acknowledgements)
* [License](#license)

<!-- vim-markdown-toc -->

## Description

`aiac` is a command line tool to generate IaC (Infrastructure as Code) templates
via [OpenAI](https://openai.com/)'s API. The CLI allows you to ask the model to generate templates
for different scenarios (e.g. "get terraform for AWS EC2"). It will make the
request, and store the resulting code to a file, or simply print it to standard
output.

## Quick Start

First, install `aiac`:

    brew install gofireflyio/aiac/aiac

Or using `docker`:

    docker pull ghcr.io/gofireflyio/aiac

Alternatively, clone the repository and build from source:

    git clone https://github.com/gofireflyio/aiac.git
    go build


### Instructions

1. Create your OpenAI API key [here](https://beta.openai.com/account/api-keys).
1. Click “Create new secret key” and copy it.
1. Provide the API key via the `OPENAI_API_KEY` environment variable or via the `--api-key` command line flag.

By default, aiac prints the extracted code to standard output and asks if it should save or re-generate the code 

    aiac get terraform for AWS EC2

To store the resulting code to a file:

    aiac -o "aws_ec2.tf" get terraform for AWS EC2
         
To run using `docker`

    docker run \
    -it \
    -e OPENAI_API_KEY=[PUT YOUR KEY HERE] \
    ghcr.io/gofireflyio/aiac get terraform for ec2

## Example Prompts and Usecases

### Generate IaC
- `aiac get terraform for a highly available eks`
- `aiac get pulumi golang for an s3 with sns notification`
- `aiac get cloudformation for a neptundb`

### Generate Configuration Files
- `aiac get dockerfile for a secured nginx`
- `aiac get k8s manifest for a mongodb deployment`

### Generate CICD Pipelines
- `aiac get jenkins pipelie for building nodejs`
- `aiac get github action which plan and apply terraform and send slack notification`

### Policy as Code
- `aiac get opa policy that enforces readiness probe at k8s deployments`

### Shell Scripts
- `aiac get bash script that is killing all active terminal sessions`

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

## License

This code is published under the terms of the [Apache License 2.0](/LICENSE).
