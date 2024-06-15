package handler

import (
	"log"
	"net/http"

	"h24s_14-extension-backend/user"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func HandleGetMe(c echo.Context) error {
	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}

	var me user.Me
	me.Name = sess.Values["my_name"].(string)
	me.IconBase64 = sess.Values["my_icon_base64"].(string)
	me.IsVisitor = sess.Values["is_visitor"].(bool)

	return c.JSON(http.StatusOK, &me)
}
