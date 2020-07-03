package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"restaurant_backend/controllers"
)

func main() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
		}
	}
	r := httprouter.New()
	uc := controller.NewUserController(getSession())
	r.ServeFiles("/static/*filepath", http.Dir("/"))
	r.GET("/",uc.Starter)
	r.POST("/create",uc.MobileCreateUser)
	r.POST("/dishcreate",uc.DishCreate)
	r.GET("/getdishes",uc.GetDishes)
	r.GET("/placeorder",uc.PlaceOrder)
	r.POST("/login",uc.MobileLogin)
	r.GET("/image",uc.GetImage)
	r.POST("/createcomments",uc.CreateComments)
	r.GET("/getcomments",uc.GetComments)
	r.GET("/promo",uc.Promo)
	r.GET("/leader",uc.Leader)
	r.POST("/createleader",uc.CreateLeader)
	r.POST("/feedback",uc.FeedBack)
	r.POST("/notification",uc.NotifyUser)
	r.GET("/pastorder",uc.PastOrders)
	r.GET("/lastorder",uc.LastOrder)
	http.ListenAndServe(":8080", r)
}
func getSession() *sql.DB{
	s, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	return s
}