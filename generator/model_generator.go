package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/jinzhu/inflection"
	"github.com/ycoe/gorm-generator/helper"
	"os"
	"strings"
)

func GenerateModel(tableName string, columns []map[string]string, dir string) string {
	index := strings.LastIndex(dir, "/")
	daoPackage := dir[index+1 : len(dir)]
	var codes []jen.Code
	idType := ""
	for i, col := range columns {
		t := col["Type"]
		column := col["Field"]
		var st *jen.Statement
		st = jen.Id(helper.SnakeCase2CamelCase(column, true))
		columnType := getCol(st, t)
		if idType == "" && i == 0 {
			idType = columnType
		}
		st.Tag(map[string]string{"json": column})
		st.Comment(col["Comment"])
		codes = append(codes, st)
	}
	f := jen.NewFile(daoPackage)
	f.HeaderComment("Code generated by https://github.com/ycoe/gorm-generator. DO NOT EDIT.")
	f.ImportAlias("time", "time")
	entityName := helper.SnakeCase2CamelCase(inflection.Singular(tableName), true)
	f.Type().Id(entityName).Struct(codes...)

	genEncoding(f, entityName)

	_ = os.MkdirAll(dir, os.ModePerm)
	fileName := dir + "/" + inflection.Singular(tableName) + ".go"
	fmt.Println(fileName)
	if err := f.Save(fileName); err != nil {
		fmt.Println(err.Error())
	}
	return idType
}

/**
func (s QyWxStaff) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s QyWxStaff) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, &s)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
 */
func genEncoding(f *jen.File, name string) {
	e := strings.ToLower(name[0:1])
	f.Func().Params(
		jen.Id(e).Id(name),
	).Id("MarshalBinary").Params().Params(
		jen.Id("[]byte"),
		jen.Id("error"),
	).Block(
		jen.Return().Qual("encoding/json", "Marshal").Call(
			jen.Id(e),
		),
	).Line()

	f.Func().Params(
		jen.Id(e).Id(name),
	).Id("UnmarshalBinary").Params(
		jen.Id("data").Id("[]byte"),
	).Params(
		jen.Id("error"),
	).Block(
		jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
			jen.Id("data"),
			jen.Id("&" + e),
		),

		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Qual("gitee.com/inngke/go-base-service/common/loggers", "Error").Call(
				jen.Id("err").Dot("Error").Call(),
			),
			jen.Return(
				jen.Id("err"),
			),
		),

		jen.Return(
			jen.Nil(),
		),
	)
}

func getCol(st *jen.Statement, t string) string {
	prefix := strings.Split(t, "(")[0]
	switch prefix {
	case "int", "tinyint", "smallint", "mediumint":
		st.Int32()
		return "int32"
	case "bigint":
		st.Int64()
		return "int64"
	case "float", "decimal":
		st.Float32()
		return "float32"
	case "date", "time", "timestamp", "year", "datetime":
		st.Id("*").Qual("time", "Time")
		return "time.Time"
	default:
		st.String()
		return "string"
	}
}
