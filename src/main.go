package main

import (
	//"fmt"
	"dbOps"
	"tgHandler"
)



func main() {
	dbOps.Hehe()
	bot := tgHandler.BotInit()
	dbOps.NewDb()  
	
	
	bot.SimpleStart()
}
