# Identeco-go

Go implementation of [identeco](https://github.com/dmsi/identeco). Can be deployed in AWS Lambda or as standalone HTTP service.

Goals:

-   [x] Implement a service which issues JWT tokens
-   [x] Use assymetric JWT-signing method
-   [x] Rotate keys periodically
-   [x] Deployment on `AWS Lambda` / `go1.x` runtime
-   [x] Deployment as standalone HTTP service
-   [ ] CI/CD github actions
-   [x] It is **NOT** designed to run at scale

# Principal design

## AWS Lambda

Here is an example how it is designed for [AWS Lambda](https://github.com/dmsi/identeco#principal-design). The lambda uses `s3` and `DynamoDB` as storage backends for keys and users.

## Standalone HTTP server

Standalone HTTP server is not tied up with AWS services so `MongoDB` is used as storage backend for both keys and users (it is still possible to use `s3` and `DynamoDB`).

# Pre-reqs

-   nodejs (tested on v16.19.1)
-   serverless installed globally (tested on 3.33.0)
-   golang 1.x (tested on go1.20.5)

```sh
npm install -g serverless
```

# Operations

## Deploy serverless in AWS Lambda

The `serverless` framework is used as infrastructure and deployment orchestrator.
The deployment manifests are located under `deployment/awslambda` directory. So all `serverless` commands must be executed from this directory.

> **Note** before you deploy change `provider.profile` to match your desired AWS profile or delete in order to use the default profile.
> Optionally change `provider.region` to reflect region of your choice.

Deploy whole stack (default stage is 'dev')

```bash
npm install
cd deployments/awslambda
serverless deploy
serverless invoke -f rotateKeys
```

> **Note** rotateKeys function is trigerred periodically by CloudWatch events but in order to
> rotate keys the first time it needs to be triggered manually right after the deployment.

Serverless will create AWS `cloudformation` with all the resources specified in `serverless.yml`.
Example output

```bash
$ serverless deploy

Deploying identeco to stage dev (eu-west-1)

✔ Service deployed to stack identeco-dev (58s)

endpoints:
  GET - https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev/.well-known/jwks.json
  POST - https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev/register
  POST - https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev/login
  GET - https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev/refresh
functions:
  getJwks: identeco-dev-getJwks (17 MB)
  register: identeco-dev-register (17 MB)
  login: identeco-dev-login (17 MB)
  refresh: identeco-dev-refresh (17 MB)
  rotateKeys: identeco-dev-rotateKeys (17 MB)

Monitor all your API routes with Serverless Console: run "serverless --console"
```

### Environment variables

All environment variables are defined in `serverless.yml` manifest in `provider.environment` section. Those variables are accesible in each lambda function.

## Run standalone HTTP server locally

**TBD**

## Deploy as standalone HTTP server in AWS EC2

**TBD**

### Environment variables

**TBD**

### Run python test

> **Note** `python3.9+` is required. The `venv` module can be used in order to localize the dependencies.

The following snippet can be used in order to run the tests in `bash` environment

```bash
cd ./test

# Setup pyton venv and activate it
python -m venv myenv
source myenv/bin/activate

# Setup python dependencies
pip install -r requirements.txt

# Run the test
export IDENTECO_API_ENDPOINT=https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev
python apitest.py
```

> **Note** `IDENTECO_API_ENDPOINT` env variable must be set prior running the test.
> The value must be taken from the `serverless deploy` output including stage (i.e. `/dev`)
> but **excluding** the tailing `/` symbol.
> For example: `export IDENTECO_API_ENDPOINT=https://3yhosi5j8l.execute-api.eu-west-1.amazonaws.com/dev`

## Remove

Remove whole stack

> **Note** Manually remove all object from s3 bucket before stack deletion.
> i.e. `aws s3 rm s3://identeco-keys --recursive`

This will remove all underlying resources from the `cloudformation` stack.

```bash
$ serverless remove
```

## Deploy a single lambda function

The following will deploy `register` function

```bash
$ serverless deploy function -f register
```

# Features

-   Registraion of username/password
-   Using assymetric RS256 JWK algorithm
-   Automatic keys rotation

# Known Issues and Limitations

-   Supports only authentication (`username` claim), i.e. identeco confirms that the owner of the claim has `username`
-   No email confirmation
-   No password restrictions
-   No OpenID support
-   Not enough information in logs when error happens

# Roadmap

## v0.1.0-alpha - v0.1.3-dev

-   [x] Port from [identeco](https://github.com/dmsi/identeco)
-   [x] Basic functionality

## v0.1.4-alpha

-   [x] Add slog logger
-   [x] Move serverless to deployment/awslambda directory
-   [x] Wrap errors to provide more context
-   [x] Inject all dependencies in `pkg/runtime`
-   [x] Move AWS Lambda handlers to `pkg/runtime/awslambda`
-   [x] Refactor and separate business logic from AWS Lambda handlers
-   [x] Implement main() for each handler/http server in `cmd`, check [go project layout](https://github.com/golang-standards/project-layout)
-   [x] Implement MongoDb storage backend for users and keys
-   [x] Implement runtime and cmd for standalone HTTP server
-   [x] Use crypto rand for private key generation
-   [ ] Put together the documentation
-   [ ] Revisit `register` it should not return tokens, should return 204
-   [ ] Change module name `github.com/dmsi/identeco` to `github.com/dmsi/identeco-go`
