package facades

import (
	"github.com/goravel/framework/contracts/http/client"
)

// Http returns the HTTP client factory. The scaffold returned the narrower
// client.Request interface; the factory also has Fake() and the Assert*
// methods that tests need. See docs/chapters/21-http-client.md.
func Http() client.Factory {
	return App().MakeHttp()
}
