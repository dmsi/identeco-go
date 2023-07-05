# The goal

-   [ ] make this as some sort of library (identeco-core) ? maybe
-   [x] add logger
-   [x] move serverless to deployment/aws-lambda directory or just deployment for now
-   [x] wrap errors? everywhere? wrap-in-wrap?
-   [ ] dynamodb backend for keys
-   [ ] apitest in go
-   [ ] solve the situation with profile and region in serverless.yml
-   [ ] add pkg/runtime where configure runtime. I.e. read env variables and create dependencies.
-   [x] call runtime from `init()` from each handler.
-   [x] move `apitest` directory somewhere, check https://github.com/golang-standards/project-layout
-   [x] in accordance with https://github.com/golang-standards/project-layout move handlers to cmd?
-   [ ] better errors
-   [x] move handlers under runtime/awslambda?
-   [ ] add mongodb storage backend
-   [ ] runtime/httpserver
-   [ ] run in docker + add http server
-   [ ] remove plain rand from key generation
-   [ ] protect private key with password?
-   [ ] configure req/res fields, like "username", "access" or "accessToken" or "access_token" etc

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