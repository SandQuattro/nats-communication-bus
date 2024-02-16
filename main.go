package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

// Delivery is an at-most-once. For MQTT users, this is referred to as Quality of Service (QoS) 0.
// There are two circumstances when a published message won’t be delivered to a subscriber:
// The subscriber does not have an active connection to the server (i.e. the client is temporarily offline for some reason)
// There is a network interruption where the message is ultimately dropped
// Messages are published to subjects which can be one or more concrete tokens, e.g. greet.bob.
// Subscribers can utilize wildcards to show interest on a set of matching subjects.
func main() {

	// Use the env variable if running in the container,
	// otherwise use the default - nats://127.0.0.1:4222
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	// Create an unauthenticated connection to NATS.
	nc, _ := nats.Connect(url)

	// Drain is a safe way to ensure all buffered messages that were published are sent and all buffered messages
	// received on a subscription are processed before closing the connection.
	defer nc.Drain()

	// Messages are published to subjects.
	// Although there are no subscribers, this will be published successfully.
	nc.Publish("greet.joe", []byte("hello"))

	// Let’s create a subscription on the greet.* wildcard.
	sub, _ := nc.SubscribeSync("greet.*")

	// For a synchronous subscription, we need to fetch the next message.
	// However.. since the publish occured before the subscription was established, this is going to timeout.
	msg, _ := sub.NextMsg(10 * time.Millisecond)
	fmt.Println("subscribed after a publish...")
	fmt.Printf("msg is nil? %v\n", msg == nil)

	// Publish a couple messages.
	nc.Publish("greet.joe", []byte("hello"))
	nc.Publish("greet.pam", []byte("hello"))

	// Since the subscription is established, the published messages will immediately be broadcasted to all subscriptions.
	// They will land in their buffer for subsequent NextMsg calls.
	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	nc.Publish("greet.bob", []byte("hello"))

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)
}
