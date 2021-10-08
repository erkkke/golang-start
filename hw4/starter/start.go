package starter

import (
	"context"
	// "encoding/json"
	"fmt"
	model "hw/4/model"
	// "os"
	"time"
)



func Run() {
	skinList := map[int]model.NamingWithUrl {
		1 : {Name: "Ak-47 | Case Hardened", Url: model.AK47_CASE_HARDENED},
		2 : {Name: "Ak-47 | Bloodsport", Url: model.AK47_BLOODSPORT},
		3 : {Name: "Ak-47 | Aquamarine", Url: model.AK47_AQUAMARINE},
	}
	ctx := context.Background()
	

	fmt.Println("Welcome to cs:go skins price parser!\n\nPlease, choose the skin that you want to monitor the price")
	chooseSkin(ctx, skinList)
}


func chooseSkin(ctx context.Context,skinList map[int]model.NamingWithUrl) {
	
	for {
		fmt.Println("[0] Quit")
		for k, v := range skinList {
			fmt.Printf("[%v] %v\n", k, v.Name)
		}
		var input int
		fmt.Scanln(&input)
		if input > 0 && input <= len(skinList) {
			printPrice(ctx, skinList, input)
		} else if input == 0 {
			return
		}
		
	}
	
}


func printPrice(ctx context.Context, skinList map[int]model.NamingWithUrl, input int) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		fmt.Scanln()
		cancel()
	}()

	fmt.Println("Type anything to go back")
	for {
		select {
		case <- time.After(5 * time.Second):
			prices := GetPrice(skinList[input].Name, skinList[input].Url)
			// formatted, err := json.Marshal(prices)
			// if err != nil {
			// 	fmt.Println("error:", err)
			// }
			// os.Stdout.Write(formatted)
			for i, val := range prices {
				fmt.Printf("%v. %v | Price: %v\n", i + 1, val.Name, val.Price)
			}
			fmt.Println()
			
		case <-ctx.Done():
			return
		}
	}
}