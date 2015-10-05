package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/unrolled/render"
	"html/template"
	"net/http"
)

type HomeScreen struct {
	Carousel []Video
	Latest   []Video
}

type Playlist struct {
	Id     int
	Videos []Video
}

type Video struct {
	Id       string
	Headline string
	Section  struct {
		Display_name string
	}
	Summary  string
	Images   []Image
	Playlist struct {
		Id string
	}
}

type VideoDetail struct {
	Id           int
	Display_name string
	Headline     string
	Summary      string
	Byline       string
	Images       []Image
	Renditions   []Rendition
}

type Image struct {
	Width  int
	Height int
	Url    string
}

type Rendition struct {
	Height      int
	Url         string
	Video_codec string
}

type PageDetail struct {
	VideoDetail VideoDetail
	Playlist    Playlist
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	r := render.New(render.Options{
		UnEscapeHTML:    true,
		IsDevelopment:   true,
		HTMLContentType: "application/javascript",
		Funcs: []template.FuncMap{
			{
				"getImage": func(images []Image, index int) string {
					return images[index].Url
				},
				"getVideoURL": func(renditions []Rendition) string {
					for _, element := range renditions {
						if element.Height == 1080 && element.Video_codec == "H264" {
							return element.Url
						}
					}
					return ""
				},
			},
		},
	})

	e := echo.New()

	e.Get("/main", func(c *echo.Context) error {
		r.HTML(c.Response().Writer(), http.StatusOK, "main", nil)
		return nil
	})
	e.Get("/home", func(c *echo.Context) error {
		latest := Playlist{}
		getJson("http://www.nytimes.com/svc/video/api/playlist/1194811622182", &latest)
		homePage := HomeScreen{
			Carousel: latest.Videos[0:5],
			Latest:   latest.Videos[5:20],
		}
		fmt.Println(len(homePage.Carousel))
		r.HTML(c.Response().Writer(), http.StatusOK, "home", homePage)
		return nil
	})

	e.Get("/video/:id/:pid", func(c *echo.Context) error {
		video := VideoDetail{}
		getJson("http://www.nytimes.com/svc/video/api/video/"+c.P(0), &video)
		playlist := Playlist{}
		getJson("http://www.nytimes.com/svc/video/api/playlist/"+c.P(1), &playlist)
		pageDetail := PageDetail{
			VideoDetail: video, Playlist: playlist,
		}
		fmt.Println("http://www.nytimes.com/svc/video/api/playlist/" + c.P(1))
		r.HTML(c.Response().Writer(), http.StatusOK, "video", pageDetail)
		return nil
	})
	e.Static("/js", "static/js")
	e.Run(":3000")
}