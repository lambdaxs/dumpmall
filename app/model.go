package app

import "encoding/json"

type Order struct {
    ID int64 `json:"id"`
    AID int64 `json:"aid" gorm:"column:aid"'`
    UserName string `json:"user_name"`
    Phone string `json:"phone"`
    Address string `json:"address"`
    Items []Item `json:"items" gorm:"-"`
    ItemInfo string `json:"item_info"`
    Status string `json:"status"`
    Timestamp int64 `json:"timestamp"`
}

//出库
func (o *Order)Decode() {
    _ = json.Unmarshal([]byte(o.ItemInfo), &o.Items)
    o.ItemInfo = ""
}

type Item struct {
    ID int64 `json:"id"`
    AID int64 `json:"aid" gorm:"column:aid"'`
    Name string `json:"name"`
    Desc string `json:"desc"`
    Price int64 `json:"price"`
    Count int64 `json:"count"` // 购买数量
    Number int64 `json:"number"` //库存
    ImageUrl string `json:"image_url"`
    Show string `json:"show"` //上下架
}