//go:generate gorunpkg github.com/99designs/gqlgen --verbose

package graph

import (
	"github.com/charlysan/goneo4jgql/internal/app/service"
)

// Resolver is the main gql resolver
type Resolver struct {
	Service service.Service
}
