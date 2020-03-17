package main

import (
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/server"
    "github.com/lambdaxs/dumpmall/app"
    "time"
)


func main() {
    httpServer := server.HttpServer{
        Host:        "127.0.0.1",
        Port:        9001,
    }

    //清理任务
    timer := time.NewTicker(time.Second*60*5)
    go func() {
        for range timer.C {
            app.ClearTask()
        }
    }()

    httpServer.StartEchoServer(func(srv *echo.Echo) {
        srv.POST("/api/order_list", app.OrderList) //订单列表
        srv.POST("/api/submit_order", app.SubmitOrder) //提交订单

        //admin
        srv.POST("/api/complete_order", app.CompleteOrder) //完成订单
        srv.POST("/api/item_list", app.ItemList) //商品列表
        srv.POST("/api/item_update", app.ItemUpdate) //编辑商品
        srv.POST("/api/share_order_text", app.ShareOrderText) //分享订单列表
    })
}