package tokenutil

// TokenGenerator generates unique token for each request.
// It should be abstracted away as a single service.
// It should be distributed in order to make our backend system scalable.
type TokenGenerator interface {
	// Generate returns an unique token from this generator.
	// If no token could be generated before timeout, it returns error.
	// tokens are isolated in different namespace
	Generate(namespace string) (string, error)
}
