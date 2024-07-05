package orderTracking

type OrderList struct {
	Id          int    `json:"id"`
	Title       int    `json:"title"`
	Description string `json:"description"`
}

type UserList struct {
	Id     int
	UserId int
	ListId int
}

type OrderItem struct {
	Id     int `json:"id"`
	ListId int `json:"listId"`
	ItemId int `json:"itemId"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}
