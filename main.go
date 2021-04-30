package main

import (
	"github.com/urfave/cli"
	"github.com/ycoe/gorm-generator/generator"
	"os"
)

func main3() {
	generator.GenBaseDao("finance", "dao", nil, "")
}

func main() {
	app := cli.NewApp()
	app.Usage = "generate model for jinzhu/gorm"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1:3306",
			Usage: "mysql server,default;127.0.0.1:3306",
		},
		cli.StringFlag{
			Name:  "username, u",
			Value: "root",
			Usage: "Username of mysql",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Password of mysql",
		},
		cli.StringFlag{
			Name:  "database, d",
			Value: "",
			Usage: "select database",
		},
		cli.StringFlag{
			Name:  "table, t",
			Usage: "table name",
			Value: "ALL",
		},
		cli.StringFlag{
			Name:  "dir",
			Usage: "path which models will be stored",
			Value: "models",
		},
		cli.StringFlag{
			Name:  "daodir, dd",
			Usage: "path which dao will be stored",
			Value: "dao",
		},
		cli.StringFlag{
			Name:  "appid",
			Usage: "your appId, eg: helloworld",
		},
		cli.StringFlag{
			Name:  "tablePrefix, tp",
			Usage: "table prefix",
			Value: "",
		},
		cli.StringFlag{
			Name:  "dp",
			Usage: "file dao.go path",
			Value: "gitee.com/inngke/proto/common",
		},
	}
	app.Action = generator.Generate
	_ = app.Run(os.Args)
}
