# aiac

**AI-generated IaC templates via ChatGPT.**

<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Caveats](#caveats)
* [Quick Start](#quick-start)
* [Example Prompts](#example-prompts)
* [Example Output](#example-output)
* [Acknowledgements](#acknowledgements)
* [License](#license)

<!-- vim-markdown-toc -->

## Description

`aiac` is a command line tool to generate IaC (Infrastructure as Code) templates
via [OpenAI](https://openai.com/)'s [ChatGPT](https://chat.openai.com/). The CLI allows you to ask ChatGPT to generate templates
for different scenarios (e.g. "generate terraform for AWS EC2"). It then makes
the request on your behalf, accepts the response, and extracts the example code
from it. This extracted code is either printed to standard output, or saved to a
file. Optionally, the CLI will append the complete message (i.e. the full
Markdown response including ChatGPT's explanations) to a separate markdown file.

## Caveats

- ChatGPT's API is not public, and is likely to change frequently, which
  may break this program. Please inform us via the [issues page](https://github.com/gofireflyio/aiac/issues) if this happens.
- ChatGPT may rate limit your requests, and is prone to answer slowly or not at
  all when under heavy load.
- ChatGPT's responses to the same prompt may differ between executions. If you
  are unhappy with the results, try again or modify your prompt.
- You will currently have to manually copy a session token from an actual browser
  session. Directions in the [Quick Start](#quick-start) section.

## Quick Start

First, install `aiac`:

    go get github.com/gofireflyio/aiac

Alternatively, clone the repository and build from source:

    git clone https://github.com/gofireflyio/aiac.git
    go build

You will need to provide `aiac` with a session token. Since ChatGPT doesn't
currently support programmatic usage, you will need to do this via your browser
(this is hopefully temporary until we can implement OpenAI authentication).
To get a token, follow these steps:

1. Login to [ChatGPT](https://chat.openai.com/) via your web browser.
2. Open the Web Developer Tools (usually Ctrl+Shift+I).
3. Go to the "Storage" tab, and move to the list of "Cookies".
4. Find the cookie called "__Secure-next-auth.session-token".
5. Copy its value. This is the session token.

![](/authentication.jpg)

By default, `aiac` simply prints the extracted code to standard output

    aiac --session-token=TOKEN generate terraform for AWS EC2

To store the resulting code to a file, and to append the explanations to a
markdown file, run:

    aiac --session-token=TOKEN \
         --output-file="aws_ec2.tf" \
         --readme-file="README.md" \
         generate terraform for AWS EC2

## Example Prompts

The following prompts are known to work:

- generate Terraform for AWS EC2
- generate Dockerfile for NodeJS
- generate GitHub action for deploying Terraform
- generate Python code for Pulumi that deploys Azure VPC

## Example Output

Command line prompt:

```sh
aiac --session-token=TOKEN generate Dockerfile for NodeJS
```

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

## Acknowledgements

Thanks to [acheong08/ChatGPT](https://github.com/acheong08/ChatGPT) for helping with reverse engineering the ChatGPI
API.

## License

This code is published under the terms of the [Apache License 2.0](/LICENSE).
