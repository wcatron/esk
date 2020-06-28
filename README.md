# Introducing Event Sourcing Kit or ESK

ESK is a simple tool for building a project following the event sourcing pattern. The core tool is a broker which enables publishing and subscribing to topics via web sockets and storing events on those topics. Think Kafka but lightweight, written in go, with simplicity in mind. It's currently in pre-pre-alpha and we're looking for help!

## Current Features

- A basic broker which can be connected to via WebSockets
- Subscribe and publish to topics
- Write events to a file for each topic
- On subscribe receive past events sent to the topic
- Simple cli tooling (`./esk --port 8080`)
- 
- Typescript client library

### Planned Features

- Support for clients to subscribe to a topic with a cursor to reduce redundant INFORM messages
- [Cluster support](./ClusterSupport.md) for scalability and high availability
- And much much more...

## Getting Started

*This section only contains currently implemented functionality*

Clone the repository then build the application and the demo playground:

`go build`
`cd demo && yarn`

Start the server by running:

`./esk`

Open the browser to see a simple page to publish and subscribe. Open multiple windows and subscribe to a topic `SUB /topic`. Then publish to the topic from one window `PUB /topic "Hello"` and see the message appear in the console of the other window. Observe every time you subscribe to a topic the broker will send all previous messages sent on the client from oldest to newest.

## Building

The goal is to setup a something like [pkger](https://github.com/markbates/pkger) to bundle the demo code however currently it wont work in the CI process.

## Testing

Run all tests in all packages:

`go test ./...`
