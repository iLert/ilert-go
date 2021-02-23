# ilert-go

**The official iLert Go api bindings.**

## Create an incident (manually)

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go"
)

func main() {

	client := ilert.NewClient()

	var apiKey = "alert source API Key"
	event := &ilert.Event{
		APIKey:      apiKey,
		EventType:   ilert.EventTypes.Alert,
		Summary:     "My test incident summary",
		IncidentKey: "123456",
	}

	input := &ilert.CreateEventInput{Event: event}
	result, err := client.CreateEvent(input)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}

	log.Println("Incident key:", result.EventResponse.IncidentKey)
}
```

## Ping heartbeat

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go"
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

## Getting help

We are happy to respond to [GitHub issues][issues] as well.

[issues]: https://github.com/iLert/ilert-go/issues/new

<br>

#### License

<sup>
Licensed under <a href="LICENSE-APACHE">Apache License, Version
2.0</a>
</sup>

<br>

<sub>
Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in ilert-rust by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.
</sub>
