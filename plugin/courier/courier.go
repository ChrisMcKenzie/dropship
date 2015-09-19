package courier

import (
	"net/http"

	"github.com/ChrisMcKenzie/dropship/model"
	"github.com/gin-gonic/gin"
)

type (
	Courier interface {
		Authorize(*gin.Context) (*model.Authentication, error)
		GetKind() string
		ParseHook(*http.Request) (*model.Deployment, error)
		GetRepos(user *model.User) ([]*model.Repo, error)
		GetScript(user *model.User, repo *model.Repo, deployment *model.Deployment) ([]byte, error)
		Activate(repo *model.Repo, hook string) error
	}
)

// List of registered plugins.
var couriers []Courier

// Register registers a plugin by name.
//
// All plugins must be registered when the application
// initializes. This should not be invoked while the application
// is running, and is not thread safe.
func Register(courier Courier) {
	couriers = append(couriers, courier)
}

// List Registered remote plugins
func Registered() []Courier {
	return couriers
}

// Lookup gets a plugin by name.
func Lookup(name string) Courier {
	for _, courier := range couriers {
		if courier.GetKind() == name {
			return courier
		}
	}
	return nil
}
