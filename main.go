package main

import (
	"context"
	"fmt"
	"net/http"
	"personal-web-week-3/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
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
}

var node string
var react string
var next string
var typesc string

func main() {
	connection.DatabaseConnect()

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

	// map data project
	data, _ := connection.Conn.Query(context.Background(), "SELECT id_project, name_project, start_date, end_date, desc_project,technologies, duration FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}
		var techArray pgtype.VarcharArray
		err := data.Scan(&each.Id, &each.Title, &each.Start, &each.End, &each.Desc, &techArray, &each.Duration)
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
	alert := fmt.Sprintf("<script>alert('Input data berhasil!')</script>")

	projets := map[string]interface{}{
		"Project": result,
		"message":alert,
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
	if c.FormValue("node") == "on" {
		node = "node"
	}
	if c.FormValue("react") == "on" {
		react = "react"
	}
	if c.FormValue("next") == "on" {
		next = "next"
	}
	if c.FormValue("typesc") == "on" {
		typesc = "typesc"
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
	_, err := connection.Conn.Exec(context.Background(), "insert into tb_project(name_project,start_date,end_date,desc_project,technologies,duration)values($1,$2,$3,$4,$5,$6)", title, start, end, desc, techs, diff.Hours()/24)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	success := "User data has been submitted successfully."
	c.HTML(http.StatusOK,"<script>alert('" + success + "')")
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	// id, _ := strconv.Atoi(c.Param("id"))
	tmpl, err := template.ParseFiles("views/edit-project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	var editProject = Project{}
	// for i, data := range dataProject {
	// 	if id == i {
	// 		editProject = Project{
	// 			Id:    id,
	// 			Title: data.Title,
	// 			Desc:  data.Desc,
	// 			Start: data.Start,
	// 			End:   data.End,
	// 		}
	// 	}
	// }
	data := map[string]interface{}{
		"Edit": editProject,
	}
	return tmpl.Execute(c.Response(), data)
}

func editFunc(c echo.Context) error {
	// id, _ := strconv.Atoi(c.Param("id"))
	// title := c.FormValue("title")
	// start := c.FormValue("start")
	// end := c.FormValue("end")
	// desc := c.FormValue("desc")

	// layout := "2006-01-02"
	// date1, _ := time.Parse(layout, start)
	// date2, _ := time.Parse(layout, end)
	// diff := date2.Sub(date1)
	// dataProject[id].Title = title
	// dataProject[id].Desc = desc
	// dataProject[id].Start = start
	// dataProject[id].End = end

	// dataProject[id].Duration = int(diff.Hours())
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailproject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tmpl, err := template.ParseFiles("views/detail-project.html")
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var detailProject = Project{}
	errs := connection.Conn.QueryRow(context.Background(), "SELECT id_project,name_project,start_date,end_date,desc_project,technologies,duration FROM tb_project WHERE id_project=$1", id).Scan(&detailProject.Id, &detailProject.Title, &detailProject.Start, &detailProject.End, &detailProject.Desc, &detailProject.Techno, &detailProject.Duration)
	if errs != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// for i, data := range dataProject {
	// 	if id == i {
	// 		detailProject = Project{
	// 			Title:    data.Title,
	// 			Desc:     data.Desc,
	// 			Start:    data.Start,
	// 			End:      data.End,
	// 			Duration: data.Duration,
	// 		}
	// 	}
	// }
	data := map[string]interface{}{
		"Project": detailProject,
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id_project=$1", id)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError,map[string]string{"message":err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func contactme(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}
