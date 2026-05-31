# go-common-domain examples

These examples are small standalone programs for the value objects in this module.

Run any example with `go run` and an explicit path:

```sh
go run ./_examples/finance/money
go run ./_examples/finance/percentagerate
go run ./_examples/person/fullname
go run ./_examples/auth/password
```

Coverage by folder:

- `auth/password`: password validation, verification, protected string output, and database hash scan/value.
- `domain/error`: wrapped domain errors and `errors.Is`.
- `finance/currency`: currency normalization, minor-unit metadata, and database scan.
- `finance/money`: inferred currency scale, minor-unit money, arithmetic, rounding, JSON, relational reconstitution, and percentage rates.
- `finance/percentagerate`: basis-point rates, money application, and canonical basis-point database scan/value.
- `geography/countrycode`: country-code normalization.
- `identifier`: integer identifiers from numbers, strings, and database values.
- `person/fullname`: multi-part names, JSON, and relational reconstitution.
- `person/contact/phonenumber`: phone-number normalization.
- `web/domainname`: domain-name normalization.
- `web/email`: email normalization.
- `web/ipaddress`: IP normalization and IPv4/IPv6 helpers.
- `web/url`: URL normalization and parsed components.
