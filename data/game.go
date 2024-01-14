package data

type Image struct {
	ImageId int
	Medium  string
	Small   string
	Large   string
}

type Game struct {
	CoverImage  string
	Images      []Image
	Name        string
	Description string
	Price       string
	GameId      int
}
