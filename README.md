# gosqlutil

example:

	andmap :=map[string]interface{}{
		"id not in": []string{"dsaf", "131", "132","133","134","135"},
		"status": "0",
		"name like":`"%麻将%"`,
		"region_id >":0,
		"id !=":130,
	}
	ormap :=map[string]interface{}{
		"name like":`"%北京%"`,
	}

	tgame,err:=sqlutil.Sql().Table("t_game").Find().Where(andmap,"and").Where(ormap,"or").QueryBuild()
	fmt.Println(err)
	fmt.Println(tgame)
//  <nil>
//select * from t_game where id not in  ('dsaf','131','132','133','134','135')  AND status= '0' AND name like '"%麻将%"' AND region_id> //'0' AND id!='130' OR name like '"%北京%"'

	querystr,err:=sqlutil.Sql().Table("t_game").Find("id,name").RSToL("id=83 or id=84 and name ='贵阳捉鸡'","where").
		          Group("id,region_id").Having("id>1").Order("id","desc").Limit(0,10).QueryBuild()
	fmt.Println(err)
	fmt.Println(querystr)
//<nil>
//select id,name from t_game where id=83 or id=84 and name ='贵阳捉鸡' group by id,region_id having (id>1) order by id desc limit 0,10 
	
	test,err:=sqlutil.Sql().Where(andmap,"and").Where(ormap,"or").Find("id").Table("t_game").QueryBuild()
	fmt.Println(err)
	fmt.Println(test)
//<nil>
//select id from t_game where id not in  ('dsaf','131','132','133','134','135')  AND status= '0' AND name like '"%麻将%"' AND region_id> //'0' AND id!='130' OR name like '"%北京%"'
	
	user := User{UserName:"ceshi",PassWord:"ceshi",CreateTime:"2018-10-18 18:00:00"}
	insert,err:=sqlutil.Sql().Table("user").Insert(user,"update_time","delete_time").InsertBuild()
	fmt.Println(err)
	fmt.Println(insert)
//<nil>
//insert into user  (id,user_name,pass_word,create_time,status) values ('0','ceshi','ceshi','2018-10-18 18:00:00','0') 
	update,err:=sqlutil.Sql().Table("user").Update(user,"update_time","delete_time").RSToL(fmt.Sprintf("id=%v",1),"where").UpdateBuild()
	fmt.Println(err)
	fmt.Println(update)
	//<nil>
	//update user  SET id='0',user_name='ceshi',pass_word='ceshi',create_time='2018-10-18 18:00:00',status='0' where id=1 

	delete,err:= sqlutil.Sql().Table("user").RSToL(fmt.Sprintf("id=%v",1),"where").DeleteBuild()
	fmt.Println(err)
	fmt.Println(delete)
//<nil>
//delete from user where id=1 

	var users []User
	user2 := User{UserName:"ceshi2",PassWord:"ceshi2",CreateTime:"2018-10-18 18:00:22"}
	users=append(users,user)
	users=append(users,user2)
	batchinsert,err:=sqlutil.Sql().Table("user").BatchInsert(users,"update_time","delete_time").InsertBuild()
	fmt.Println(err)
	fmt.Println(batchinsert)
  //<nil>
//insert into user (id,user_name,pass_word,create_time,status)  values ('0','ceshi','ceshi','2018-10-18 18:00:00','0'),				//('0','ceshi2','ceshi2','2018-10-18 18:00:22','0')
  
  
 
