# Introducing Event Sourcing Kit or ESK

ESK is a simple library for building a project following the [event sourcing pattern](https://martinfowler.com/eaaDev/EventSourcing.html). The core tool is a broker which enables publishing and subscribing to topics via web sockets and storing events on those topics. Think Kafka but lightweight, written in go, with simplicity in mind. It's currently in pre-pre-alpha and we're looking for help!

## Getting Started

### Installation

Homebrew (macOS)

```
brew install wcatron/esk/esk
```

From Source: 

*See below contributing guide*

### Start Instance

Navigate to a directory where your data files can be added and run the following.

```
esk
```

It's that simple! Open a browser to view the playground and start publishing and subscribing to events. Checkout a working demo using our [Pointer Demo](https://github.com/wcatron/esk-pointer-demo).

## Current Features

- A basic broker which can be connected to via WebSockets
- Subscribe and publish to topics
- Write events to binary files for each topic
- Every event is a line in the file
- On subscribe receive past event sent to the topic
- Simple cli tooling (`./esk --port 8080`)
- Clients can subscribe to a topic with a cursor and receive messages since that point
- Client library: [typescript](https://github.com/wcatron/esk-client-typescript)

### Planned Features

- Subscribe only to the last event.
- [Cluster support](./ClusterSupport.md) for scalability and high availability
- And much much more...

## Contributing

*This project is in very early stages. This section only contains currently implemented functionality*

Clone the repository then build the application and the demo playground:

`go build`
`cd demo && yarn`

Start the server by running:

`./esk`

Open the browser to see a simple page to publish and subscribe. Open multiple windows and subscribe to a topic `SUB /topic`. Then publish to the topic from one window `PUB /topic "Hello"` and see the message appear in the console of the other window. Observe every time you subscribe to a topic the broker will send all previous messages sent on the client from oldest to newest.

## Testing

Run all tests in all packages:

`go test ./...`
