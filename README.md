# micro-backends

[![Netlify Status](https://api.netlify.com/api/v1/badges/f2393a44-27a4-4fbb-9c10-63d1f8659b38/deploy-status)](https://app.netlify.com/sites/micro-backends/deploys)
[![Continuous Integration](https://github.com/kucera-lukas/micro-backends/actions/workflows/ci.yml/badge.svg)](https://github.com/kucera-lukas/micro-backends/actions/workflows/ci.yml)
[![Continuous Integration](https://github.com/kucera-lukas/micro-backends/actions/workflows/cd.yml/badge.svg)](https://github.com/kucera-lukas/micro-backends/actions/workflows/cd.yml)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/stegoer/server/main.svg)](https://results.pre-commit.ci/latest/github/kucera-lukas/micro-backends/main)

micro-backends is using Go, RabbitMQ, PostgreSQL, MongoDB, gRPC and GraphQL.

## Installation

```sh
git clone git@github.com:kucera-lukas/micro-backends.git
```

### Environment variables

Each service has a separate `.env` file.
Copy the contents of `.env.example` file and fill in the required values.

```sh
cp .env.example .env
```

### Docker

Installation using Docker is recommended because application requires
many services to run at once and `docker-compose` simplifies the process significantly.

```sh
make build
```

## Development

```sh
make dev
```

```sh
# build images before starting containers
make dev-build
```

### Make

```sh
make help
```

### Tools

Tools is a small Go module which contains CLI tooling which is used mainly via `Makefile` commands.

## Contributing

```sh
pre-commit install
```

## License

Developed under the [MIT](https://github.com/kucera-lukas/micro-backends/blob/master/LICENSE) license.
