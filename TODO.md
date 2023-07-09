# The goal

-   [ ] make this as some sort of library (identeco-core) ? maybe
-   [x] add logger
-   [x] move serverless to deployment/aws-lambda directory or just deployment for now
-   [x] wrap errors? everywhere? wrap-in-wrap?
-   [ ] dynamodb backend for keys
-   [ ] apitest in go
-   [ ] solve the situation with profile and region in serverless.yml
-   [x] add pkg/runtime where configure runtime. I.e. read env variables and create dependencies.
-   [x] call runtime from `init()` from each handler.
-   [x] move `apitest` directory somewhere, check https://github.com/golang-standards/project-layout
-   [x] in accordance with https://github.com/golang-standards/project-layout move handlers to cmd?
-   [ ] better errors
-   [x] move handlers under runtime/awslambda?
-   [x] add mongodb storage backend
-   [x] runtime/httpserver
-   [ ] run in docker + add http server
-   [x] remove plain rand from key generation
-   [ ] protect private key with password?
-   [ ] configure req/res fields, like "username", "access" or "accessToken" or "access_token" etc
-   [x] rename storage modules "keydatas3", "keydatamdb" or "keydatamongo", "userdataddb" or "userdatadynamo", "keydata_mdb" ???
-   [ ] put together documentation
-   [ ] rename "github.com/dmsi/identeco-go"?
-   [ ] add `model` which consumes `storage` and provides easier-to-use serialization/deserialization
-   [ ] delete user
-   [ ] delete user testcases
-   [ ] store refresh in DB
-   [ ] TODO document env variables for both awslambda and httpserver
-   [ ] when it comes to error it's not clear how to pin-point the issue. wrap all the errors unconditionally?
-   [x] move AWS Lambda handlers to runtime
-   [x] handlers.Handler naming using for both httpserver and awslambda is confusing. One way to fix is to name packages differently i.e. httphandlers vs awslambdahandlers. Another way to fix is flattern the structure of the runtime and move handlers, routers, etc to the awslambda and httpserver packages

# Env

```
IDO_DDB_TABLE_NAME
IDO_S3_BUCKET_NAME
IDO_PRIVATE_KEY_NAME
IDO_JWKS_JSON_NAME
IDO_PRIVATE_KEY_BITS
IDO_PRIVATE_KEY_LIFETIME
IDO_ACCESS_TOKEN_LIFETIME
IDO_REFRESH_TOKEN_LIFETIME
IDO_CLAIM_ISS
```
