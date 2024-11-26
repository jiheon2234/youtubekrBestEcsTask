package main

import (
	"collectYoutubeData/config"
	"collectYoutubeData/service"
	"log"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {

	log.Println("HELLO WORLD!!")

	//now := time.Now()
	//
	//// yyyy-mm-dd 형식으로 출력
	//formattedDateTime := now.Format("2006-01-02 15:04:05")
	//fmt.Println(formattedDateTime)

	cfg := config.NewConfig()
	y := service.NewYoutubeService(cfg)
	y.Run()
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
