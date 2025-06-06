# Changelog

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

# 20.02.2024, Version 3.6.1

- fix backwards compatibility with alert actions v2 changes [#37](https://github.com/iLert/ilert-go/pull/37)
    - ensures existing scripts using one alert source with no team explicitly set to use legacy api without breaking

# 01.02.2024, Version 3.6.0

- apply alert actions v2 changes in [#29](https://github.com/iLert/ilert-go/pull/29)
    - add alertSources and teams fields, deprecate alertSourceIds

# 12.01.2024, Version 3.5.0

- add Telegram as alert action type in [#36](https://github.com/iLert/ilert-go/pull/36)

# 05.01.2024, Version 3.4.1

- replace `ThemeMode` field with `Appearance` for status page resource in [#35](https://github.com/iLert/ilert-go/pull/35)

# 03.01.2024, Version 3.4.0

- deprecate uptime monitors in [#32](https://github.com/iLert/ilert-go/pull/32)
- add new resource support hours in [#33](https://github.com/iLert/ilert-go/pull/33)
- add new status page fields in [#34](https://github.com/iLert/ilert-go/pull/34)

# 12.12.2023, Version 3.3.0

- add link templates and priority template to alert source resource in [#31](https://github.com/iLert/ilert-go/pull/31)

# 13.11.2023, Version 3.2.0

- add new trigger type `alert-escalation-ended` and new field `delaySec` to alert action resource in [#30](https://github.com/iLert/ilert-go/pull/30)

# 09.10.2023, Version 3.1.0

- add new fields `delayMin` and `routingKey` to escalation policy resource in [#26](https://github.com/iLert/ilert-go/pull/26)
- add templates and alert grouping to alert source resource in [#27](https://github.com/iLert/ilert-go/pull/27)
- add new trigger types to alert action resource in [#28](https://github.com/iLert/ilert-go/pull/28)

## 02.05.2023, Version 3.0.2

- add missing field `accountWideView` to status page resource in [#25](https://github.com/iLert/ilert-go/pull/25)

## 13.03.2023, Version 3.0.1

- add missing field `integrationKey` to metric resource in [#24](https://github.com/iLert/ilert-go/pull/24)

## 08.03.2023, Version 3.0.0 - API user preference migration: see [migration changes](https://docs.ilert.com/rest-api/api-version-history/api-user-preference-migration-2023#migrating-ilert-go-and-or-terraform) for a detailed migration guide

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
