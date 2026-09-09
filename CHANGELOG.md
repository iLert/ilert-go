# Changelog

## 09.09.2026, Version 3.25.0

- add `PurchaseSeat` to `CreateUserInput`, mapping to the `purchase-seat` query parameter of `POST /users`. The API then buys a license for the account instead of validating its license quota, which charges the account, prorated for the rest of the billing period. The purchase is unconditional, the API does not check whether a free seat is available first, so the field is only sent when it is set explicitly; the account also needs an active paid subscription and admin seat purchase enabled or the API rejects the request [#80](https://github.com/iLert/ilert-go/pull/80)
- add a `Bool` pointer helper alongside `String`, `Int64` and `Int`, so the `*bool` input fields (`PurchaseSeat`, `SendNoInvitation`, `AbortOnGaps`, `ExcludeOverrides`) can be set without hand-rolling a pointer [#80](https://github.com/iLert/ilert-go/pull/80)
- add `TelemetrySource` with the full set of operations behind `/api/telemetry-sources`: `CreateTelemetrySource`, `GetTelemetrySource`, `GetTelemetrySources`, `SearchTelemetrySource`, `UpdateTelemetrySource`, `DeleteTelemetrySource`, `RotateTelemetrySourceIntegrationKey` and `GetTelemetrySourceStats`, along with `TelemetrySourceType`, `TelemetrySourceStatus` and `TelemetrySourceInclude`. `Teams` is only part of a response when `TelemetrySourceInclude.Teams` is requested, the same way the heartbeat monitor handles `integrationUrl` [#81](https://github.com/iLert/ilert-go/pull/81)
- add `Labels`, `IconUrl`, `Links`, `PublicStatus` and `Dependencies` to `Service`, and the `ServiceLink` and `ServiceDependency` types. `Links`, `PublicStatus` and `Dependencies` are only part of a response when requested through the new `ServiceInclude.Links`, `ServiceInclude.PublicStatus` and `ServiceInclude.Dependencies` [#81](https://github.com/iLert/ilert-go/pull/81)
- add `GetServiceDependencies`, `GetServiceDependency`, `CreateServiceDependency`, `SetServiceDependencies` and `DeleteServiceDependency`. Service dependencies are not writable through the service create and update payloads, they have their own endpoints and their own edge id [#81](https://github.com/iLert/ilert-go/pull/81)
- add `ManagedBy`, returned on services, service dependencies and telemetry sources for entities managed outside of the ilert web app [#81](https://github.com/iLert/ilert-go/pull/81)
- `TelemetrySource.Labels` carries the reserved `ilert.com/discovered-by` key that the API stamps onto every source from its name. The API documents it as read-only derived metadata that clients should omit when sending labels back; the SDK returns it as-is, callers are expected to drop it from their own writes [#81](https://github.com/iLert/ilert-go/pull/81)
- add `ServiceDependency.Type` with `ServiceDependencyType` (`HARD`, `SOFT`). The field is absent from the published OpenAPI spec but the API validates it as an enum, defaults it to `HARD` and returns it on every read, so without it a `SOFT` dependency could neither be created nor read back [#81](https://github.com/iLert/ilert-go/pull/81)
- add `Include` to `SearchServiceInput`, so a service looked up by name can request `links`, `publicStatus` and the other optional properties the endpoint already accepted [#81](https://github.com/iLert/ilert-go/pull/81)
- `ServiceDependency.InvalidAfter` is documented as read-only. The spec marks it writable, but the API ignores it on both POST and PUT in every date format tried and always returns null (verified against the API on 09.09.2026) [#81](https://github.com/iLert/ilert-go/pull/81)
- verified against the API on 09.09.2026: `Service.Labels` and `Service.Teams` are left untouched when absent from an update and cleared by an explicit empty value, while `Service.Links` and `Service.IconUrl` are cleared by absence. The pointer types above are what let a caller tell those two behaviours apart [#81](https://github.com/iLert/ilert-go/pull/81)
- `Service.Labels`, `Service.Links` and `TelemetrySource.Labels`/`TelemetrySource.Teams` are pointer-typed, so a `nil` disappears from the payload and leaves the value untouched while a pointer to an empty map or slice clears it. Without the pointer these fields would either be undroppable, sending `null` on every write of a service whose links or labels are maintained in the web app, or undroppable in the other direction, unable to express "remove them all" — the two halves of the bug fixed for `Teams` in 3.23.1 and 3.24.0 [#81](https://github.com/iLert/ilert-go/pull/81)

## 17.08.2026, Version 3.24.0

- add `LinkTextTemplate` to alert source `LinkTemplate`. The API superseded the plain `text` display name with the templatable `linkTextTemplate` and no longer returns `text` for link templates that use it, so those link templates previously read back with an empty `Text`. `Text` is kept and still accepted by the API as the legacy fallback, but is now deprecated; at least one of the two must be set or the API rejects the request [#77](https://github.com/iLert/ilert-go/pull/77)
- add `NONE` to `SupportStatus` and `SupportStatusAll`. The API accepts and returns it for support hour exceptions that suspend support entirely, but it was missing from the SDK, so callers validating against `SupportStatusAll` rejected exceptions the API itself had created [#78](https://github.com/iLert/ilert-go/pull/78)
- **Source-compatibility note:** `CallFlow.Teams` and `EventFlow.Teams` change type from `[]TeamShort` to `*[]TeamShort`, matching `AlertAction.Teams`. Both endpoints clear the teams on an explicit `null` as well as on an empty array, so the field has to disappear from the payload entirely to mean "leave untouched", which a bare slice cannot express once it is able to marshal an empty array. Callers assigning a slice need to take its address; callers reading teams back through `CallFlowOutput`/`EventFlowOutput` are unaffected. [#79](https://github.com/iLert/ilert-go/pull/79)
- fix `Teams` never being able to clear the teams of a call flow, escalation policy, event flow, heartbeat monitor, incident template, schedule or support hour, the same bug fixed for `AlertSource.Teams` in 3.23.1: the field was tagged `omitempty`, so an empty slice was dropped from the payload and the API left the existing teams in place. The tag is now `json:"teams"` on the five endpoints that ignore a `null`, and a pointer on the two that do not, so an empty slice marshals to `"teams":[]` on all of them. Callers that relied on a `nil` `Teams` meaning "leave untouched" are unaffected. [#79](https://github.com/iLert/ilert-go/pull/79)
- `DeploymentPipeline.Teams`, `Metric.Teams` and `MetricDataSource.Teams` are left tagged `omitempty` on purpose: those endpoints replace the teams with an empty list when the field is absent, so a `nil` slice already clears them [#79](https://github.com/iLert/ilert-go/pull/79)

## 03.08.2026, Version 3.23.1

- fix `AlertSource.Teams` never being able to clear an alert source's teams: the field was tagged `omitempty`, so an empty slice was dropped from the payload and the API left the existing teams in place (it only clears them on an explicit empty array; an omitted or `null` field is a no-op). The tag is now `json:"teams"` so an empty slice marshals to `"teams":[]`. Callers that relied on a `nil` `Teams` meaning "leave untouched" are unaffected, since `nil` still marshals to `null`, which the API ignores. [#76](https://github.com/iLert/ilert-go/pull/76)

## 22.07.2026, Version 3.23.0

- **Source-compatibility note:** the output `CompanyID` field on `AlertActionOutputParams` and `ConnectionOutputParams` changes type from `int64` to `string` (see the Autotask fix below). This is source-breaking for code referencing those fields as `int64`, but it ships in a minor release because the fields never unmarshalled successfully before, so no working caller could depend on the `int64` form. [#73](https://github.com/iLert/ilert-go/pull/73)
- add documented ServiceNow (`closeCode`, `assignmentGroup`, `ownerGroup`, `service`, `serviceOffering`, `contactType`) and Autotask (`noteType`, `notePublish`, `status`) alert action params, service `alias` and call flow VOICEMAIL `disableTranscription` [#73](https://github.com/iLert/ilert-go/pull/73)
- add `VIEWER` to `UserRole` and `TeamMemberRoles` [#73](https://github.com/iLert/ilert-go/pull/73)
- fix reading Autotask alert actions and connections: `companyId`/`queueId` are returned by the API as JSON strings and previously failed to unmarshal; `queueId` now tolerates both encodings, and the output `companyId` is typed `string` to match its API contract and the create/update params (it never unmarshalled successfully before) [#73](https://github.com/iLert/ilert-go/pull/73)

## 08.06.2026, Version 3.22.0

- add `servicesTemplate` (dynamic services mapping), default `services` and `autoCreateServices` to alert source [#70](https://github.com/iLert/ilert-go/pull/70)
- add `acceptAlertOnAnswer` field to call flow CREATE_ALERT node metadata [#71](https://github.com/iLert/ilert-go/pull/71)

## 05.06.2026, Version 3.21.0

- add event flow integration CRUD [#67](https://github.com/iLert/ilert-go/pull/67)
- add setup status to alert source [#68](https://github.com/iLert/ilert-go/pull/68)
- add severityTemplate (dynamic severity mapping) and default severity to alert source [#69](https://github.com/iLert/ilert-go/pull/69)
- distinguish error causes instead of the opaque `An error occurred`: non-JSON error responses now report status, `Server`, request-id and a body snippet, and transport failures are classified (timeout / DNS / TLS / connection)
- detect intermediary blocks (WAF / proxy / load balancer): a non-JSON `403`/`503` is treated as a transient upstream block and retried, while a genuine JSON `4xx` (e.g. invalid credentials) now fails fast instead of being retried until timeout
- classify `404`/`400` responses by status code even when the body is not JSON, so `*NotFoundAPIError` / `*BadRequestAPIError` callers keep working
- keep transient `409`/`423`/`425` responses retryable (concurrent edit / eventual consistency)
- add `ILERT_DEBUG` (parsed with `strconv.ParseBool`) and a `WithDebug` option to enable full request/response tracing; the `Authorization` header is masked in the debug output

## 02.06.2026, Version 3.20.0

- add `AddAlertSourceToAlertAction` and `RemoveAlertSourceFromAlertAction` methods for non-destructive alert-source attachment [#66](https://github.com/iLert/ilert-go/pull/66)

## 22.04.2026, Version 3.19.0

- add reroute param for alert action [#65](https://github.com/iLert/ilert-go/pull/65)

## 18.03.2026, Version 3.18.1

- add groups to status page [#64](https://github.com/iLert/ilert-go/pull/64)

## 13.03.2026, Version 3.18.0

- add teams to escalation rule in escalation policy [#63](https://github.com/iLert/ilert-go/pull/63)
- add custom http headers for alert action [#62](https://github.com/iLert/ilert-go/pull/62)

## 11.03.2026, Version 3.17.0

- add event flows [#61](https://github.com/iLert/ilert-go/pull/61)

## 06.03.2026, Version 3.16.0

- add holiday exceptions for support hours [#60](https://github.com/iLert/ilert-go/pull/60)
- add new call flow node type [#59](https://github.com/iLert/ilert-go/pull/59)

## 24.02.2026, Version 3.15.2

- fix alert source active field not propagating to API when set to false [#58](https://github.com/iLert/ilert-go/pull/58)

## 06.01.2026, Version 3.15.1

- add missing integration type for new heartbeat alert sources [#57](https://github.com/iLert/ilert-go/pull/57)
  - deprecate legacy heartbeat integration type

## 01.10.2025, Version 3.15.0

- add logs and comments to events [#56](https://github.com/iLert/ilert-go/pull/56)

## 03.09.2025, Version 3.14.0

- add new integration types [#55](https://github.com/iLert/ilert-go/pull/55)

## 17.06.2025, Version 3.13.0

- add call flows [#52](https://github.com/iLert/ilert-go/pull/52)

## 01.07.2025, Version 3.12.3

- fix type of `labels` in Event [#54](https://github.com/iLert/ilert-go/pull/54)

## 01.07.2025, Version 3.12.2

- add new `labels` fields [#53](https://github.com/iLert/ilert-go/pull/53)

## 06.06.2025, Version 3.12.1

- change heartbeat monitor search to always return integration url [#51](https://github.com/iLert/ilert-go/pull/51)

## 24.04.2025, Version 3.12.0

- new alert source fields [#50](https://github.com/iLert/ilert-go/pull/50)
  - add new heartbeat monitor resource
  - deprecate legacy email and heartbeat fields
  - add new email integration type

## 08.04.2025, Version 3.11.2

- fix typo in json mapping for send notification field in incident template resource [#49](https://github.com/iLert/ilert-go/pull/49)

## 29.11.2024, Version 3.11.1

- add type GitLab for deployment pipeline [#48](https://github.com/iLert/ilert-go/pull/48)

## 07.11.2024, Version 3.11.0

- add new resource deployment pipeline [#47](https://github.com/iLert/ilert-go/pull/47)

## 01.10.2024, Version 3.10.1

- fix to ensure backwards compatibility when using conditions field [#46](https://github.com/iLert/ilert-go/pull/46)

## 01.10.2024, Version 3.10.0

- add conditions field to alert action [#45](https://github.com/iLert/ilert-go/pull/45)

## 16.09.2024, Version 3.9.1

- add additional body template fields [#44](https://github.com/iLert/ilert-go/pull/44)

## 22.08.2024, Version 3.9.0

- add new api resources/fields pt.1 [#43](https://github.com/iLert/ilert-go/pull/43)
  - alert action
    - deprecate `delaySec` in favor of more specific `escalationEndedDelaySec` and `notResolvedDelaySec`
    - new trigger type `AlertNotResolved`
    - new alert action type `SlackWebhook`
  - alert source
    - new alert grouping type `intelligentGrouping`
      - add field `scoreThreshold`
    - add event filter
    - add includes for POST and PUT API calls
  - status page
    - add email login via `emailWhitelist`
    - add `announcement` fields
    - add `metrics`

## 09.05.2024, Version 3.8.1

- add region to user [#41](https://github.com/iLert/ilert-go/pull/41)

## 06.05.2024, Version 3.8.0

- add send-no-invitation option to user create api [#40](https://github.com/iLert/ilert-go/pull/40)

## 06.05.2024, Version 3.7.1

- readd removed connector and alert action for microsoft teams simple webhook [#39](https://github.com/iLert/ilert-go/pull/39)

## 25.04.2024, Version 3.7.0

- remove connectors and alert actions deprecated via api in [#38](https://github.com/iLert/ilert-go/pull/38)
  - adds support for alert actions and connectors previously missing

## 20.02.2024, Version 3.6.1

- fix backwards compatibility with alert actions v2 changes [#37](https://github.com/iLert/ilert-go/pull/37)
  - ensures existing scripts using one alert source with no team explicitly set to use legacy api without breaking

## 01.02.2024, Version 3.6.0

- apply alert actions v2 changes in [#29](https://github.com/iLert/ilert-go/pull/29)
  - add alertSources and teams fields, deprecate alertSourceIds

## 12.01.2024, Version 3.5.0

- add Telegram as alert action type in [#36](https://github.com/iLert/ilert-go/pull/36)

## 05.01.2024, Version 3.4.1

- replace `ThemeMode` field with `Appearance` for status page resource in [#35](https://github.com/iLert/ilert-go/pull/35)

## 03.01.2024, Version 3.4.0

- deprecate uptime monitors in [#32](https://github.com/iLert/ilert-go/pull/32)
- add new resource support hours in [#33](https://github.com/iLert/ilert-go/pull/33)
- add new status page fields in [#34](https://github.com/iLert/ilert-go/pull/34)

## 12.12.2023, Version 3.3.0

- add link templates and priority template to alert source resource in [#31](https://github.com/iLert/ilert-go/pull/31)

## 13.11.2023, Version 3.2.0

- add new trigger type `alert-escalation-ended` and new field `delaySec` to alert action resource in [#30](https://github.com/iLert/ilert-go/pull/30)

## 09.10.2023, Version 3.1.0

- add new fields `delayMin` and `routingKey` to escalation policy resource in [#26](https://github.com/iLert/ilert-go/pull/26)
- add templates and alert grouping to alert source resource in [#27](https://github.com/iLert/ilert-go/pull/27)
- add new trigger types to alert action resource in [#28](https://github.com/iLert/ilert-go/pull/28)

## 02.05.2023, Version 3.0.2

- add missing field `accountWideView` to status page resource in [#25](https://github.com/iLert/ilert-go/pull/25)

## 13.03.2023, Version 3.0.1

- add missing field `integrationKey` to metric resource in [#24](https://github.com/iLert/ilert-go/pull/24)

# 08.03.2023, Version 3.0.0 - API user preference migration: see [migration changes](https://docs.ilert.com/rest-api/api-version-history/api-user-preference-migration-2023#migrating-ilert-go-and-or-terraform) for a detailed migration guide

- removed notification settings fields from user resource
- add user contacts
  - email
  - phone number
- add user notification preferences
  - alert (alert creation)
  - duty (on-call)
  - subscription (subscriber to incident, service, status page)
  - update (alert update changes)

## 17.01.2023, Version 2.6.0

- add metrics
- add metric data sources
- add series

## 13.01.2023, Version 2.5.1

- fix missing json params in status page group

## 12.01.2023, Version 2.5.0

- add status page groups
- add status page structure for groups

## 11.01.2023, Version 2.4.0

- add pagination for alert actions
- add pagination for alert sources
- add pagination for connectors
- add pagination for escalation policies
- add pagination for schedules
- add pagination for teams
- add pagination for uptime monitors
- add pagination for users
- enforce pagination for alerts
- enforce pagination for automation rules
- enforce pagination for incident templates
- enforce pagination for incidents
- enforce pagination for services
- enforce pagination for status pages

## 28.11.2022, Version 2.3.3

- fix typo in dingtalk alert action type

## 25.11.2022, Version 2.3.2

- add alert action alert filter

## 24.11.2022, Version 2.3.1

- add missing connector types for Terraform provider

## 21.11.2022, Version 2.3.0

- deprecate automation rule
- add automation rule as connector type under alert action
- add status page ip filter
- add dingtalk as an alert action / connector

## 17.11.2022, Version 2.2.3

- add dingtalk alert action and connector support

## 08.09.2022, Version 2.2.2

- add fields for multiple responders in escalation rule (escalation policy)

## 06.09.2022, Version 2.2.1

- fix alert action search method
- fix connector search method
- change user search

## 04.09.2022, Version 2.2.0

- add search by name for all entities

## 30.08.2022, Version 2.1.1

- add various validation lists to schedule

## 26.08.2022, Version 2.1.0

- add schedules

## 12.07.2022, Version 2.0.5

- another fix for alert source support hours check

## 12.07.2022, Version 2.0.4

- fix alert source support hours check

## 12.07.2022, Version 2.0.3

- fix error codes on create resources

## 08.07.2022, Version 2.0.2

- add various lists
- fix alert source creation
- add deprecated hints
- fix user alert notification type list

## 04.07.2022, Version 2.0.1

- fix internal version

## 30.06.2022, Version 2.0.0 - API Version upgrade

- add automation rules
- add incident templates
- add services
- add statuspages
- migrating API v1 to versionless
  - rename incident v1 to alert and update fields
  - add incident
  - renaming connection to alert action, therefore deprecating connection
  - update alert source fields
  - update event fields
  - update uptime monitor fields
  - update user fields
- update examples

## 16.04.2022, Version 1.6.5

- add not_found error types

## 16.04.2022, Version 1.6.4

- fix connector types

## 13.04.2022, Version 1.6.3

- add retryable errors

## 18.01.2022, Version 1.6.2

- add new alert source types

## 18.01.2022, Version 1.6.1

- fix ssl uptime monitor check params

## 18.01.2022, Version 1.6.0

- add more uptime monitor check params

## 14.04.2021, Version 1.5.1

- add auto raise incidents prop to support hours

## 13.04.2021, Version 1.5.0

- add proxy option to the client

## 09.04.2021, Version 1.4.1

- fix connection output types

## 08.04.2021, Version 1.4.0

- return generic api error for each resource output if an api error occurred
- add new alert source types
- add new connector types
- add new connection types

## 07.04.2021, Version 1.3.1

- add default retry option to the client

## 07.04.2021, Version 1.3.0

- add retry exponential backoff option
- upgrade to golang 1.6

## 12.03.2021, Version 1.2.2

- fix team visibility

## 12.03.2021, Version 1.2.1

- add team member role list

## 23.02.2021, Version 1.2.0

- add teams
- add url option for create event
- add autotask types

## 09.11.2020, Version 1.1.3

- add integration url to alert source definition

## 08.11.2020, Version 1.1.2

- add jira alert source type

## 03.11.2020, Version 1.1.1

- add all connector types list

## 02.11.2020, Version 1.1.0

- add connections
- add connectors
- add application default environment variables
- remove type property from escalation rule

## 20.10.2020, Version 1.0.2, 1.0.3, 1.0.4, 1.0.5, 1.0.6, 1.0.7, 1.0.8, 1.0.9

- fix alert source type definition
- fix incident type definition
- fix uptime monitor type definition
- fix user type definition
- fix escalation policy type definition
- add auto resolution timeout option to alert source

## 19.10.2020, Version 1.0.0, 1.0.1

- add user agent option for tools like terraform
- add events
- add users
- add alert sources
- add escalation policies
- add schedules
- add numbers
- add uptime monitors
- add incidents
- add heartbeats
