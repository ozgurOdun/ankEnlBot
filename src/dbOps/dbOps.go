package dbOps

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"time"
)

var o orm.Ormer
var w io.Writer

type Task struct {
	Uid          int       `orm:"auto"` // if the primary key is not id, you need to add tag `PK` for your customized primary key.
	Task         string    `orm:"size(160)"`
	CreationTime time.Time `orm:"auto_now_add;type(datetime)"`
	Ownerid      int       `orm:"default(1)"`
	Done         bool      `orm:"default(false)"`
}

type User struct {
	Uid          int       `orm:"pk"`
	UserName     string    `orm:"size(100)"`
	FirstName    string    `orm:"size(100);NULL"`
	LastName     string    `orm:"size(100);NULL"`
	CreationTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func Hehe() {
	fmt.Printf("Running\n")
}

func NewDb() {
	o = orm.NewOrm()
	o.Using("default")

}

func init() {
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "database/new.db")
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Task))
}

func AddNewUser(id int, uname string, fname string, lname string) bool {
	user := new(User)
	user.Uid = id
	user.UserName = uname
	user.FirstName = fname
	user.LastName = lname
	_, err := o.Insert(user)
	if err != nil {
		fmt.Printf("Error occured when inserting new user", err)
		return false
	}

	return true
}

func CheckUser(uname string) bool {
	user := new(User)
	err := o.QueryTable("user").Filter("username", uname).One(user)

	if err == orm.ErrMultiRows {
		fmt.Println("returned multirows not one")
		return true
	} else if err == orm.ErrNoRows {
		fmt.Println("no row found")
		return false
	} else {
		fmt.Println(user.Uid, user.UserName)
		return true
	}

}

func GetUserByName(uname string) *User {
	user := new(User)
	err := o.QueryTable("user").Filter("username", uname).One(user)

	if err == orm.ErrMultiRows {
		fmt.Println("returned multirows not one")
		return nil
	} else if err == orm.ErrNoRows {
		fmt.Println("no row found")
		return nil
	} else {
		fmt.Println(user.Uid, user.UserName)
		return user
	}
}

func GetUserById(uid int) *User {
	user := new(User)
	err := o.QueryTable("user").Filter("uid", uid).One(user)

	if err == orm.ErrMultiRows {
		fmt.Println("returned multirows not one")
		return nil
	} else if err == orm.ErrNoRows {
		fmt.Println("no row found")
		return nil
	} else {
		fmt.Println(user.Uid, user.UserName)
		return user
	}
}

func AddNewTask(oId int, text string) int64 {
	task := new(Task)
	task.Ownerid = oId
	task.Task = text
	task.Done = false
	id, err := o.Insert(task)
	if err != nil {
		fmt.Printf("Error occured when inserting new task\n", err)
		return -1
	} else {
		fmt.Printf("New task added\n")
	}

	return id
}

func ListUndoneTasks() ([]*Task, int64) {
	var tasks []*Task
	num, _ := o.QueryTable("task").Filter("done", false).All(&tasks)
	fmt.Printf("Returned Rows Num: %d\n", num)

	return tasks, num
}

func ListUndoneTasksByUser(uname string) ([]*Task, int64) {
	var tasks []*Task
	user := GetUserByName(uname)
	num, _ := o.QueryTable("task").Filter("done", false).Filter("ownerid", user.Uid).All(&tasks)
	fmt.Printf("Returned Rows Num: %d\n", num)

	return tasks, num
}

func MarkTaskAsDone(id int) bool {
	task := new(Task)
	task.Uid = id
	if o.Read(task) == nil {
		task.Done = true
		if num, err := o.Update(task); err == nil {
			fmt.Println(num)
		} else {
			return false
		}
	} else {
		return false
	}

	return true
}

func MarkTaskAsUndone(id int) bool {
	task := new(Task)
	task.Uid = id
	if o.Read(task) == nil {
		task.Done = false
		if num, err := o.Update(task); err == nil {
			fmt.Println(num)
		} else {
			return false
		}
	} else {
		return false
	}

	return true
}

func DeleteTask(id int) bool {
	task := new(Task)
	task.Uid = id
	if num, err := o.Delete(task); err == nil {
		fmt.Println(num)
	} else {
		return false
	}

	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
