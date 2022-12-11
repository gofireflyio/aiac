# aiac

**AI-generated IaC templates via ChatGPT.**

<!-- vim-markdown-toc GFM -->

* [Description](#description)
* [Quick Start](#quick-start)
    * [Usage via OpenAI API](#usage-via-openai-api)
    * [Usage via ChatGPT](#usage-via-chatgpt)
* [Example Prompts](#example-prompts)
* [Example Output](#example-output)
* [Acknowledgements](#acknowledgements)
* [License](#license)

<!-- vim-markdown-toc -->

## Description

`aiac` is a command line tool to generate IaC (Infrastructure as Code) templates
via [OpenAI](https://openai.com/)'s API or via ChatGPT. The CLI allows you to ask the model to generate templates
for different scenarios (e.g. "generate terraform for AWS EC2"). It will make the
request, and store the resulting code to a file, or simply print it to standard
output.

When using ChatGPT, the server returns a Markdown file with code and explanations.
The CLI will extract the code in this case, and optionally store the entire
Markdown of explanations to a separate file.

## Quick Start

First, install `aiac`:

    go generate github.com/gofireflyio/aiac

Alternatively, clone the repository and build from source:

    git clone https://github.com/gofireflyio/aiac.git
    go build

### Usage via OpenAI API

You will need to provide `aiac` with an API key. Create your API key [here](https://beta.openai.com/account/api-keys).
You can either provide the API key via the `--api-key` command line flag, or via
the `OPENAI_API_KEY` environment variable.

By default, `aiac` simply prints the extracted code to standard output

    aiac --api-key=API_KEY generate terraform for AWS EC2

To store the resulting code to a file:

    aiac --api-key=API_KEY \
         --output-file="aws_ec2.tf" \
         get terraform for AWS EC2

### Usage via ChatGPT

There are several caveats to using `aiac` in ChatGPT mode:

- ChatGPT's API is not public, and is likely to change frequently, which
  may break this program. Please inform us via the [issues page](https://github.com/gofireflyio/aiac/issues) if this happens.
- ChatGPT may rate limit your requests, and is prone to answer slowly or not at
  all when under heavy load.
- You will currently have to manually copy a session token from an actual browser
  session in order to authenticate (instructions follow).

You will need to provide `aiac` with a session token. Since ChatGPT doesn't
currently support programmatic usage, you will need to do this via your browser
(this is hopefully temporary until we can implement OpenAI authentication).
To get a token, follow these steps:

1. Login to [ChatGPT](https://chat.openai.com/) via your web browser.
2. Open the Web Developer Tools (usually Ctrl+Shift+I).
3. Go to the "Storage" tab, and move to the list of "Cookies".
4. Find the cookie called "__Secure-next-auth.session-token".
5. Copy its value. This is the session token. You can store it in the
   `CHATGPT_SESSION_TOKEN` environment variable, or provide it via the
   `--session-token` command line flag.

![](/authentication.jpg)

Then run:

    aiac --chat-gpt \
         --session-token=TOKEN \
         --output-file="ec2.tf" \
         --readme-file="README.md" \
         get terraform for AWS EC2

## Example Prompts

The following prompts are known to work:

- get Terraform for AWS EC2
- get Dockerfile for NodeJS with comments
- get GitHub action for deploying Terraform
- get Python code for Pulumi that deploys Azure VPC

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

## Acknowledgements

Thanks to [acheong08/ChatGPT](https://github.com/acheong08/ChatGPT) for helping with reverse engineering the ChatGPI
API.

## License

This code is published under the terms of the [Apache License 2.0](/LICENSE).
