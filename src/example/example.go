// This is an example plugin that will be available to www.yoursite.com/misc
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/student020341/go-server-core/src/lib/RouterModule"
)

/*
	HandleWeb is expected to be implemented by every plugin
	void function to respond directly to the client
	most can be handled by the included router class's Handle method.
	path is the remaining path after the application name is removed,
	ex: www.site/com/misc/hello will be ["hello"]
*/

func HandleWeb(w http.ResponseWriter, r *http.Request, path []string) {

	router.Handle(w, r, path)
}

// GetName is expected to be implemented by every plugin. This is the root of your sub application.
// visiting www.site.com/{GetName()} will enter this plugin's HandleWeb function
func GetName() string {
	return "misc"
}

// router included with server core
var router RouterModule.SubRouter

// use init to setup your router
func init() {
	// www.site.com/misc/file/anything.asdf
	// shows route that takes a writer, a request, and args - expected to respond to client directly
	router.Register("/file/*", "GET", func(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
		// will attempt to serve the given path ex: project-root/files/misc/anything.asdf
		// r.URL.Path[11:] turns /misc/file/anything.asdf into anything.asdf
		// files/misc/hello.html will be included with the base repo, visit www.yoursite.com/misc/hello.html to check it out
		http.ServeFile(w, r, "./files/misc/"+r.URL.Path[11:])
	})

	// www.site.com/misc/code/200
	// shows a route that takes args only - expected to return a json-like interface to be returned to client as json
	router.Register("/code/:code", "*", func(args map[string]interface{}) interface{} {
		// get route arguments
		route := args["route"].(map[string]string)
		// get :code from route arguments
		status, err := strconv.Atoi(route["code"])

		var code int
		var msg string
		if err != nil {
			code = 500
			msg = err.Error()
		} else {
			code = status
			msg = "testing status code"
		}

		// special arg HTTPStatusCode will overwrite the status code returned by the included router
		return map[string]interface{}{
			"HTTPStatusCode": code,
			"status":         msg,
		}
	})
}

// use main to test your application outside of the server environment
func main() {
	fmt.Println("hello from example")
}