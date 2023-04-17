package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Id       int
	Title    string
	Start    string
	End      string
	Desc     string
	Node     bool
	React    bool
	Next     bool
	Typesc   bool
	Duration int
}

var dataProject = []Project{
	{
		Title:    "Sapi terbang jualan kuah sayur",
		Start:    "2023-04-15",
		End:      "2023-06-21",
		Desc:     "Lorem ipsum dolor, sit amet consectetur adipisicing elit. Quia temporibus mollitia impedit repudiandae? Non, assumenda! Dignissimos sed, accusamus aspernatur itaque at ut perferendis consequatur dolorem ad temporibus corrupti molestiae id, iusto saepe sint iure. Praesentium assumenda, eveniet placeat explicabo perferendis, consequuntur nesciunt optio corrupti quaerat error porro officiis. Cupiditate, vitae!",
		Node:     true,
		React:    true,
		Next:     false,
		Typesc:   true,
		Duration: 234,
	},
	{
		Title:    "Kuda terbang jualan kuah sayur",
		Start:    "2023-02-15",
		End:      "2023-08-21",
		Desc:     "Lorem ipsum dolor, sit amet consectetur adipisicing elit. Quia temporibus mollitia impedit repudiandae? Non, assumenda! Dignissimos sed, accusamus aspernatur itaque at ut perferendis consequatur dolorem ad temporibus corrupti molestiae id, iusto saepe sint iure. Praesentium assumenda, eveniet placeat explicabo perferendis, consequuntur nesciunt optio corrupti quaerat error porro officiis. Cupiditate, vitae!",
		Node:     false,
		React:    true,
		Next:     false,
		Typesc:   true,
		Duration: 10,
	},
}

func main() {
	e := echo.New()
	e.Static("/public", "public")
	//Routing
	e.GET("/", home)
	e.GET("/contact-me", contactme)
	e.GET("/add-project", formproject)
	e.POST("/add", addProject)
	e.GET("/detail-project/:id", detailproject)
	e.GET("/delete-project/:id", deleteProject)
	e.GET("/edit-project/:id", editProject)
	e.POST("/edit/:id", editFunc)

	e.Logger.Fatal(e.Start("localhost:5000"))
	
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	projets := map[string]interface{}{
		"Project": dataProject,
	}
	
	return tmpl.Execute(c.Response(), projets)
}

func formproject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/my-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func addProject(c echo.Context) error {
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	desc := c.FormValue("desc")
	node := c.FormValue("node") == "on"
	react := c.FormValue("react") == "on"
	next := c.FormValue("next") == "on"
	typesc := c.FormValue("typesc") == "on"
	layout := "2006-01-02"
	date1, _ := time.Parse(layout, start)
	date2, _ := time.Parse(layout, end)
	diff := date2.Sub(date1)

	var addProject = Project{
		Title:    title,
		Desc:     desc,
		Start:    start,
		End:      end,
		Node:     node,
		React:    react,
		Next:     next,
		Typesc:   typesc,
		Duration: int(diff.Hours()),
	}

	fmt.Println(start, end, diff)
	dataProject = append(dataProject, addProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tmpl, err := template.ParseFiles("views/edit-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	var editProject = Project{}
	for i, data := range dataProject {
		if id == i {
			editProject = Project{
				Id: id,
				Title:  data.Title,
				Desc:   data.Desc,
				Start:  data.Start,
				End:    data.End,
				Node:   data.Node,
				React:  data.React,
				Next:   data.Next,
				Typesc: data.Typesc,
			}
		}
	}
	data := map[string]interface{}{
		"Edit": editProject,
	}
	return tmpl.Execute(c.Response(), data)
}

func editFunc(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	desc := c.FormValue("desc")
	node := c.FormValue("node") == "on"
	react := c.FormValue("react") == "on"
	next := c.FormValue("next") == "on"
	typesc := c.FormValue("typesc") == "on"

	layout := "2006-01-02"
	date1, _ := time.Parse(layout, start)
	date2, _ := time.Parse(layout, end)
	diff := date2.Sub(date1)
	dataProject[id].Title = title
	dataProject[id].Desc = desc
	dataProject[id].Start = start
	dataProject[id].End = end
	dataProject[id].Node = node
	dataProject[id].React = react
	dataProject[id].Next = next
	dataProject[id].Typesc = typesc
	dataProject[id].Duration = int(diff.Hours())
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tmpl, err := template.ParseFiles("views/detail-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	var detailProject = Project{}

	for i, data := range dataProject {
		if id == i {
			detailProject = Project{
				Title:    data.Title,
				Desc:     data.Desc,
				Start:    data.Start,
				End:      data.End,
				Node:     data.Node,
				React:    data.React,
				Next:     data.Next,
				Typesc:   data.Typesc,
				Duration: data.Duration,
			}
		}
	}
	data := map[string]interface{}{
		"Project": detailProject,
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProject = append(dataProject[:id], dataProject[id+1:]...)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func contactme(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}
