# How to build and host the Chatterino 2 API

## Prerequisites

1. Resolved links are stored in PostgreSQL, so you must have PostgreSQL installed and accessible for the user running the API. For Ubuntu, you would install it with `sudo apt install postgresql`, create a DB user for your system user (`sudo -upostgres createuser pajlada`), then create a db for the api (`sudo -upostgres createdb chatterino-api --owner pajlada`).

## Build

1. Clone the repo: `git pull https://github.com/Chatterino/api.git`
1. Move into the directory, fetch the dependencies and build the binary: `cd api && make build`
1. Edit API's configuration, see [configuration](./config.md)
1. Start the API: `./cmd/api/api`

## Install systemd unit

1. Install the pre-packaged systemd unit file: `sudo cp ./docs/chatterino-api.service /etc/systemd/system/`
1. Use your preferred editor to edit the service file with all the necessary details (incl. [API keys](./apikeys.md)): `sudo editor /etc/systemd/system/chatterino-api.service`
1. Tell systemd to reload the changes: `sudo systemctl daemon-reload`
1. Start the service and enable to start on boot: `sudo systemctl enable --now chatterino-api`
