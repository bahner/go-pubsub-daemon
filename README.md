# myspace-pubsub-daemon

This is a simple libp2p pubsub daemon. It is written because
pubsub is being removed from kubo.

You can subscribe to a libp2p pubsub topic, eg. `MYTOPIC`,
by connecting to `http://localhost:5002/api/v0/topics/MYTOPIC`
where you will receive a websocket that reads and writes
to this libp2p pubsub topic.

Easy peasy. Lemon squeezy.

If you run a server on most any computer it will connect to
peers with the same rendezvous string: `myspace`.

## TL;DR

Docker: `make all && docker-compose up`
Development: `make serve`
Install: `make install # Installs to /usr/local/myspace-pubsub-daemon by default.`

## API

A simple RESTful API at localhost:5002/api/v0. The following endpoints
exists.

- GET /topics - List topics
- GET /topics/topicName - Establishes a websocket that subscribes to the topic
- GET /topics/topicName/peers - List peers for the topics. NB! Local peers does not show up.

## Daemon

The daemon starts a libp2p node and connects via rendezvous
discovery. It listen on all interfaces on a random tcp port.

For options run `./myspace-pubsub-daemon -help`. Please note
that the default reported reflects your current environment
settings.

To build the binary, simply run `make build`

## Configuration

Some configuration can be edited in `.env`. `make` picks up
on changes here.

There is not much in way of configuration. You can set the
listen address and port for the web sockets. The default is
to listen on port 127.0.0.1:5002.

The evaluation order for settings is:

- defaults
- environment variables
- command line parameters

## Environment variables

NB! If these variables are set in your environment they
will be defaults for the application. The defaults are set
in `.env` and are as follows.

```bash
MYSPACE_PUBSUB_DAEMON_PORT="5002" # The port the daemon is listens on
MYSPACE_PUBSUB_DAEMON_ADDR="127.0.0.1" # The interface the daemon binds to
MYSPACE_PUBSUB_DAEMON_LOG_LEVEL="info" # Log level for the libp2p part of the daemon

IMAGE=docker.io/bahner/myspace-pubsub-daemon:latest # The name of the docker image to be used or built

GO=go # The name of go binary to use
BUILD_IMAGE=golang:alpine # Image to used for building the daemon in docker.
```

## Docker

Inside the docker container the service listens on 0.0.0.0.
This is hardcoded in the Dockerfile.

The `MYSPACE_PUBSUB_DAEMON_ADDR` variable is used by
docker-compose for binding to a port on the host system.

To start the service in docker simply run: `docker-compose up`.

If you want to create your own version or changes, edit `.env`
and run:

```bash
make image
docker-compose up -d
```

## systctl.sh

This a helper script, to set required limits in sysctl.
You probably wanna do this somewhere else. Use if if you
receive fatal resource errors from the daemon.

## Client

You can connect and chat with the client. But this is *not*
intended to be used by hoomans as a chat tool.
The client is just a debug/inspection tool.

For options please type `./client/client -help`

2023-07-16: bahner
