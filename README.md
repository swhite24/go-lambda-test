# go-lambda-test

Example REST API built on lambda, using golang with the [gin](https://github.com/gin-gonic/gin) framework.

Part of the service declaration is provisioning a DynamoDB table and updating the lambda execution role with permissions for access.

## Development

The entry-point handler is configured to run locally based on the value of `GIN_MODE` in the environment. If `GIN_MODE` is set to `release`, it is assumed that the application is being executed in lambda. Otherwise, a local http server will be started on port `3000`.

A [realize](https://github.com/oxequa/realize) configuration is also provided.

To run locally:

```bash
realize start
```

## Deployment

Full deployment:

```bash
make deploy
```

Quick deploy of handler function:

```bash
make deploy-quick
```

## License

MIT. See [LICENSE](LICENSE).
