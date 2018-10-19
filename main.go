package main

import (
	"fmt"
	"sqlutil/src/sqlstring"
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
		"id not in": []string{"dsaf", "131", "132","133","134","135"},
		"status": "0",
		"name like":`"%麻将%"`,
		"region_id >":0,
		"id !=":130,
		"create_time between":[]string{"2018-10-18 18:00:00","2018-10-18 19:00:00"},
	}
	ormap :=map[string]interface{}{
		"name like":`"%北京%"`,
	}

	query,_:=sqlutil.Sql().Table("user").Find().QueryBuild()
	fmt.Println("1:"+query)

	tgame,err:=sqlutil.Sql().Table("t_game").Find().Where(andmap,"and").Where(ormap,"or").QueryBuild()
	fmt.Println(err)
	fmt.Println("2:"+tgame)

	//RSToL 替代某个位置的sql语句
	querystr,err:=sqlutil.Sql().Table("t_game").Find("id,name").RSToL("id=83 or id=84 and name ='贵阳捉鸡'","where").
		          Group("id,region_id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild()
	fmt.Println(err)
	fmt.Println("3:"+querystr)

	//ISToL 拼接某个位置的sql语句
	querystr2,err:=sqlutil.Sql().Table("t_game").Find("id,name").RSToL("name ='贵阳捉鸡'","where").ISToL("and id=83 or id=84 ","where").
		Group("id,region_id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild()
	fmt.Println(err)
	fmt.Println("4:"+querystr2)

	querystr3,err:=sqlutil.Sql().Table("t_game").Find("id,name").Where(andmap,"and").ISToL("and id=83 or id=84 ","where").
		Group("id,region_id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild()
	fmt.Println(err)
	fmt.Println("5:"+querystr3)


	test,err:=sqlutil.Sql().Where(andmap,"and").Where(ormap,"or").Find("id").Table("t_game").QueryBuild()
	fmt.Println(err)
	fmt.Println("6:"+test)

	user := User{UserName:"ceshi",PassWord:"ceshi",CreateTime:"2018-10-18 18:00:00"}
	insert,err:=sqlutil.Sql().Table("user").Insert(user,"update_time","delete_time").InsertBuild()
	fmt.Println(err)
	fmt.Println("7:"+insert)

	update,err:=sqlutil.Sql().Table("user").Update(user,"update_time","delete_time").RSToL(fmt.Sprintf("id=%v",1),"where").UpdateBuild()
	fmt.Println(err)
	fmt.Println("8:"+update)

	delete,err:= sqlutil.Sql().Table("user").RSToL(fmt.Sprintf("id=%v",1),"where").DeleteBuild()
	fmt.Println(err)
	fmt.Println("9:"+delete)

	var users []User
	user2 := User{UserName:"ceshi2",PassWord:"ceshi2",CreateTime:"2018-10-18 18:00:22"}
	users=append(users,user)
	for i:=0;i<100;i++ {
		users=append(users,user2)
	}

	batchinsert,err:=sqlutil.Sql().Table("user").BatchInsert(users,"update_time","delete_time","id").InsertBuild()
	fmt.Println(err)
	fmt.Println("10:"+batchinsert)

}
