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
func GenerateDao(orgTableName string, appId, tableName, dir, idType string) {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by model-generator.")
	f.ImportAlias(appId+"/model", "model")
	f.ImportAlias(appId+"/proto", appId)
	f.ImportAlias("time", "time")

	genGetEntityDao(f, orgTableName, tableName)
	genEntityDaoStruct(f, tableName)
	genGetDb(f, tableName, orgTableName)
	genCreateFun(f, appId, tableName, idType)
	genGetByIdFun(f, appId, tableName, idType)

	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".dao.go"
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(fileName)
}

/**
func GetUserFundDao(tx ...*gorm.DB) *UserFundDao {
	if len(tx) == 0 {
		err := GetDao().Ping()
		if err != nil {
			logger.Error(err)
		}
		return &UserFundDao{
			client: GetDao().client,
		}
	} else {
		return &UserFundDao{
			client: tx[0],
		}
	}
}
*/
func genGetEntityDao(f *jen.File, orgTableName string, tableName string) {
	tableDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Func().Id("Get"+tableDaoName).Params(
		jen.Id("tx").Id("...*").Qual("gorm.io/gorm", "DB"),
	).Id("*"+tableDaoName).Block(
		jen.If(
			jen.Id("len").Call(jen.Id("tx")).Op("==").Id("0"),
		).Block(
			jen.Id("err").Op(":=").Id("GetDao").Call().Dot("Ping").Call(),
			jen.If(
				jen.Id("err").Op("!=").Nil().Block(
					jen.Qual("github.com/micro/go-micro/v2/logger", "Error").Call(
						jen.Id("err"),
					),
				),
			),
			jen.Return(
				jen.Id("&"+tableDaoName).Values(
					jen.Dict{
						jen.Id("client"): jen.Id("GetDao").Call().Dot("client"),
					},
				),
			),
		).Else().Block(
			jen.Return(
				jen.Id("&"+tableDaoName).Values(
					jen.Dict{
						jen.Id("client"): jen.Id("tx[0]"),
					},
				),
			),
		),
	)
}

/**
type AccountDao struct {
	client *gorm.DB
}
*/
func genEntityDaoStruct(f *jen.File, tableName string) {
	tableEntityDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Type().Id(tableEntityDaoName).Struct(
		jen.Id("client").Id("*").Qual("gorm.io/gorm", "DB"),
	)
}

/**
func (d *UserFundDao) GetDb() *gorm.DB {
	return d.client.Table("f_user_fund")
}
 */
func genGetDb(f *jen.File, tableName string, orgTableName string) {
	tableEntityDaoName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true) + "Dao"
	f.Func().Params(
		jen.Id("d").Id("*" + tableEntityDaoName),
	).Id("GetDb").Params().Id("*gorm.DB").Block(
		jen.Return(
			jen.Id("d").Dot("client").Dot("Table").Call(
				jen.Lit(orgTableName),
			),
		),
	)
}

/**
//通过ID获取
func (d *UserFundDao) GetById(id int32, fields ...string) (userFund model.UserFund, err error) {
	db := d.GetDb()
	if len(fields) > 0 {
		db = db.Select(fields)
	}

	err = db.First(&userFund, id).Error
	return
}
*/
func genGetByIdFun(f *jen.File, appId, tableName, idType string) {
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	entityVarName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), false)
	entityDaoName := entityName + "Dao"

	f.Comment("通过ID获取").Line().Func().Params(
		jen.Id("d").Id("*"+entityDaoName),
	).Id("GetById").Params(
		jen.Id("id").Id(idType),
		jen.Id("fields").Id("...string"),
	).Params(
		jen.Id(entityVarName).Id("model").Dot(entityName),
		jen.Id("err").Id("error"),
	).Block(
		jen.Id("db").Op(":=").Id("d").Dot("GetDb").Call(),
		jen.If(
			jen.Id("len").Call(
				jen.Id("fields"),
			).Op(">").Id("0"),
		).Block(
			jen.Id("db").Op("=").Id("db").Dot("Select").Call(
				jen.Id("fields"),
			),
		),
		jen.Id("err").Op("=").Id("db").Dot("First").Call(
			jen.Id("&").Id(entityVarName),
			jen.Lit("id=?"),
			jen.Id("id"),
		).Dot("Error"),
		jen.Return(),
	)
}

/**
func (dao *AccountDao) Create(entity *model.Account) error {
	result := d.GetDb().Create(entity)
	return result.Error
}
*/
func genCreateFun(f *jen.File, appId, tableName, idType string) {
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	entityDaoName := entityName + "Dao"

	f.Comment("创建").Line().Func().Params(
		jen.Id("d").Id("*"+entityDaoName),
	).Id("Create").Params(
		jen.Id("entity").Id("*").Qual(appId+"/model", entityName),
	).Params(
		jen.Id("error"),
	).Block(
		jen.Id("result").Op(":=").Id("d").Dot("GetDb").Call().Dot("Create").Call(
			jen.Id("entity"),
		),
		jen.Return(
			jen.Id("result").Dot("Error"),
		),
	).Line()
}
