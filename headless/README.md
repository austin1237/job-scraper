# lol-counter-source-api
A lambda that invokes a headless browser to render a page (including it's javascript) and passes along the rendered html.

## Why is this lambda using a container deployment rather than the standard zip deployment?
[Pupeteer](https://pptr.dev/) requires a chrome/chromium binary which execeeded the standard [lambda size limit](https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html#function-configuration-deployment-and-execution). Using a container image greatly increases the limit and allows for the binary to be deployed. Currently this service also uses [@sparticuz/chromium](https://github.com/Sparticuz/chromium) due to the standard pupeeteer chromium install having permissions issues when running in the deployed aws env.

## Prerequisites
You must have the following installed/configured on your system for this to work correctly<br />
1. [Docker](https://www.docker.com/)
2. [Docker-Compose](https://docs.docker.com/compose/)

## Development Environment
The development environment uses a pinned version of [aws's node 18 image](https://gallery.ecr.aws/lambda/nodejs) to mimic the running lambda. 

```bash
docker-compose up
```

The output is similar to what you would see in cloudwatch logs ex.

```bash
headless-lambda-1  | 18 Aug 2023 09:47:04,515 [INFO] (rapid) exec '/var/runtime/bootstrap' (cwd=/var/task, handler=)
```

The endpoint of the local container is localhost:3000/2015-03-31/functions/function/invocations send a POST request with the following body
```json
{
    "queryStringParameters": {
	"url": "https://www.google.com"
}}
```