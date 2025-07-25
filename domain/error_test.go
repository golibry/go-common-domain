package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ErrorTestSuite contains tests for domain error functionality
type ErrorTestSuite struct {
	suite.Suite
}

// TestErrorSuite runs the error test suite
func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) TestItCanCreateNewError() {
	err := NewError("test error: %s", "message")

	s.Equal("test error: message", err.Error())
	s.Nil(err.Unwrap(), "Expected nil unwrap for non-wrapped error")
}

func (s *ErrorTestSuite) TestItCanCreateNewErrorWithWrap() {
	originalErr := errors.New("original error")
	wrappedErr := NewErrorWithWrap(originalErr, "wrapped: %s", "context")

	expectedMsg := "wrapped: context: original error"
	s.Equal(expectedMsg, wrappedErr.Error())
	s.Equal(originalErr, wrappedErr.Unwrap(), "Expected unwrap to return original error")
}

func (s *ErrorTestSuite) TestItIsCompatibleWithErrorsIs() {
	originalErr := errors.New("original error")
	wrappedErr := NewErrorWithWrap(originalErr, "wrapped error")

	// Test that errors.Is works with our wrapped error
	s.True(
		errors.Is(wrappedErr, originalErr),
		"errors.Is should find the original error in the wrapped error",
	)

	// Test with a different error
	differentErr := errors.New("different error")
	s.False(
		errors.Is(wrappedErr, differentErr),
		"errors.Is should not find a different error in the wrapped error",
	)
}

func (s *ErrorTestSuite) TestItIsCompatibleWithErrorsAs() {
	originalErr := NewError("original domain error")
	wrappedErr := NewErrorWithWrap(originalErr, "wrapped error")

	var domainErr *Error
	s.True(
		errors.As(wrappedErr, &domainErr),
		"errors.As should find the domain error in the wrapped error",
	)
	s.Equal(wrappedErr, domainErr, "errors.As should return the wrapping domain error")
}

func (s *ErrorTestSuite) TestItIsCompatibleWithErrorsUnwrap() {
	originalErr := errors.New("original error")
	wrappedErr := NewErrorWithWrap(originalErr, "wrapped error")

	unwrapped := errors.Unwrap(wrappedErr)
	s.Equal(originalErr, unwrapped, "errors.Unwrap should return the original error")

	// Test with non-wrapped error
	simpleErr := NewError("simple error")
	unwrappedSimple := errors.Unwrap(simpleErr)
	s.Nil(unwrappedSimple, "errors.Unwrap should return nil for non-wrapped error")
}

func (s *ErrorTestSuite) TestItCanHandleChainedWrapping() {
	baseErr := errors.New("base error")
	firstWrap := NewErrorWithWrap(baseErr, "first wrap")
	secondWrap := NewErrorWithWrap(firstWrap, "second wrap")

	// Test that errors.Is can find the base error through multiple wraps
	s.True(
		errors.Is(secondWrap, baseErr),
		"errors.Is should find the base error through multiple wraps",
	)
	s.True(errors.Is(secondWrap, firstWrap), "errors.Is should find the first wrap error")

	// Test the error message includes all levels
	expectedMsg := "second wrap: first wrap: base error"
	s.Equal(expectedMsg, secondWrap.Error())
}

func (s *ErrorTestSuite) TestItCanHandleMixedWrappingWithFmtErrorf() {
	domainErr := NewError("domain error")
	fmtWrapped := fmt.Errorf("fmt wrapped: %w", domainErr)
	domainWrapped := NewErrorWithWrap(fmtWrapped, "domain wrapped")

	// Test that errors.Is can find the original domain error
	s.True(
		errors.Is(domainWrapped, domainErr),
		"errors.Is should find the original domain error through mixed wrapping",
	)

	// Test that we can unwrap step by step
	firstUnwrap := errors.Unwrap(domainWrapped)
	s.Equal(fmtWrapped, firstUnwrap, "First unwrap should return the fmt wrapped error")

	secondUnwrap := errors.Unwrap(firstUnwrap)
	s.Equal(domainErr, secondUnwrap, "Second unwrap should return the original domain error")
}
