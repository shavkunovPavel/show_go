package main

import (
	"admigo/common"
	c "admigo/controllers"
	"admigo/controllers/api"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func getRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir(common.Env().Static))
	router.GET("/", c.Index)
	router.GET("/users", c.Users)
	router.GET("/items", c.Items)
	router.GET("/shopping", c.Shopping)
	router.GET("/ico", c.Ico)
	router.GET("/btc", c.ChartBTC)
	router.GET("/eth", c.ChartETH)
	router.GET("/person", c.Person)
	router.GET("/app", c.App)

	router.POST("/login", c.Login)
	router.POST("/logout", c.Logout)
	router.POST("/signup", c.SignupAccount)
	router.GET("/confirm/*filepath", c.ConfirmUser)

	router.POST("/users", api.UsersList)
	router.GET("/users/:id", api.User)
	router.DELETE("/users/:id", api.DeleteUser)
	router.POST("/users/update", api.EditUser)

	router.GET("/roles", api.RolesList)

	router.POST("/items", api.ItemsList)
	router.GET("/items/:id", api.Item)
	router.POST("/items/update", api.UserCanDo(api.EditItem))
	router.DELETE("/items/:id", api.UserCanDo(api.DeleteItem))

	router.GET("/eth-prices", api.EthPrices)
	return
}
