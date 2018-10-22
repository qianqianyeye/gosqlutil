# gosqlutil

example:

andmap :=map[string]interface{}{
		"id not in": []string{"dsaf", "131", "132","133","134","135"},
		"status": "0",
		"user_name like":"%测试2%",
		"id !=":130,
		"create_time between":[]string{"2018-10-17 18:00:00","2018-10-20 19:00:00"},
	}
	ormap :=map[string]interface{}{
		"user_name like":"%测试2%",
	}

	//查询

	var user1 []User
	query,_:=sqlutil.Sql().Table("user").Find().QueryBuild(&user1)
	fmt.Println("1:"+query)
	fmt.Println("user1:",user1)

	var user2 []User
	tgame,err:=sqlutil.Sql().Table("user").Find().Where(andmap,"and").Where(ormap,"or").QueryBuild(&user2)
	fmt.Println(err)
	fmt.Println("2:"+tgame)
	fmt.Println("user2:",user2)

	var user3 []User
	//RSToL 替代某个位置的sql语句
	querystr,err:=sqlutil.Sql().Table("user").Find("id,user_name").RSToL("id=3 or id=4 and user_name ='ceshi2'","where").
		          Group("id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild(&user3)
	fmt.Println(err)
	fmt.Println("3:"+querystr)
	fmt.Println("user3:",user3)

	var user4 []User
	//ISToL 拼接某个位置的sql语句
	querystr2,err:=sqlutil.Sql().Table("user").Find("id,user_name").RSToL("user_name ='ceshi'","where").ISToL("and id=8 ","where").
		Group("id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild(&user4)
	fmt.Println(err)
	fmt.Println("4:"+querystr2)
	fmt.Println("user4:",user4)

	var user5 []User
	querystr3,err:=sqlutil.Sql().Table("user").Find("id,user_name").Where(andmap,"and").ISToL("and id=8 or id=9 ","where").
		Group("id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild(&user5)
	fmt.Println(err)
	fmt.Println("5:"+querystr3)
	fmt.Println("user5:",user5)


	var user6 []User
	test,err:=sqlutil.Sql().Where(andmap,"and").Where(ormap,"or").Find("id").Table("user").QueryBuild(&user6)
	fmt.Println(err)
	fmt.Println("6:"+test)
	fmt.Println("user6:",user6)

	//插入
	user := User{UserName:"ceshi",PassWord:"ceshi",CreateTime:"2018-10-18 18:00:00"}
	insert,lastid,err:=sqlutil.Sql().Table("user").Insert(user,"update_time","delete_time").InsertBuild()
	fmt.Println(err)
	fmt.Println("7:"+insert)
	fmt.Println("lastid:",lastid)

	//修改
	update,rowseffect,err:=sqlutil.Sql().Table("user").Update(user,"id","update_time","delete_time").RSToL(fmt.Sprintf("id=%v",1),"where").UpdateBuild()
	fmt.Println(err)
	fmt.Println("8:"+update)
	fmt.Println("rowseffect:",rowseffect)

	update2,rowseffect2,err :=sqlutil.Sql().Table("user").RSToL(fmt.Sprint("user_name='dfas'"),"update").RSToL(fmt.Sprint("id=5"),"where").UpdateBuild()
	fmt.Println(err)
	fmt.Println("9:"+update2)
	fmt.Println("rowseffect:",rowseffect2)

	//删除
	delete,drowseffect,err:= sqlutil.Sql().Table("user").RSToL(fmt.Sprintf("id=%v",2),"where").DeleteBuild()
	fmt.Println(err)
	fmt.Println("9:"+delete)
	fmt.Println("delete rowseffect:",drowseffect)

	var users []User
	user0 := User{UserName:"ceshi2",PassWord:"ceshi2",CreateTime:"2018-10-18 18:00:22"}
	for i:=0;i<10;i++ {
		users=append(users,user0)
	}
	//批量插入
	batchinsert,lastid,err:=sqlutil.Sql().Table("user").BatchInsert(users,"update_time","delete_time","id").InsertBuild()
	fmt.Println(err)
	fmt.Println("10:"+batchinsert)
	fmt.Println("lastid:",lastid)

	//事物
	tx:=sqlutil.Sql().TxStart()
	txuser1 := User{UserName:"txceshi",PassWord:"txceshi",CreateTime:"2018-10-18 18:00:00"}
	txuser2 := User{UserName:"txceshi2",PassWord:"txceshi2",CreateTime:"2018-10-18 18:00:00"}
	txuser3 := User{UserName:"txceshi3",PassWord:"txceshi3",CreateTime:"2018-10-18 18:00:00"}
	tx.Table("user").Insert(txuser1,"update_time","delete_time").TxInsertBuild()
	tx.Table("user").Insert(txuser2,"update_time","delete_time").TxInsertBuild()
	tx.Table("user").Update(txuser3,"id","update_time","delete_time").RSToL(fmt.Sprint("id=1"),"where").TxUpdateBuild()
	s,err:=tx.TxBuild()
	fmt.Println(err)
	for _,v:=range s{
		fmt.Println(v)
	}
