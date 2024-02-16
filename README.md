PubSub Demo using NATS

Delivery is an at-most-once. For MQTT users, this is referred to as Quality of Service (QoS) 0.
There are two circumstances when a published message won’t be delivered to a subscriber:
The subscriber does not have an active connection to the server (i.e. the client is temporarily offline for some reason)
There is a network interruption where the message is ultimately dropped
Messages are published to subjects which can be one or more concrete tokens, e.g. greet.bob. Subscribers can utilize wildcards to show interest on a set of matching subjects.
