package main

// paht to aps modified to use my implementation
import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/iris-contrib/gothic"
	"github.com/kataras/iris"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"github.com/silviofdez/oauth/client/goth/aps"
)

// UserNotLoggedIn - Return in Json format user not logged message
func UserNotLoggedIn() (int, string) {
	error, _ := json.Marshal(map[string]string{"code": "401", "message": "Please login into the platform"})
	return 401, string(error)
}

// StandardErrorWithStatusCode - Return in Json format standar error message
func StandardErrorWithStatusCode(code int, errorMessage error) (int, string) {
	error, _ := json.Marshal(map[string]string{"code": iris.StatusText(code), "message": errorMessage.Error()})
	return code, string(error)
}

func main() {

	api := iris.New()

	goth.UseProviders(
		aps.New("bawdy-reindeers-14-56dd2bcc2ba94", "6454acedc7024fdfa743c5407da7ad44", "http://localhost:3000/auth/aps/callback"),
		gplus.New("72983246488-upvsod3t92stf9o9ojvqvqrip0t3anln.apps.googleusercontent.com", "M8D_euTcQ9WC2NJdTwVqwX5R", "http://localhost:3000/auth/gplus/callback"),
	)

	api.Get("/", func(ctx *iris.Context) {
		user := ctx.Session().Get("user")
		if user == nil {
			ctx.JSON(UserNotLoggedIn())
		} else {
			ctx.JSON(200, ctx.Session().Get("user").(goth.User))
		}
	})

	api.Get("/auth/:provider", func(ctx *iris.Context) {
		err := gothic.BeginAuthHandler(ctx)
		if err != nil {
			ctx.JSON(StandardErrorWithStatusCode(iris.StatusInternalServerError, err))
			return
		}
	})

	api.Get("/auth/:provider/callback", func(ctx *iris.Context) {
		user, err := gothic.CompleteUserAuth(ctx)
		if err != nil {
			ctx.JSON(StandardErrorWithStatusCode(iris.StatusUnauthorized, err))
			return
		}
		ctx.Session().Set("user", user)
		ctx.Redirect("/", iris.StatusOK)
	})

	api.Get("/login", func(c *iris.Context) {
		c.Write("You are session with ID: %s", c.Session().ID())
	})

	api.Get("/logout", func(ctx *iris.Context) {
		//destroy, removes the entire session and cookie
		ctx.SessionDestroy()
		ctx.Redirect("/", iris.StatusAccepted)
	})

	w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

	fmt.Fprintf(w, "Configured Routes:\n")
	fmt.Fprintf(w, "\tNAME\tMETHOD\tPATH\n")

	for _, route := range api.Lookups() {
		fmt.Fprintf(w, "\t%s\t%s\t%s\n", route.Name(), route.Method(), route.Path())
	}
	w.Flush()

	api.Listen(":3000")

}
