package main

import (
	"Hello/src/sqlstring"
	"fmt"
)
type User struct {
	ID int64 `sql:"id" json:"id"`
	UserName string `sql:"user_name" json:"user_name"`
	PassWord string `sql:"pass_word" json:"pass_word"`
	CreateTime string `sql:"create_time" json:"create_time"`
	UpdateTime string `sql:"update_time" json:"update_time"`
	DeleteTime string `sql:"delete_time" json:"delete_time"`
	Status int64 `sql:"status" json:"status"`
}


func main()  {
	andmap :=map[string]interface{}{
		"status": "0",
		"name like":`"%麻将%"`,
		"id in": []int64{130, 131, 132,133,134,135},
		"region_id >":0,
		"id !=":130,
	}
	ormap :=map[string]interface{}{
		"name like":`"%北京%"`,
	}

	tgame,_:=sqlutil.Sql().Table("t_game").Find().Where(andmap,"and").Where(ormap,"or").QueryBuild()
	fmt.Println(tgame)

	querystr,err:=sqlutil.Sql().Table("t_game").Find().RSToL("id=83 or id=84 and name ='贵阳捉鸡'","where").
		          Group("id,region_id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild()
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(querystr)

	test,_:=sqlutil.Sql().Where(andmap,"and").Where(ormap,"or").Find("id").Table("t_game").QueryBuild()
	fmt.Println(test)

	user := User{UserName:"ceshi",PassWord:"ceshi",CreateTime:"2018-10-18 18:00:00"}

	insert,_:=sqlutil.Sql().Table("user").Insert(user,"update_time","delete_time").InsertBuild()
	fmt.Println(insert)

	update,_:=sqlutil.Sql().Table("user").Update(user,"update_time","delete_time").RSToL(fmt.Sprintf("id=%v",1),"where").UpdateBuild()
	fmt.Println(update)

	delete,_:= sqlutil.Sql().Table("user").RSToL(fmt.Sprintf("id=%v",1),"where").DeleteBuild()
	fmt.Println(delete)

}
