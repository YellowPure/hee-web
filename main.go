package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

// func ListenAndServe(address string, handler Handler) error

// func indexHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "Hello Path=%s", req.URL.Path)
// }

// func helloHandler(w http.ResponseWriter, req *http.Request) {
// 	for k, v := range req.Header {
// 		fmt.Fprintf(w, "req header[%q]=%q\n", k, v)
// 	}
// }

// func main() {
// 	r := gee.New()
// 	// r.GET("/", func(w http.ResponseWriter, req *http.Request) {
// 	// 	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
// 	// })

// 	// r.POST("/hello", func(w http.ResponseWriter, req *http.Request) {
// 	// 	for k, v := range req.Header {
// 	// 		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
// 	// 	}
// 	// })

// 	r.GET("/", func(c *gee.Context) {
// 		c.HTML(http.StatusOK, "<h1>hello world</h1>")
// 	})

// 	r.GET("/hello", func(c *gee.Context) {
// 		// expect /hello?name=geektutu
// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
// 	})

// 	r.POST("/login", func(c *gee.Context) {
// 		c.JSON(http.StatusOK, gee.H{
// 			"username": c.PostForm("username"),
// 			"password": c.PostForm("password"),
// 		})
// 	})

// 	r.Run(":9999")
// 	// http.HandleFunc("/", indexHandler)
// 	// http.HandleFunc("/hello", helloHandler)
// 	// log.Fatal(http.ListenAndServe(":9999", engine))
// }

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func test() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()
	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after")
}

func main() {

	r := gee.Default()
	// r.GET("/", func(c *gee.Context) {
	// 	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	// })

	// r.GET("/hello", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	// })
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"ttttyelow"}
		c.String(http.StatusOK, names[100])
	})
	v1 := r.Group("/v1")
	{

		v1.GET("/hello/:name", func(c *gee.Context) {
			// c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9999")
}
