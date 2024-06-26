package pojo

import "time"

// 商品
type Item struct {
	ItemID    int       `json:"item_id" xorm:"'id' Double autoincr"`
	Name      string    `json:"name" xorm:"name"`
	Price     float64   `json:"price" xorm:"price"`
	IsActive  int       `json:"is_active" xorm:"is_active"`
	CreatedAt time.Time `json:"created_at" xorm:"created_at TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated_at TIMESTAMP"`
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted_at TIMESTAMP"`
}

// responseData
type ResponseData struct {
	ItemID int     `json:"item_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

// 增加商品信息请求参数
type AddItemReq struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// 增加商品信息响应参数
type AddItemResp struct {
	Item_info ResponseData `json:"item_info"`
}

// 修改商品信息请求参数
type UpdateItemReq struct {
	ItemID int     `json:"item_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

// 修改商品信息响应参数
type UpdateItemResp struct {
	StoreInfo ResponseData `json:"store_info"`
}

// 获取商品信息响应参数
type GetItemResp struct {
	StoreInfo ResponseData `json:"store_info"`
}

// 删除商品信息响应参数
type DeleteItemResp struct {
	DeleteTime string `json:"delete_time"`
}
