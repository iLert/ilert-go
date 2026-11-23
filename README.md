# ilert-go

**The official ilert Go api bindings.**

## Create an alert (manually)

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v4"
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
	"github.com/iLert/ilert-go/v4"
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
	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithProxy("http://proxyserver:8888"), ilert.WithAPIToken(apiToken))
	...
}
```

## Early access

Some APIs of ilert are in [early access](https://docs.ilert.com/developer-docs/rest-api/api-reference/introduction/early-access): actions, agents and credentials. They may change or be removed, even in a minor version of ilert-go, so a client has to opt in to each one it uses, by its path under `/api/beta`: `"actions"` for `/api/beta/actions`, `"agents"` for `/api/beta/agents` and `"credentials"` for `/api/beta/credentials`. Without the opt-in, their operations fail with an `EarlyAccessDisabledError` before anything is sent.

```go
package main

import (
	"log"
	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(
		ilert.WithAPIToken(apiToken),
		ilert.WithEarlyAccess("actions", "agents", "credentials"),
	)
	...
}
```

The name of an API that leaves early access stays accepted by `WithEarlyAccess`, without effect, so code that opts in to it keeps working.

## Versions overview

If you want to use older legacy versions of ilert-go, you can access previous major versions using one of the commands below.

| Version       | Description                                                                                                                       | Command                               |
| ------------- | --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| >= 4.0.0      | Incidents split into operational incidents and status updates - [changes](CHANGELOG.md)                                          | `go get github.com/iLert/ilert-go/v4` |
| 3.0.0 - 3.26.0 | API user preference migration - [changes](https://docs.ilert.com/developer-docs/api-version-history/api-user-preference-migration-2023) | `go get github.com/iLert/ilert-go/v3` |
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
