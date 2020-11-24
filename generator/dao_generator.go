package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/ycoe/gorm-generator/helper"
	"os"
	"strings"
)

/**
package dao

import model "finance/model"

type AccountDao struct {
	dao *Dao
}

func (account *AccountDao) Create (entity *model.Account) (uint, error) {
	result := account.dao.client.Table("accounts").Create(entity)
	return entity.ID, result.Error
}
*/
func GenerateDao(orgTableName string, appId, tableName, dir string) {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by model-generator.")
	f.ImportAlias(appId+"/model", "model")
	f.ImportAlias(appId+"/proto", appId)
	f.ImportAlias("time", "time")

	genDaoVar(f, orgTableName, tableName)
	genGetEntityDao(f, orgTableName, tableName)
	genEntityDaoStruct(f, tableName)
	genCreateFun(f, appId, tableName)

	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".dao.go"
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(fileName)
}

/**
var accountDao *AccountDao = nil
*/
func genDaoVar(f *jen.File, orgTableName string, tableName string) {
	tableDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	tableDaoVarName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), false) + "Dao"
	f.Var().Id(tableDaoVarName).Id("*" + tableDaoName).Op("=").Nil()
}

/**
func GetAccountDao() *AccountDao {
	err := GetDao().Ping()
	if err != nil {
		logger.Error(err)
	}
	if accountDao != nil && err == nil {
		return accountDao
	}
	table := GetDao().client.Table("finance_account")
	accountDao = &AccountDao{
		db: table,
	}

	return accountDao
}
*/
func genGetEntityDao(f *jen.File, orgTableName string, tableName string) {
	tableDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	tableDaoVarName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), false) + "Dao"
	f.Func().Id("Get"+tableDaoName).Params().Id("*"+tableDaoName).Block(
		jen.Id("err").Op(":=").Id("GetDao").Call().Dot("Ping").Call(),
		jen.If(
			jen.Id("err").Op("!=").Nil().Block(
				jen.Qual("github.com/micro/go-micro/v2/logger", "Error").Call(
					jen.Id("err"),
				),
			),
		),
		jen.If(
			jen.Id(tableDaoVarName).Op("!=").Nil().Op("&&").Id("err").Op("==").Nil().Block(
				jen.Return(
					jen.Id(tableDaoVarName),
				),
			),
		),
		jen.Id("table").Op(":=").Id("GetDao").Call().Dot("client").Dot("Table").Call(
			jen.Lit(orgTableName),
		),
		jen.Id(tableDaoVarName).Op("=").Id("&"+tableDaoName).Values(
			jen.Dict{
				jen.Id("db"): jen.Id("table"),
			},
		),
		jen.Return(jen.Id(tableDaoVarName)),
	)
}

/**
type AccountDao struct {
	db *gorm.DB
}
*/
func genEntityDaoStruct(f *jen.File, tableName string) {
	tableEntityDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Type().Id(tableEntityDaoName).Struct(
		jen.Id("db").Id("*").Qual("gorm.io/gorm", "DB"),
	)
}

/**
func (dao *AccountDao) Create(entity *model.Account) (uint, error) {
	result := d.db.Create(entity)
	return entity.ID, result.Error
}
*/
func genCreateFun(f *jen.File, appId, tableName string) {
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	entityDaoName := entityName + "Dao"
	f.Func().Params(
		jen.Id("d").Id("*"+entityDaoName),
	).Id("Create").Params(
		jen.Id("entity").Id("*").Qual(appId+"/model", entityName),
	).Params(
		jen.Id("uint"),
		jen.Id("error"),
	).Block(
		jen.Id("result").Op(":=").Id("d").Dot("db").Dot("Create").Call(
			jen.Id("entity"),
		),
		jen.Return(
			jen.Id("entity").Dot("ID"),
			jen.Id("result").Dot("Error"),
		),
	)
}