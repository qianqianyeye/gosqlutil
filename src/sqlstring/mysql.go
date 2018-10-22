package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var (
	driver string
	mOpen  int
	mIdle  int
)

func Init(mysqlDirver string, maxOpen, maxIdle int) error {
	driver = mysqlDirver
	maxOpen = mOpen
	maxIdle = mIdle
	db, _ = sql.Open("mysql", driver)
	//设置最大打开的连接数，默认值为0表示不限制,若不限制并发太高有可能导致连接mysql出现too many connections的错误
	db.SetMaxOpenConns(maxOpen)
	//设置闲置的连接数，当开启的一个连接使用完成后可以放在池里等候下一次使用
	db.SetMaxIdleConns(maxIdle)
	return db.Ping()

}

/*
 *	获取mysql_adm操作对象
 */
func GetClient() (*sql.DB, error) {
	err := db.Ping()
	//断线重连5次
	if err != nil {
		for i := 0; i < 5; i++ {
			err = Init(driver, mOpen, mIdle)
			if err == nil {
				return db, nil
			}
			time.Sleep(1 * time.Second)
		}
		return nil, err
	}
	return db, nil
}

// 执行语句
func Exec(sql string) (int64, error) {
	o, err := GetClient()
	if err != nil {
		return 0, err
	}
	ret, err := o.Exec(sql)
	if nil != err {
		return 0, err
	}

	return ret.RowsAffected()
}

// 执行事务
func ExecTrans(sqls []string) error {
	o, err := GetClient()
	if err != nil {
		return err
	}
	tx, err := o.Begin()
	if nil != err {
		return err
	}

	for i := 0; i < len(sqls); i++ {
		row, err := tx.Query(sqls[i])
		if nil != err {
			tx.Rollback()
			return err
		}

		row.Close()
	}

	return tx.Commit()
}

// 插入
func Insert(sql string) (int64, error) {
	o, err := GetClient()
	if err != nil {
		return 0, err
	}
	result, err := o.Exec(sql)
	if nil != err {
		return 0, err
	}
	rowsaffect,_ :=result.RowsAffected()
	lastid,err:=result.LastInsertId()
	lastid=lastid+rowsaffect-1
	return lastid,err
}

// 查询
func Query(sql string, val interface{}) error {
	var tagMap map[string]int
	var tp, tps reflect.Type
	var n, i int
	var err error
	var ret reflect.Value
	// 检测val参数是否为我们所想要的参数
	tp = reflect.TypeOf(val)
	if reflect.Ptr != tp.Kind() {
		return errors.New("is not pointer")
	}

	if reflect.Slice != tp.Elem().Kind() {
		return errors.New("is not slice pointer")
	}

	tp = tp.Elem()
	tps = tp.Elem()
	if reflect.Struct != tps.Kind() {
		return errors.New("is not struct slice pointer")
	}

	tagMap = make(map[string]int)
	n = tps.NumField()
	for i = 0; i < n; i++ {
		tag := tps.Field(i).Tag.Get("sql")
		if len(tag) > 0 {
			tagMap[tag] = i + 1
		}
	}

	// 执行查询
	ret, err = queryAndReflect(sql, tagMap, tp)
	if nil != err {
		return err
	}

	// 返回结果
	reflect.ValueOf(val).Elem().Set(ret)

	return nil
}

// 查询并构建返回
func queryAndReflect(sql string,
	tagMap map[string]int,
	tpSlice reflect.Type) (reflect.Value, error) {

	var ret reflect.Value

	o, err := GetClient()
	if err != nil {
		return ret, err
	}
	// 执行sql语句
	rows, err := o.Query(sql)
	if nil != err {
		return reflect.Value{}, err
	}

	defer rows.Close()
	// 开始枚举结果
	cols, err := rows.Columns()
	if nil != err {
		return reflect.Value{}, err
	}

	ret = reflect.MakeSlice(tpSlice, 0, 50)
	// 构建接收队列
	scan := make([]interface{}, len(cols))
	row := make([]interface{}, len(cols))
	for r := range row {
		scan[r] = &row[r]
	}

	for rows.Next() {
		feild := reflect.New(tpSlice.Elem()).Elem()
		// 取得结果

		err = rows.Scan(scan...)
		// 开始遍历结果
		for i := 0; i < len(cols); i++ {
			n := tagMap[cols[i]] - 1
			if n < 0 {
				continue
			}
			switch feild.Type().Field(n).Type.Kind() {
			case reflect.Bool:
				if nil != row[i] {
					feild.Field(n).SetBool("false" != string(row[i].([]byte)))
				} else {
					feild.Field(n).SetBool(false)
				}
			case reflect.String:
				if nil != row[i] {
					feild.Field(n).SetString(string(row[i].([]byte)))
				} else {
					feild.Field(n).SetString("")
				}
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				if nil != row[i] {
					v, e := strconv.ParseFloat(string(row[i].([]byte)), 0)
					if nil == e {
						feild.Field(n).SetFloat(v)
					}
				} else {
					feild.Field(n).SetFloat(0)
				}
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				fallthrough
			case reflect.Int:
				if nil != row[i] {
					byRow, ok := row[i].([]byte)
					if ok {
						v, e := strconv.ParseInt(string(byRow), 10, 64)
						if nil == e {
							feild.Field(n).SetInt(v)
						}
					} else {
						v, e := strconv.ParseInt(fmt.Sprint(row[i]), 10, 64)
						if nil == e {
							feild.Field(n).SetInt(v)
						}
					}
				} else {
					feild.Field(n).SetInt(0)
				}
			}
		}

		ret = reflect.Append(ret, feild)
	}

	return ret, nil
}

// 获取单一值
func GetString(sql string, args ...interface{}) (string, error) {

	o, err := GetClient()
	if err != nil {
		return "", err
	}

	row, err := o.Query(format(sql, args...))
	if nil != err {
		return "", err
	}

	defer row.Close()
	if row.Next() {
		var col interface{}
		err = row.Scan(&col)
		if nil != err {
			return "", err
		}

		return string(col.([]byte)), nil
	}

	return "", nil
}

// 获取int
func GetInt64(sql string, args ...interface{}) (int64, error) {
	str, err := GetString(sql, args...)
	if nil != err {
		return 0, err
	}

	return strconv.ParseInt(str, 10, 64)
}

func format(fsql string, args ...interface{}) string {
	sql := fsql
	i := 0
	for ; i < len(args); i++ {
		key := fmt.Sprintf("{%d}", i)
		val := strings.Replace(fmt.Sprint(args[i]), "'", "''", -1)
		sql = strings.Replace(sql, key, val, -1)
	}

	return sql
}
