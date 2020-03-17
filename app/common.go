package app

import "github.com/labstack/echo"

type Response struct {
    Code int `json:"code"`
    Data interface{} `json:"data"`
    Message string `json:"message"`
}

func OutputData(c echo.Context, code int, data interface{}) error {
    c.Response().Header().Set("Access-Control-Allow-Origin", "*")
    return c.JSON(200, Response{
        Code: code,
        Data: data,
        Message:"",
    })
}

func OutputError(c echo.Context, code int, err error) error {
    c.Response().Header().Set("Access-Control-Allow-Origin", "*")
    return c.JSON(200, Response{
        Code: code,
        Data: nil,
        Message:err.Error(),
    })
}
