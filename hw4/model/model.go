package model

const (
	AK47_CASE_HARDENED = "https://steamcommunity.com/market/listings/730/AK-47%20%7C%20Case%20Hardened%20%28Factory%20New%29"
	AK47_BLOODSPORT = "https://steamcommunity.com/market/listings/730/AK-47%20%7C%20Bloodsport%20%28Factory%20New%29"
	AK47_AQUAMARINE = "https://steamcommunity.com/market/listings/730/AK-47%20%7C%20Aquamarine%20Revenge%20%28Factory%20New%29"
)

type Skin struct {
	Name string
	Price string
}

type NamingWithUrl struct {
	Name string
	Url string
}