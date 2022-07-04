# ilert-go

**The official iLert Go api bindings.**

## Legacy API

> If you want to use the old client with v1 resources, please install it with the following command: `go get github.com/iLert/ilert-go`


## Create an incident (manually)

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v2"
)

func main() {

	// We strongly recommend to enable a retry logic if an error occurs
	client := ilert.NewClient(ilert.WithRetry(10, 5*time.Second, 20*time.Second))

	var apiKey = "alert source API Key"
	event := &ilert.Event{
		APIKey:      apiKey,
		EventType:   ilert.EventTypes.Alert,
		Summary:     "My test alert summary",
		AlertKey: "123456",
	}

	input := &ilert.CreateEventInput{Event: event}
	result, err := client.CreateEvent(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}

	log.Println("Event processed!")
}
```

## Ping heartbeat

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v2"
)

func main() {

	client := ilert.NewClient()

	var apiKey = "heartbeat API Key"
	result, err := client.PingHeartbeat(&ilert.PingHeartbeatInput{
		APIKey: ilert.String(apiKey),
		Method: ilert.String(ilert.HeartbeatMethods.HEAD),
	})

	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}

	log.Println("Heartbeat is ok!")
}
```

## Using proxy

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v2"
)

func main() {
	client := ilert.NewClient(ilert.WithProxy("http://proxyserver:8888"))
	...
}
```

## Getting help

We are happy to respond to [GitHub issues][issues] as well.

[issues]: https://github.com/iLert/ilert-go/issues/new

<br>

#### License

<sup>
Licensed under <a href="LICENSE">Apache License, Version
2.0</a>
</sup>

<br>

<sub>
Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in ilert-go by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.
</sub>
