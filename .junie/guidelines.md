## Overview
The project's purpose is to provide common domain concepts used in many applications. For 
example, the full name or the domain name are widely accepted concepts.

## Testing Practices

### Test Organization
- Use table-driven tests for multiple test cases
- Use testify/suite for all tests
- Group related tests in test suites
- Use descriptive test names that explain the scenario
- Test package behavior and follow the test name convention like TestItCanDoSomething
- Maximize the return on investment to ease maintenance
- Do not go more than once through a tested logic path 

### Error Testing
- Test both success and failure scenarios
- Use `errors.Is()` and `errors.As()` in tests
- Test error wrapping and unwrapping behavior
- Verify error messages are meaningful

### Test Coverage
- Test all public methods and functions
- Test edge cases and boundary conditions
- Test error conditions and invalid inputs
- Test JSON marshaling/unmarshaling if applicable

## Code Documentation

### Function Documentation
- Document all exported functions and methods
- Start with the function name
- Explain what the function does, not how it does it
- Document parameters and return values when not obvious

```go
// NewFullName creates a new instance of FullName with validation and normalization
func NewFullName(firstName, middleName, lastName string) (FullName, error) {
    // implementation
}
```

### Comment Style
- Use complete sentences in comments
- Keep comments concise but informative
- Update comments when code changes
- Avoid obvious comments that just repeat the code

## Validation and Normalization

### Input Validation
- Validate all inputs at domain boundaries
- Use separate validation functions for reusability
- Return meaningful error messages
- Handle Unicode properly with `unicode/utf8` package

### Data Normalization
- Normalize data consistently across the application
- Use `strings.Builder` for efficient string construction
- Handle edge cases in normalization (multiple spaces, special characters)
- Separate normalization from validation

## General Go Best Practices

### Code Organization
- Keep functions small and focused on a single responsibility
- Use early returns to reduce nesting
- Group related functionality together
- Separate concerns clearly
- Do not use both value and pointer receivers in struct methods. Use one or the other

### Concurrency
- Follow Go's concurrency patterns: "Don't communicate by sharing memory; share memory by communicating"
- Use channels for communication between goroutines
- Use sync package primitives when appropriate
- Always handle goroutine lifecycle properly

### Dependencies
- Keep dependencies minimal and well-justified
- Prefer standard library when possible
- Use semantic versioning for your modules
- Regularly update dependencies for security patches

### Code Style
- Use `gofmt` to format code consistently
- Use `golint` and `go vet` for code quality checks
- Follow the Go Code Review Comments guidelines
- Use meaningful variable names, even if they're longer