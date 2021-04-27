package generator

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/ycoe/gorm-generator/database"
	"strings"
)

func Generate(c *cli.Context) error {
	dbSns := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		c.String("u"), c.String("p"), c.String("d"))
	db := database.GetDB(dbSns)
	appId := c.String("appid")
	daoDir := c.String("daodir")
	daoPackage := c.String("dp")
	tableName := c.String("t")
	tablePrefix := c.String("tablePrefix")
	if tableName == "ALL" {
		tableNames := make([]string, 0)
		tables := db.GetDataBySql("show tables")
		for _, table := range tables {
			orgTableName := table["Tables_in_"+c.String("d")]
			tableName := getTableName(orgTableName, tablePrefix)
			tableNames = append(tableNames, tableName)
			columns := db.GetDataBySql("SHOW FULL COLUMNS FROM  " + orgTableName)
			idType := GenerateModel(tableName, columns, c.String("dir"))
			GenerateDao(orgTableName, appId, tableName, daoDir, idType, daoPackage)
		}

		//生成dao.go
		//index := strings.LastIndex(daoDir, "/")
		//daoPackage := daoDir[index+1:]
		//GenBaseDao(appId, daoPackage, tableNames, tablePrefix)
	} else {
		for _, table := range strings.Split(tableName, ",") {
			orgTableName := table
			columns := db.GetDataBySql("desc " + tableName)
			tableName := getTableName(tableName, tablePrefix)
			idType := GenerateModel(tableName, columns, c.String("dir"))
			GenerateDao(orgTableName, appId, tableName, daoDir, idType, daoPackage)
		}
	}
	return nil
}

func getTableName(orgTableName, tablePrefix string) string {
	if len(tablePrefix) == 0 {
		return orgTableName
	}
	index := strings.LastIndex(orgTableName, tablePrefix)
	if index == 0 {
		return orgTableName[len(tablePrefix):]
	}
	//fmt.Printf("当前表名%s 不是以前缀%s 开头！\n", orgTableName, tablePrefix)
	return orgTableName
}
