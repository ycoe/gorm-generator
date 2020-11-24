package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/ycoe/gorm-generator/helper"
	"os"
)

func GenBaseDao(appId, packageName string, tableNames []string, prefix string) {
	f := jen.NewFile(packageName)
	genDaoStruct(f)
	genVarDefaultDao(f)
	genGetDao(f)
	genInit(appId, f)
	genNewDao(f, tableNames)
	genPing(f)
	genDisconnect(f)

	filename := "./" + packageName + "/dao.go"
	fmt.Println(filename)
	_ = os.MkdirAll("./"+packageName, os.ModePerm)
	if err := f.Save(filename); err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("%#v\n", f)
}

func genDisconnect(f *jen.File) {
	/**
	func (d *Dao) Disconnect() error {
		return d.client.DB().Close()
	}
	*/
	f.Func().Params(
		jen.Id("d *Dao"),
	).Id("Disconnect").Params().Id("error").Block(
		jen.Id("db, _").Op(":=").Id("d").Dot("client").Dot("DB").Call(),
		jen.Return(
			jen.Id("db").Dot("Close").Call(),
		),
	)
}

func genPing(f *jen.File) {
	/**
	func (d *Dao) Ping() error {
		return d.client.DB().Ping()
	}
	*/
	f.Func().Params(
		jen.Id("d *Dao"),
	).Id("Ping").Params().Id("error").Block(
		jen.Id("db, _").Op(":=").Id("d").Dot("client").Dot("DB").Call(),
		jen.Return(
			jen.Id("db").Dot("Ping").Call(),
		),
	).Line()
}

/**
// newDao 创建 Dao 实例
func newDao(c *conf.Config) (*Dao, error) {
	var (
		d   Dao
		err error
	)
	if d.client, err = gorm.Open(c.DB.DriverName, c.DB.URL); err != nil {
		return nil, err
	}
	d.client.SingularTable(true)       //表名采用单数形式
	d.client.DB().SetMaxOpenConns(100) //SetMaxOpenConns用于设置最大打开的连接数
	d.client.DB().SetMaxIdleConns(10)  //SetMaxIdleConns用于设置闲置的连接数
	//d.client.LogMode(true)

	if err = d.client.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.FinanceAccount{},
	).Error; err != nil {
		_ = d.client.Close()
		return nil, err
	}

	return &d, nil
}
*/
func genNewDao(f *jen.File, tableNames []string) {
	var codes []jen.Code
	for _, tableName := range tableNames {
		entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
		codes = append(codes, jen.Line().Id("&").Qual("finance/model", entityName).Block())
	}

	f.Line().Func().Id("newDao").Params(
		jen.Id("c *conf.DbConfig"),
	).Params(
		jen.Id("*Dao"),
		jen.Id("error"),
	).Block(
		jen.Var().Id("d Dao").Line().Var().Id("err error"),
		jen.If(
			jen.Id("d.client, err").Op("=").Id("gorm").Dot("Open").Call(
				jen.Line().Qual("gorm.io/driver/mysql", "Open").Call(
					jen.Id("c.DB.URL"),
				),
				jen.Id("&gorm").Dot("Config").Values(
					jen.Dict{
						jen.Line().Id("NamingStrategy"): 	jen.Qual("gorm.io/gorm/schema", "NamingStrategy").Values(
							jen.Dict{
								jen.Id("TablePrefix"): 	jen.Id("c.DB.TablePrefix"),
								jen.Id("SingularTable"): 	jen.True(),
							},
						),
					},
				),
			),
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Return(
				jen.Nil(),
				jen.Id("err"),
			),
		).Line(),

		jen.Id("sqlDB, _").Op(":=").Id("d").Dot("client").Dot("DB").Call().Line(),

		jen.Id("sqlDB").Dot("SetMaxIdleConns").Call(
			jen.Id("c").Dot("DB").Dot("MaxIdleConns"),
		).Comment("SetMaxIdleConns 设置空闲连接池中连接的最大数量"),
		jen.Id("sqlDB").Dot("SetMaxOpenConns").Call(
			jen.Id("c").Dot("DB").Dot("MaxOpenConns"),
		).Comment("SetMaxOpenConns 设置打开数据库连接的最大数量"),
		jen.Id("sqlDB").Dot("SetConnMaxLifetime").Call(
			jen.Qual("time", "Hour"),
		).Comment("SetConnMaxLifetime 设置了连接可复用的最大时间"),

		jen.Id("d").Dot("client").Dot("Set").Call(
			jen.Lit("gorm:table_options"),
			jen.Lit("ENGINE=InnoDB"),
		),
		//jen.If(
		//	jen.Id("err").Op("=").Id("d").Dot("client").Dot("Set").Call(
		//		jen.Lit("gorm:table_options"),
		//		jen.Lit("ENGINE=InnoDB"),
		//	).Dot("AutoMigrate").Call(codes...),
		//	jen.Id("err").Op("!=").Nil(),
		//).Block(
		//	jen.Return(
		//		jen.Nil(),
		//		jen.Id("err"),
		//	),
		//),
		jen.Return(
			jen.Id("&d"),
			jen.Nil(),
		),
	).Line()
}

/**
func Init(c *conf.Config) (err error) {
	defaultDao, err = newDao(c)
	return
}
*/
func genInit(appId string, f *jen.File) {
	f.Line().Func().Id("Init").Params(
		jen.Id("c *").Qual(appId+"/conf", "DbConfig"),
	).Params(
		jen.Id("err").Id("error"),
	).Block(
		jen.Id("defaultDao, err").Op("=").Id("newDao").Call(
			jen.Id("c"),
		).Line().Return(),
	)
}

/**
func GetDao() *Dao {
	return defaultDao
}
*/
func genGetDao(f *jen.File) {
	f.Line().Func().Id("GetDao").Params().Params(
		jen.Id("*Dao"),
	).Block(
		jen.Return(
			jen.Id("defaultDao"),
		),
	)
}

func genVarDefaultDao(f *jen.File) *jen.Statement {
	return f.Var().Id("defaultDao").Id("*Dao")
}

func genDaoStruct(f *jen.File) *jen.Statement {
	return f.Type().Id("Dao").Struct(
		jen.Id("client").Id("*").Qual("gorm.io/gorm", "DB"),
	)
}
