package sqlutil


import (
	"fmt"
	"strings"
	"git.jiaxianghudong.com/go/logs"
	"errors"
	"reflect"
)

type SqlModel struct {

	err          error
	findsql 	 string
	tablesql     string
	wheresql     string
	ordersql     string
	limitsql     string
	groupsql     string
	havingsql    string
	orsql        string

	insertsql    string
	updatesql    string
	deletesql    string
}

func (sqlModel *SqlModel)Find(args...interface{}) *SqlModel{
	sqlModel.findsql = "select "
	if len(args)==0 {
		sqlModel.findsql +="* from "
		return sqlModel
	}
	for _,v:=range args {
		sqlModel.findsql+=fmt.Sprint(v)+","
	}
	sqlModel.findsql =deleteLastString(sqlModel.findsql)
	sqlModel.findsql+=" from "
	return sqlModel
}

func (sqlModel *SqlModel)Table(args... interface{})*SqlModel{
	if len(args)==0 {
		return sqlModel
	}
	sqlModel.tablesql=""
	for _,v:=range args {
		sqlModel.tablesql+=fmt.Sprint(v)+","
	}
	sqlModel.tablesql=deleteLastString(sqlModel.tablesql)+" "
	return sqlModel
}

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

//替代指定sql位置字符串
func (sqlModel *SqlModel)RSToL(sql string,location string) *SqlModel  {
	switch location {
	case "where":
		sqlModel.wheresql= sql+" "
		break
	case "select":
		sqlModel.findsql="select "+sql +" "
		break
	case "having":
		sqlModel.havingsql="having "+sql+" "
		break
	}
	return sqlModel
}

//拼接指定位置字符串
func (sqlModel *SqlModel)ISToL(sql string,location string) *SqlModel {
	switch location {
	case "where":
		sqlModel.wheresql+= sql+" "
		break
	case "select":
		sqlModel.findsql += sql +" "
		break
	case "having":
		sqlModel.havingsql += sql+" "
		break
	}
	return sqlModel
}
func (sqlModel *SqlModel)Where(where map[string]interface{},operation string) *SqlModel  {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return sqlModel
		}

		if sqlModel.wheresql != "" && operation == "and"{
			sqlModel.wheresql += " AND "
		}else if sqlModel.wheresql != "" && operation == "or" {
			sqlModel.wheresql += " OR "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					sqlModel.wheresql += fmt.Sprint(k, " IS NOT NULL")
				} else {
					sqlModel.wheresql += fmt.Sprint(k, " IS NULL")
				}
			default:
				sqlModel.wheresql += fmt.Sprint(k, "="+fmt.Sprint(v))
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				sqlModel.wheresql += fmt.Sprint(k, "="+fmt.Sprint(v))
				break
			case ">":
				sqlModel.wheresql  += fmt.Sprint(k, ">"+fmt.Sprint(v))
				break
			case ">=":
				sqlModel.wheresql += fmt.Sprint(k, ">="+fmt.Sprint(v))
				break
			case "<":
				sqlModel.wheresql += fmt.Sprint(k, "<"+fmt.Sprint(v))
				break
			case "<=":
				sqlModel.wheresql += fmt.Sprint(k, "<="+fmt.Sprint(v))
				break
			case "!=":
				sqlModel.wheresql += fmt.Sprint(k, "!="+fmt.Sprint(v))
				break
			case "<>":
				sqlModel.wheresql += fmt.Sprint(k, "!="+fmt.Sprint(v))
				break
			case "in":
				s:=fmt.Sprint(v)
				if len(s)<2 {
					sqlModel.err=errors.New("sql in parm is not array!")
					logs.Error("sql in parm is not array!")
					break
				}
				s =strings.Replace(flsub(s)," ",",",-1)
				sqlModel.wheresql += fmt.Sprint(k, " in ("+s+") ")
				break
			case "like":
				sqlModel.wheresql += fmt.Sprint(k, " like "+fmt.Sprint(v) )
			}
			break
		}
	}
	return sqlModel
}


func (sqlModel *SqlModel)Group(args...interface{}) *SqlModel {
	if len(args)==0 {
		return sqlModel
	}
	sqlModel.groupsql="group by "
	for _,v:=range args {
		sqlModel.groupsql += fmt.Sprint(v)+","
	}
	sqlModel.groupsql=deleteLastString(sqlModel.groupsql)+" "
	return sqlModel
}

func (sqlModel *SqlModel)Having(args...interface{})*SqlModel {
	if len(args)==0 {
		return sqlModel
	}
	sqlModel.havingsql="having ("
	str:=""
	for _,v := range args {
		str += str+ fmt.Sprint(v) +" and"
	}
	str = deleteLastNString(str,3)
	sqlModel.havingsql += str+" "
	return sqlModel
}

func (sqlModel *SqlModel)Order(args...interface{}) *SqlModel  {
	if len(args)==0 {
		return sqlModel
	}
	if len(args)>2 {
		sqlModel.err=errors.New("order parm to much! examp(  order('id,cid',desc))")
	}
	sqlModel.ordersql="order by "
	for _,v:=range args {
		sqlModel.ordersql+=fmt.Sprint(v)+" "
	}
	return sqlModel
}

func (sqlModel *SqlModel)Limit(args...int) *SqlModel{
	if len(args)>2 {
		sqlModel.err=errors.New("limit parm too much ")
		return sqlModel
	}
	if len(args)==0 {
		return sqlModel
	}
	sqlModel.limitsql="limit "
	for _,v:=range args {
		sqlModel.limitsql += fmt.Sprint(v)+","
	}
	sqlModel.limitsql = deleteLastString(sqlModel.limitsql)+" "
	return sqlModel
}

func (sql *SqlModel)QueryBuild() (string,error) {
	if sql.err!=nil {
		return "",sql.err
	}
	result:=sql.findsql+sql.tablesql
	if sql.wheresql!="" {
		result+="where "+sql.wheresql
	}
	result+=sql.groupsql+sql.havingsql+sql.ordersql+sql.limitsql
	return result,nil
}

//args:指定哪些字段不插入
func (sqlModel *SqlModel) Insert(obj interface{},args...interface{}) *SqlModel {
	defer func() {
		if e := recover(); e != nil{
			sqlModel.err=errors.New("insert error !")
			fmt.Println(e)
		}
	}()
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	column := " ("
	values := "values ("
	if len(args)>0 {
		argsMap:= make(map[string]interface{})
		for _,v:=range args {
			argsMap[fmt.Sprint(v)]=v
		}
		for i := 0; i < t.NumField(); i++ {
			sqlfiled:=t.Field(i).Tag.Get("sql")
			if _,ok :=argsMap[sqlfiled];ok {
				continue
			}
			column+=sqlfiled+","
			values += "'"+fmt.Sprint(v.Field(i).Interface())+"',"
		}
	}else {
		for i := 0; i < t.NumField(); i++ {
			column+=t.Field(i).Tag.Get("sql")+","
			values += "'"+fmt.Sprint(v.Field(i).Interface())+"',"
		}
	}
	column=deleteLastString(column)+") "
	values=deleteLastString(values)+") "
	sqlModel.insertsql=column+values
	return sqlModel
}

func (sqlModel *SqlModel)InsertBuild()(string,error)  {
	if sqlModel.err!=nil {
		return "",sqlModel.err
	}
	if sqlModel.tablesql=="" {
		sqlModel.err=errors.New("please add table name")
		return "",sqlModel.err
	}
	return "insert into "+sqlModel.tablesql+sqlModel.insertsql,nil
}

func (sqlModel *SqlModel)Update(obj interface{},args... interface{}) *SqlModel {
	defer func() {
		if e := recover(); e != nil{
			sqlModel.err=errors.New("update error !")
			fmt.Println(e)
		}
	}()
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	sqlModel.updatesql=" SET "
	if len(args)>0 {
		argsMap:= make(map[string]interface{})
		for _,v:=range args {
			argsMap[fmt.Sprint(v)]=v
		}
		for i := 0; i < t.NumField(); i++ {
			sqlfiled:=t.Field(i).Tag.Get("sql")
			if _,ok :=argsMap[sqlfiled];ok {
				continue
			}
			sqlModel.updatesql+=t.Field(i).Tag.Get("sql")+"='"+fmt.Sprint(v.Field(i).Interface())+"',"
		}
	}else {
		for i := 0; i < t.NumField(); i++ {
			sqlModel.updatesql+=t.Field(i).Tag.Get("sql")+"='"+fmt.Sprint(v.Field(i).Interface())+"',"
		}
	}
	sqlModel.updatesql=deleteLastString(sqlModel.updatesql)
	return sqlModel
}

func (sqlModel *SqlModel)UpdateBuild()(string,error)  {
	if sqlModel.err!=nil {
		return "",sqlModel.err
	}
	if sqlModel.tablesql=="" {
		sqlModel.err=errors.New("please add table name")
		return "",sqlModel.err
	}
	result :="update "+sqlModel.tablesql+sqlModel.updatesql
	if sqlModel.wheresql!="" {
		result += " where "+sqlModel.wheresql
	}
	return result,nil
}

func (sqlModel *SqlModel)DeleteBuild() (string,error) {
	if sqlModel.err != nil{
		return "",sqlModel.err
	}
	if sqlModel.tablesql=="" {
		sqlModel.err=errors.New("please add table name")
		return "",sqlModel.err
	}
	if sqlModel.wheresql=="" {
		sqlModel.err=errors.New("please add where condition")
		return "",sqlModel.err
	}
	return "delete from "+sqlModel.tablesql+"where "+sqlModel.wheresql,nil
}

//删除字符串最后n个
func deleteLastNString(s string,n int) string{
	return  string([]rune(s)[:len(s)-1-n])
}

//删除字符串最后一个
func deleteLastString(s string) string {
	return  string([]rune(s)[:len(s)-1])
}
//删除字符串首尾
func flsub(s string) string{
	return  string([]rune(s)[1:len(s)-1])
}

//结构体转Map
func StructToMap(stru interface{}) map[string]interface{} {
	t := reflect.TypeOf(stru).Elem()
	v := reflect.ValueOf(stru)
	resultmap := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		resultmap[t.Field(i).Tag.Get("sql")] =v.Field(i).Interface()
	}
	return resultmap
}

func Sql() *SqlModel {
	sqlmodel :=&SqlModel{}
	return sqlmodel
}