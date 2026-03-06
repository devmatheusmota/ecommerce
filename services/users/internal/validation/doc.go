// Package validation provides input validation for the users service.
//
// Input validation (required fields, format) runs in handlers, not in use cases.
// Handlers call ValidateRegisterInput, ValidateLoginInput, etc. and return 400
// before delegating to the use case. Use cases contain domain/business logic only.
//
// Reuse policy: Rules used by more than one flow (e.g. email format for register
// and login) MUST live in common.go and be reused. Do not duplicate the same
// checks in multiple validators — add shared helpers in common.go and call them
// from ValidateRegisterInput, ValidateLoginInput, etc.
package validation
