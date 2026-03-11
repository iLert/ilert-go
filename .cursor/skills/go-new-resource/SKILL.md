---
name: go-new-resource
description: Creates new ilert-go API resource files by mirroring established patterns in call_flow.go and support_hour.go, including declaration order, var-based enums, CRUD method structure, apiRoutes wiring, and Go coding style. Use when adding a new resource file such as event_flow.go or extending the client with a new entity.
---

# Create ilert-go Resource

## When To Use

Use this skill when:
- adding a new `ilert-go` resource file like `<entity>.go`
- cloning an existing resource pattern (for example `call_flow.go`)
- wiring a new endpoint into `apiRoutes` in `client.go`
- keeping style identical to existing `ilert-go` files

## Source Of Truth Files

- `call_flow.go` (most complete node-based resource pattern)
- `support_hour.go` (clean CRUD baseline)
- `team.go` (enum var + `All` slice style)
- `client.go` (`apiRoutes` definition + initialization style)

## Non-Negotiable Conventions

1. Keep import block style consistent:
   - `encoding/json`
   - `errors`
   - `fmt`
   - `net/url`
   - `strconv`
   - include only what is used
2. Use GoDoc comments for exported declarations.
3. Use explicit JSON tags on all API structs.
4. Use `var` struct blocks for enum-like values, not `const` blocks.
5. Add corresponding `<EnumName>All` slices after each enum var block.
6. Use operation wrappers with `_ struct{}` sentinel field.
7. Use operation order exactly:
   - `Create<Entity>`
   - `Get<Entity>`
   - `Get<Entities>`
   - `Search<Entity>`
   - `Update<Entity>`
   - `Delete<Entity>`
8. Use early-return input validation with `errors.New(...)`.
9. Use request pattern:
   - `c.httpClient.R()`
   - `SetBody(...)` for create/update
   - `fmt.Sprintf(...)` endpoint composition
10. Use response pattern:
    - `getGenericAPIError(resp, <expectedStatus>)`
    - `json.Unmarshal(resp.Body(), ...)`
    - return typed output wrapper

## File Construction Order

Build the new `<entity>.go` in this exact order:

1. Domain model structs (`<Entity>`, `<Entity>Output`, nested models).
2. Enum-like var groups and `<Name>All` slices.
3. Input/output structs for create/get/list/search/update/delete.
4. Client methods in matching CRUD/search order.

## Standard Operation Shapes

- `Create`:
  - validate `input != nil`
  - validate payload pointer is not nil
  - `POST` base route
  - expect `201`
- `Get` by id:
  - validate `input != nil`
  - validate id pointer
  - `GET %s/%d`
  - expect `200`
- `Get` list:
  - use `url.Values{}`
  - include `start-index` and `max-results` when provided
  - `GET %s?%s`
  - expect `200`
- `Search` by name:
  - validate name pointer
  - `GET %s/name/%s`
  - expect `200`
- `Update`:
  - validate entity pointer and id pointer
  - `PUT %s/%d`
  - expect `200`
- `Delete`:
  - validate id pointer
  - `DELETE %s/%d`
  - expect `204`

## Route Wiring Rule

When adding a new resource:
- add a new field in `apiRoutes` struct in `client.go`
- add matching initialization entry in the same block
- keep naming pattern aligned with existing route keys

## Style Guardrails

- Keep names entity-prefixed (`EventFlowNode`, `CreateEventFlowInput`).
- Keep optional request values as pointers where existing style does that.
- Keep list outputs as slices of pointers when existing style does that.
- Preserve short, direct error strings (for example `"input is required"`).
- Do not refactor unrelated files while adding a new resource.

## Minimal Verification Checklist

- `gofmt -w <new_file>.go client.go`
- `go test ./...`
- confirm:
  - operation order matches existing resources
  - enum var + `All` slices exist
  - route key exists in both struct and initializer
  - JSON tags and pointer semantics are consistent

