package account

import (
	"strings"
	"time"
)

type Premium struct {
	Start  int64 `json:"start"`
	Finish int64 `json:"finish"`
}
type Like struct {
	Ts int64 `json:"ts"`
	ID int64 `json:"id"`
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

func (a *Account) EmailDomain() string {
	if a.Email == "" {
		return ""
	}
	return strings.Split(a.Email, "@")[1]
}
func (a *Account) UTCBirthYear() int {
	return time.Unix(a.Birth, 0).UTC().Year()
}
func (a *Account) PhoneCode() string {
	if a.Phone == "" {
		return ""
	}
	l := strings.IndexByte(a.Phone, '(')
	r := strings.IndexByte(a.Phone, ')')
	return a.Phone[l+1 : r]
}
func (a *Account) LikesIds() []int64 {
	var ans []int64

	for _, l := range a.Likes {
		ans = append(ans, l.ID)
	}

	return ans
}
