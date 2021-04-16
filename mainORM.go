package main

import (
	"database/sql"
	"fmt"
	"gee"
	"geeorm"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

const (
	userName = "pbroot"
	password = "1qaz!QAZ"
	ip       = "bms-rlzj-dev-pub.rwlb.rds.aliyuncs.com"
	port     = "3306"
	dbName   = "mall"
)

type User struct {
	Name string
}

func mainORM() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	engine, _ := geeorm.NewEngine("mysql", path)
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	//u1 := &User{Name: "潘金莲"}
	//u2 := &User{Name: "李瓶儿"}
	//affected, err := s.Insert(u1, u2)
	//if err != nil || affected < 1 {
	//	log.Println("failed to insert record")
	//}

	//affected,_ := s.Where("Name = ?","潘金莲").Update("name","爱丽丝")
	//if affected < 1 {
	//	log.Println("failed to update name")
	//}

	affected, _ := s.Where("Name = ?", "爱丽丝").Delete()
	count, _ := s.Count()
	if affected < 1 || count < 1 {
		log.Println("failed to delete or count")
	}

	var users []User
	if err1 := s.Limit(10).Find(&users); err1 != nil || len(users) < 1 {
		log.Printf("failed to query all")
	}
	fmt.Println(users)
}

func sql_test() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	engine, _ := geeorm.NewEngine("mysql", path)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`name`) VALUES (?),(?)", "张三", "李四").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success,%d affected\n", count)

}

func db_test() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	db, _ := sql.Open("mysql", path)
	defer func() { _ = db.Close() }()
	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User(`Name`) VALUE (?),(?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}

func test2() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})

	//index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}

func test1() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Run(":9999")
}

func onlyForV1() gee.HandlerFunc {
	return func(c *gee.Context) {
		//Start timer
		t := time.Now()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v1", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func test() {
	r := gee.New()
	r.Use(gee.Logger())

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v1 := r.Group("/v1")
	v1.Use(onlyForV1())
	{

		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	r.Run(":9999")
}
