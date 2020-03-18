package app

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/jinzhu/gorm"
    "github.com/lambdaxs/go-server/driver/mysql_client"
    "time"
)

var DB *gorm.DB

func init() {
    sqlConf := mysql_client.MysqlDB{
        DSN:          "root:123456@tcp(127.0.0.1:3306)/mall?charset=utf8&parseTime=True&loc=Local",
        Log:          true,
    }
    db,err := sqlConf.ConnectGORMDB()
    if err != nil {
        panic(err)
    }
    DB = db

}

//订单列表
func OrderList(c echo.Context) error {
    reqModel := new(struct{
        AID int64 `json:"aid"`
        Page int64 `json:"page"`
        Status string `json:"status"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }

    limit := 10
    offset := (reqModel.Page - 1)*int64(limit)
    list := make([]Order, 0)
    sql := DB.Table("order").Where("aid = ? AND status != ?", reqModel.AID, "回收")
    if reqModel.Status != "" {
       sql = sql.Where("status = ?", reqModel.Status)
    }
    sql.Offset(offset).Limit(limit).Find(&list)
    for i,_ := range list {
        list[i].Decode()
    }
    return OutputData(c, 0, list)
}

//UserName string `json:"user_name"`
//Phone string `json:"phone"`
//Address string `json:"address"`
//Items []Item `json:"items" gorm:"-"`
//ItemInfo string `json:"item_info"`
//Status string `json:"status"`
//Timestamp int64 `json:"timestamp"`

//提交订单
func SubmitOrder(c echo.Context) error {
    reqModel := new(struct{
        AID int64 `json:"aid"`
        UserName string `json:"user_name"`
        Phone string `json:"phone"`
        Address string `json:"address"`
        ItemInfo string `json:"item_info" form:"item_info"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    data := &Order{
       AID:       reqModel.AID,
       UserName:  reqModel.UserName,
       Phone:     reqModel.Phone,
       Address:   reqModel.Address,
       ItemInfo: reqModel.ItemInfo,
       Status:    "已提交",
       Timestamp: time.Now().Unix(),
    }
    DB.Table("order").Create(data)
    return OutputData(c, 0, true)
}

//查询订单
func SearchOrder(c echo.Context) error {
    reqModel := new(struct{
        Aid int64 `json:"aid" form:"aid"`
        Phone string `json:"phone" form:"phone"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    list := make([]Order, 0)
    DB.Table("order").Where("aid = ? AND phone = ?", reqModel.Aid, reqModel.Phone).Find(&list)
    return OutputData(c, 0, list)
}

//完成订单
func CompleteOrder(c echo.Context) error  {
    reqModel := new(struct{
        OrderID int64 `json:"order_id"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    DB.Table("order").Where("id = ?", reqModel.OrderID).Update("status", "已完成")
    return OutputData(c, 0, true)
}

//商品列表展示
func ItemList(c echo.Context) error {
    reqModel := new(struct{
        AID int64 `json:"aid"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    list := make([]Item, 0)
    DB.Table("item").Where("aid = ?", reqModel.AID).Find(&list)
    return OutputData(c, 0, list)
}

//商品信息更新
func ItemUpdate(c echo.Context) error {
    reqModel := new(struct{
        ItemID int64 `json:"item_id"`
        Name string `json:"name"`
        Desc string `json:"desc"`
        Price int64 `json:"price"`
        Number int64 `json:"number"`
        ImageUrl string `json:"image_url"`
        Show string `json:"show"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    DB.Table("item").Where("id = ?", reqModel.ItemID).Update(map[string]interface{}{
        "name": reqModel.Name,
        "desc": reqModel.Desc,
        "price": reqModel.Price,
        "number": reqModel.Number,
        "image_url": reqModel.ImageUrl,
        "show": reqModel.Show,
    })
    return OutputData(c, 0, true)
}

//分享订单列表
func ShareOrderText(c echo.Context) error {
    reqModel := new(struct{
        AID int64 `json:"aid"`
        IDS []int64 `json:"ids"`
    })
    if err := c.Bind(reqModel);err != nil {
        return OutputError(c, 1, fmt.Errorf("请求参数错误:%s", err.Error()))
    }
    sql := DB.Table("order").Where("aid = ? AND status != ?", reqModel.AID, "回收")
    if len(reqModel.IDS) > 0 {
        sql = sql.Where("id IN (?)", reqModel.IDS)
    }
    list := make([]Order, 0)
    sql.Find(&list)
    return OutputData(c, 0 , list)
}

//订单过期任务
func ClearTask() {
    list := make([]Order, 0)
    DB.Table("order").Where("status = ?", "已完成").Find(&list)

    now := time.Now().Unix()
    for _,item := range list {
        //30min
        if now - item.Timestamp > 60*30 {
            DB.Table("order").Where("id = ?", item.ID).Update("status", "回收")
        }
    }
}

