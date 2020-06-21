# Cluster Support

When operating in a cluster the following must happen.

- When one node receives an event it must be published to all other nodes and their subscribers.
- When that event is written to the data store, it is written in the same order it is published to subscribers.
- If order is considered important then clusters will need to communicate locks for topics.