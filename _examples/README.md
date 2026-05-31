# go-common-domain examples

These examples are small standalone programs for the value objects in this module.

Run any example with `go run` and an explicit path:

```sh
go run ./_examples/finance/money
go run ./_examples/finance/percentagerate
go run ./_examples/product/identifier
go run ./_examples/person/fullname
go run ./_examples/auth/password
```

Coverage by folder:

- `auth/password`: password validation, verification, protected string output, hash access, and reconstitution.
- `domain/error`: wrapped domain errors and `errors.Is`.
- `finance/currency`: currency normalization, minor-unit metadata, and reconstitution.
- `finance/money`: inferred currency scale, minor-unit money, arithmetic, rounding, JSON, relational reconstitution, and percentage rates.
- `finance/percentagerate`: basis-point rates, money application, and basis-point reconstitution.
- `geography/address`: postal address normalization, JSON, and relational reconstitution.
- `geography/countrycode`: country-code normalization.
- `geography/geopoint`: latitude/longitude validation and JSON.
- `identifier`: integer identifiers from numbers, strings, and reconstituted persisted values.
- `identifier/uuid`: UUID identifier normalization, JSON, and reconstitution.
- `person/fullname`: multi-part names, JSON, and relational reconstitution.
- `person/contact/phonenumber`: phone-number normalization.
- `product/identifier`: EAN, CNK, and GTIN product identifiers.
- `temporal/daterange`: date-only ranges, containment, JSON, and reconstitution.
- `temporal/timerange`: timestamp ranges, durations, JSON, and reconstitution.
- `web/domainname`: domain-name normalization.
- `web/email`: email normalization.
- `web/ipaddress`: IP normalization and IPv4/IPv6 helpers.
- `web/url`: URL normalization and parsed components.
