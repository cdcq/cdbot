package main

import (
	"cdbot/src/global"
	"cdbot/src/initialize"
	"log"
)

func main() {
	if err := initialize.InitConfig(); err != nil {
		log.Fatalln("Could not init configuration, exit...")
	}
	global.LOGGER.Info("\n" +
		"===================================================================================\n" +
		" cdbot \n" +
		"===================================================================================")
	global.DATABASE = initialize.GormMysql()
	initialize.InitTables(global.DATABASE)
	engine, err := initialize.InitRouter()
	if err != nil {
		log.Fatalln("Could not init router, exit...")
	}
	err = engine.Run(":8080")
	if err != nil {
		log.Fatalln("Could not run server! exit...")
	}
}
