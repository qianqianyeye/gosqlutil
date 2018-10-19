package sqlutil


import (
	"fmt"
	"strings"
	"errors"
	"reflect"
	"git.jiaxianghudong.com/go/logs"
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
	joinsql      string
	betweensql   string

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
	case "update":
		sqlModel.updatesql=" SET "+"sql"+" "
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
	case "update":
		sqlModel.updatesql += sql+" "
		break
	case "having":
		sqlModel.havingsql += sql+" "
		break
	}
	return sqlModel
}
func (sqlModel *SqlModel) And(args...interface{}) *SqlModel{
	if len(args)==0 {
		return sqlModel
	}
	for _,v:=range args {
		sqlModel.wheresql+= " and "+fmt.Sprint(v)
	}
	return sqlModel
}

func (sqlModel *SqlModel) Or(args... interface{}) *SqlModel{
	if len(args)==0 {
		return sqlModel
	}
	for _,v:=range args {
		sqlModel.wheresql+= " or "+fmt.Sprint(v)
	}
	return sqlModel
}

func (sqlModel *SqlModel)Where(where map[string]interface{},operation string) *SqlModel  {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			if len(ks)>3 {
				return sqlModel
			}
			if ks[1]=="not"&&ks[2]=="in" {

			}else if ks[1]=="not"&&ks[2]=="between" {

			}else {
				return sqlModel
			}
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
				sqlModel.wheresql += fmt.Sprint(k, "= '"+fmt.Sprint(v)+"'")
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				sqlModel.wheresql += fmt.Sprint(k, "= '"+fmt.Sprint(v)+"'")
				break
			case ">":
				sqlModel.wheresql  += fmt.Sprint(k, "> '"+fmt.Sprint(v)+"'")
				break
			case ">=":
				sqlModel.wheresql += fmt.Sprint(k, ">= '"+fmt.Sprint(v)+"'")
				break
			case "<":
				sqlModel.wheresql += fmt.Sprint(k, "< '"+fmt.Sprint(v)+"'")
				break
			case "<=":
				sqlModel.wheresql += fmt.Sprint(k, "<= '"+fmt.Sprint(v)+"'")
				break
			case "!=":
				sqlModel.wheresql += fmt.Sprint(k, "!='"+fmt.Sprint(v)+"'")
				break
			case "<>":
				sqlModel.wheresql += fmt.Sprint(k, "!= '"+fmt.Sprint(v)+"'")
				break
			case "in":
				if len(fmt.Sprint(v))<2{
					sqlModel.err=errors.New("in parms must be arr")
					break
				}
				s:=getinstr(v)
				sqlModel.wheresql += fmt.Sprint(k, " in ("+s+") ")
				break
			case "between":
				if len(fmt.Sprint(v))<2{
					sqlModel.err=errors.New("between parms must be arr")
					break
				}
				s:=getbetweenstr(v)
				sqlModel.wheresql += fmt.Sprint(k," between ("+s+") ")
				break
			case "like":
				sqlModel.wheresql += fmt.Sprint(k, " like '"+fmt.Sprint(v)+"'" )
			}
			break
			case 3:
				k = ks[1]+" "+ks[2]
				switch k {
				case "not in":
					if len(fmt.Sprint(v))<2{
						sqlModel.err=errors.New("parms must be arr")
						break
					}
					s:=getinstr(v)
					sqlModel.wheresql += fmt.Sprint(ks[0]+" "+k, "  ("+s+") ")
					break
				case "not between":
					if len(fmt.Sprint(v))<2{
						sqlModel.err=errors.New("parms must be arr")
						break
					}
					s:=getbetweenstr(v)
					sqlModel.wheresql += fmt.Sprint(ks[0]+" "+k, "  ("+s+") ")
					break
				}
				break
		}
	}
	return sqlModel
}

func getinstr(v interface{})  string{
	s:=""
	switch v.(type) {
	case []string:
		for _,value:=range v.([]string){
			s+="'"+value+"',"
		}
		s=deleteLastString(s)
		break
	default:
		s=fmt.Sprint(v)
		s =strings.Replace(flsub(s)," ",",",-1)
		break
	}
	return s
}

func getbetweenstr(v interface{})string  {
	s:=""
	switch v.(type) {
	case []string:
		for _,value:=range v.([]string){
			s+=" '"+value+"' and"
		}
		s=deleteLastNString(s,3)
		break
	default:
		s=fmt.Sprint(v)
		s =strings.Replace(flsub(s)," "," and ",-1)
		break
	}
	return s
}

func (sqlModel *SqlModel)JoinSql()  {

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
	sqlModel.havingsql += str+") "
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

//构建查询语句
func (sql *SqlModel)QueryBuild() (string,error) {
	if sql.err!=nil {
		return "",sql.err
	}
	result:=sql.findsql+sql.tablesql
	if sql.wheresql!="" {
		result+="where "+sql.wheresql
		//if sql.betweensql!="" {
		//	result+=sql.betweensql
		//}
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

//生成批量插入语句 参数传如切片类型，args指定字段不插入
func (sqlModel *SqlModel)BatchInsert(obj interface{},args...interface{}) *SqlModel{
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj).Elem()

	if v.Kind() != reflect.Slice {
		sqlModel.err=errors.New("toslice arr not slice")
		logs.Error("toslice arr not slice")
		return sqlModel
	}

	argsMap:= make(map[string]interface{})
	for _,v:=range args {
		argsMap[fmt.Sprint(v)]=v
	}
	columnsql:="("
	for i := 0; i < t.NumField(); i++ {
		sqlfiled:=t.Field(i).Tag.Get("sql")
		if _,ok :=argsMap[sqlfiled];ok {
			continue
		}
		columnsql+=sqlfiled+","
	}
	columnsql=deleteLastString(columnsql)
	columnsql +=") "
	valuesql :=" values "
	l := v.Len()
	for i := 0; i < l; i++ {
		valuesql+="("
		values :=reflect.ValueOf(v.Index(i).Interface())
		for i:=0;i<values.NumField();i++ {
			sqlfiled:=t.Field(i).Tag.Get("sql")
			if _,ok :=argsMap[sqlfiled];ok {
				continue
			}
			valuesql+="'"+fmt.Sprint(values.Field(i).Interface())+"',"
		}
		valuesql=deleteLastString(valuesql)
		valuesql+="),"
	}
	valuesql=deleteLastString(valuesql)
	sqlModel.insertsql=columnsql+valuesql
	return sqlModel
}

//构建插入语句
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

//args指定字段不更新
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

//构建更新语句
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

//构建删除语句
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