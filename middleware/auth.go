package middleware

import (
	"log"
	"net/http"

	"h24s_14-extension-backend/user"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SetUserInformationToSess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("note_session", c)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "Failed to get the session.")
		}

		accessToken := sess.Values["access_token"]
		var myInformation user.Me
		if accessToken != nil {
			myInformation = user.GetUserInformation(accessToken.(string))
		} else {
			myInformation = user.GetUserInformation("")
		}

		sess.Values["my_name"] = myInformation.Name
		sess.Values["my_icon_base64"] = myInformation.IconBase64
		sess.Values["is_visitor"] = myInformation.IsVisitor

		//isVisitorそもそも必要なさそうなので後できれいにする
		if myInformation.IsVisitor {
			return c.String(http.StatusForbidden, "auth failed")
		}

		return next(c)
	}
}
