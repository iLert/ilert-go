# ilert-go

**The official ilert Go api bindings.**

## Create an alert (manually)

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v3"
)

func main() {

	// We strongly recommend to enable a retry logic if an error occurs
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithRetry(10, 5*time.Second, 20*time.Second), ilert.WithAPIToken(apiToken))

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
	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiKey = "heartbeat API Key"
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

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
	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithProxy("http://proxyserver:8888"), ilert.WithAPIToken(apiToken))
	...
}
```

## Versions overview

If you want to use older legacy versions of ilert-go, you can access previous major versions using one of the commands below.

| Version       | Description                                                                                                                       | Command                               |
| ------------- | --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| > 3.0.0       | API user preference migration - [changes](https://docs.ilert.com/rest-api/api-version-history/api-user-preference-migration-2023) | `go get github.com/iLert/ilert-go/v3` |
| 2.0.0 - 2.6.0 | API versionless - [changes](https://docs.ilert.com/rest-api/api-version-history#renaming-incidents-to-alerts)                     | `go get github.com/iLert/ilert-go/v2` |
| 1.0.0 - 1.6.5 | API v1 - basic legacy resources                                                                                                   | `go get github.com/iLert/ilert-go`    |

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
