# Changelog

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

## 30.06.2022, Version 2.0.0

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

## 9.11.2020, Version 1.1.3

- add integration url to alert source definition

## 8.11.2020, Version 1.1.2

- add jira alert source type

## 3.11.2020, Version 1.1.1

- add all connector types list

## 2.11.2020, Version 1.1.0

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
