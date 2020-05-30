package account

type Premium struct {
	Start  int `json:"start"`
	Finish int `json:"finish"`
}
type Like struct {
	Ts int `json:"ts"`
	ID int `json:"id"`
}

type Account struct {
	ID        int64    `json:"id"`
	Fname     string   `json:"fname"`
	Sname     string   `json:"Sname"`
	Email     string   `json:"email"`
	Interests []string `json:"interests"`
	Status    string   `json:"status"`
	Premium   Premium  `json:"premium"`
	Sex       string   `json:"sex"`
	Phone     string   `json:"phone"`
	Likes     []Like   `json:"likes"`
	Birth     int64    `json:"birth"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Joined    int64    `json:"joined"`
}
