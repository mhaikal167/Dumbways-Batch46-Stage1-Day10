package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web-week-3/connection"
	"personal-web-week-3/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id       int
	Title    string
	Start    string
	End      string
	Desc     string
	Techno   []string
	Duration int
	PostDate string
	Images string
	Author string
}
type User struct {
	IdUser   int
	Username string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin  bool
	Username string
}

var userData = SessionData{}
var node string
var react string
var next string
var typesc string

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
func main() {
	connection.DatabaseConnect()
	// new echo intance
	e := echo.New()
	// static public directory
	e.Static("/public", "public")
	e.Static("/upload", "upload")
	// init session middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))
	//Routing
	e.GET("/", home)
	e.GET("/login", loginPage)
	e.GET("/register", registerPage)
	e.GET("/contact-me", contactme)
	e.GET("/add-project", formproject)
	e.GET("/detail-project/:id", detailproject)
	e.GET("/delete-project/:id", deleteProject)
	e.GET("/edit-project/:id", editProject)
	e.GET("/logout", logout)

	e.POST("/add", middleware.UploadFile(addProject))
	e.POST("/edit/:id", middleware.UploadFile(editFunc))
	e.POST("/register-func", register)
	e.POST("/login-func", login)
	e.Logger.Fatal(e.Start("localhost:8000"))

}

func loginPage(c echo.Context) error {
	sess, _ := session.Get("session", c)
	fmt.Println(sess)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"FlashIsLogin": sess.Values["isLogin"],
		"FlashUsername" : sess.Values["username"],
	}
	delete(sess.Values, "status")
	delete(sess.Values, "message")

	sess.Save(c.Request(),c.Response())

	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal()
	}

	username := c.FormValue("user")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * from tb_user where username=$1", username).Scan(&user.Username, &user.Email, &user.Password, &user.IdUser)

	if err != nil {
		return redirectMessage(c, "Username not found,please register", false, "/login")
	}
	fmt.Println(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println(err)
	if err != nil {
		return redirectMessage(c, "Wrong Password", false, "/login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800
	sess.Values["message"] = "Login Success"
	sess.Values["statusAlert"] = true
	sess.Values["username"] = user.Username
	sess.Values["idUser"] = user.IdUser
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func registerPage(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"FlashIsLogin": sess.Values["isLogin"],
		"FlashUsername": sess.Values["username"],
	}

	delete(sess.Values, "status")
	delete(sess.Values, "message")
	sess.Save(c.Request(), c.Response())
	tmpl, err := template.ParseFiles("views/register.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), flash)
}
func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal()
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "insert into tb_user(username,email,password)values($1,$2,$3)", name, email, passwordHash)
	if err != nil {
		fmt.Println(err)
		redirectMessage(c, "Register Failed maybe Username already used, please try again", false, "/register")
	}
	return redirectMessage(c, "Register Success", true, "/login")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Values["isLogin"] = false
	sess.Save(c.Request(), c.Response())
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	return redirectMessage(c, "Logout is success", true, "/")
}
func home(c echo.Context) error {
	sess, _ := session.Get("session", c)

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// map data project if no user
	data, _ := connection.Conn.Query(context.Background(), "SELECT id_project, name_project, start_date, end_date, desc_project,technologies, duration, author, images FROM tb_project ORDER BY id_project DESC")

	var result []Project
	for data.Next() {
		var each = Project{}
		var techArray pgtype.VarcharArray
		err := data.Scan(&each.Id, &each.Title, &each.Start, &each.End, &each.Desc, &techArray, &each.Duration ,&each.Author,&each.Images)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}
		each.Techno = make([]string, len(techArray.Elements))
		for i, e := range techArray.Elements {
			if e.String == "node" {
				each.Techno[i] = "node"
			} else if e.String == "react" {
				each.Techno[i] = "react"
			} else if e.String == "next" {
				each.Techno[i] = "next"
			} else if e.String == "typesc" {
				each.Techno[i] = "typesc"
			}
		}
		result = append(result, each)
	}

	// map project if login
	user := sess.Values["username"]
	datas, _ := connection.Conn.Query(context.Background(), "SELECT id_project, name_project, start_date, end_date, desc_project,technologies, duration ,author,images FROM tb_project INNER JOIN tb_user on author = username where username=$1 ORDER BY id_project DESC",user)
	var results []Project
	for datas.Next() {
		var each = Project{}
		var techArray pgtype.VarcharArray
		err := datas.Scan(&each.Id, &each.Title, &each.Start, &each.End, &each.Desc, &techArray, &each.Duration,&each.Author,&each.Images)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}
		each.Techno = make([]string, len(techArray.Elements))
		for i, e := range techArray.Elements {
			if e.String == "node" {
				each.Techno[i] = "node"
			} else if e.String == "react" {
				each.Techno[i] = "react"
			} else if e.String == "next" {
				each.Techno[i] = "next"
			} else if e.String == "typesc" {
				each.Techno[i] = "typesc"
			}
		}
		results = append(results, each)
	}
	var totalresult  []Project
	if sess.Values["isLogin"] != true {
		totalresult = result
	}else{
		totalresult = results
	}

	fmt.Println(results,"'result'")
	projets := map[string]interface{}{
		"Project":      totalresult,
		"StatusLogin":  sess.Values["isLogin"],
		"Username":     sess.Values["username"],
		"FlashMessage": sess.Values["message"],
		"StatusAlert":  sess.Values["statusAlert"],
	}
	delete(sess.Values, "message")
	delete(sess.Values, "statusAlert")
	sess.Save(c.Request(), c.Response())
	return tmpl.Execute(c.Response(), projets)
}

func formproject(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Username = sess.Values["username"].(string)
	}

	flash := map[string]interface{}{
		"DataSession": userData,
	}
	tmpl, err := template.ParseFiles("views/my-project.html")
	if err != nil {
		return redirectMessage(c, err.Error(), false, "/add-project")
	}
	return tmpl.Execute(c.Response(), flash)
}

func addProject(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal("failed to parse form data",err)
	}
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	desc := c.FormValue("desc")
	image := c.Get("dataFile").(string)
	if c.FormValue("node") == "on" {
		node = "node"
	}else{
		node = ""
	}
	if c.FormValue("react") == "on" {
		react = "react"
	}else{
		react = ""
	}
	if c.FormValue("next") == "on" {
		next = "next"
	}else{
		next = ""
	}
	if c.FormValue("typesc") == "on" {
		typesc = "typesc"
	}else{
		typesc = ""
	}
	layout := "2006-01-02"
	date1, _ := time.Parse(layout, start)
	date2, _ := time.Parse(layout, end)
	diff := date2.Sub(date1)
	techs := []string{
		node,
		react,
		next,
		typesc,
	}

	sess, _ := session.Get("session", c)
	user:= sess.Values["username"]

	_, err = connection.Conn.Exec(context.Background(), "insert into tb_project(name_project,start_date,end_date,desc_project,technologies,duration,author,images)values($1,$2,$3,$4,$5,$6,$7,$8)", title, start, end, desc, techs, diff.Hours()/24,user,image)

	if err != nil {
		return redirectMessage(c, "Data failed to add", true, "/add-project")
	}
	return redirectMessage(c, "Data has been added", true, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Username = sess.Values["username"].(string)
	}
	tmpl, err := template.ParseFiles("views/edit-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	var editProject = Project{}
	errs := connection.Conn.QueryRow(context.Background(), "SELECT id_project,name_project, start_date, end_date, desc_project,technologies,images FROM tb_project where id_project=$1", id).Scan(&editProject.Id, &editProject.Title, &editProject.Start, &editProject.End, &editProject.Desc, &editProject.Techno ,&editProject.Images)
	if errs != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errs.Error()})
	}

	data := map[string]interface{}{
		"Edit":     editProject,
		"Contains": contains,
		"DataSession":userData,
	}
	return tmpl.Execute(c.Response(), data)
}

func editFunc(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal("failed to parse form data",err)
	}
	id, _ := strconv.Atoi(c.Param("id"))
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	desc := c.FormValue("desc")
	image := c.Get("dataFile").(string)
	if c.FormValue("node") == "on" {
		node = "node"
	} else {
		node = ""
	}
	if c.FormValue("react") == "on" {
		react = "react"
	} else {
		react = ""
	}
	if c.FormValue("next") == "on" {
		next = "next"
	} else {
		next = ""
	}
	if c.FormValue("typesc") == "on" {
		typesc = "typesc"
	} else {
		typesc = ""
	}
	layout := "2006-01-02"
	date1, _ := time.Parse(layout, start)
	date2, _ := time.Parse(layout, end)
	diff := date2.Sub(date1)
	techs := []string{
		node,
		react,
		next,
		typesc,
	}
	fmt.Println(techs)
	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name_project=$1, start_date=$2,end_date=$3, desc_project=$4, technologies=$5, duration=$6,images=$7 WHERE id_project=$8;", title, start, end, desc, techs, diff.Hours()/24,image, id)
	if err != nil {
		return redirectMessage(c, "Data failed to updated", false, "/edit-project/{{id}}")
	}
	return redirectMessage(c, "Data has been updated", true, "/")
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Username = sess.Values["username"].(string)
	}
	tmpl, err := template.ParseFiles("views/detail-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var detailProject = Project{}
	errs := connection.Conn.QueryRow(context.Background(), "SELECT id_project,name_project,start_date,end_date,desc_project,technologies,duration,images FROM tb_project WHERE id_project=$1", id).Scan(&detailProject.Id, &detailProject.Title, &detailProject.Start, &detailProject.End, &detailProject.Desc, &detailProject.Techno, &detailProject.Duration,&detailProject.Images)
	if errs != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data := map[string]interface{}{
		"Project": detailProject,
		"DataSession":userData,
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id_project=$1", id)

	if err != nil {
		return redirectMessage(c, "Data failed to delete", false, "/add-project")
	}
	return redirectMessage(c, "Data has been deleted", true, "/")
}

func contactme(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Username = sess.Values["username"].(string)
	}

	flash := map[string]interface{}{
		"DataSession": userData,
	}

	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(),flash)
}

func redirectMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}
