package starter

import (
	model "hw/4/model"
	"log"
	"net/http"
	"strings"
	"golang.org/x/net/html"
)

func GetPrice(name, url string) []model.Skin {
	var skins []model.Skin
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	
	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	result := ParsePrice(nil, doc)

	for _, price := range result {
		skins = append(skins, model.Skin{Name: name, Price: price})
	}

	return skins
}

func ParsePrice(s []string, n *html.Node) []string {
	if len(s) >= 3 {return s}
	className := "market_listing_price market_listing_price_with_fee"
	if n.Type == html.ElementNode && n.Data == "span" {
		for _, res := range n.Attr {
			if res.Key == "class" && res.Val == className{
				price := strings.ReplaceAll(n.FirstChild.Data, "	", "")
				price = strings.ReplaceAll(price, "\n", "")
				s = append(s, price)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
        s = ParsePrice(s, c)
    }

	return s
}

