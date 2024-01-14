package data

type CartGame struct {
	CartId       int
	CartGameId   int
	ParentCartId int
	UserId       int
	GameId       int
	Title        string
	CoverImage   string
	Price        string
}

type Cart struct {
	GameId int
	UserId int
}
