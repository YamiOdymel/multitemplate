package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create renderer
	r := multitemplate.NewRenderer()
	
	// Add template with multiple partials
	r.AddFromString("example", `
<!DOCTYPE html>
<html>
<head><title>{{.title}}</title></head>
<body>
	<h1>{{.title}}</h1>
	{{template "nav" .}}
	<main>{{template "content" .}}</main>
	{{template "footer" .}}
</body>
</html>

{{define "nav"}}<nav>Navigation for {{.user}}</nav>{{end}}
{{define "content"}}<p>Content area for {{.user}}</p>{{end}}
{{define "footer"}}<footer>&copy; 2024 {{.user}}</footer>{{end}}
{{define "sidebar"}}<aside>Sidebar for {{.user}}</aside>{{end}}
`)

	// Setup Gin router
	router := gin.Default()
	router.HTMLRender = r
	
	// Full template
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "example", gin.H{
			"title": "Full Template",
			"user":  "TestUser",
		})
	})
	
	// Just navigation partial
	router.GET("/nav", func(c *gin.Context) {
		c.HTML(200, "example#nav", gin.H{
			"user": "TestUser",
		})
	})
	
	// Just content partial
	router.GET("/content", func(c *gin.Context) {
		c.HTML(200, "example#content", gin.H{
			"user": "TestUser",
		})
	})
	
	// Just footer partial  
	router.GET("/footer", func(c *gin.Context) {
		c.HTML(200, "example#footer", gin.H{
			"user": "TestUser",
		})
	})
	
	// Sidebar partial (not used in main template)
	router.GET("/sidebar", func(c *gin.Context) {
		c.HTML(200, "example#sidebar", gin.H{
			"user": "TestUser",
		})
	})
	
	fmt.Println("Demo server starting on :8080")
	fmt.Println("Try these URLs:")
	fmt.Println("  http://localhost:8080/        - Full template")
	fmt.Println("  http://localhost:8080/nav     - Navigation partial only")
	fmt.Println("  http://localhost:8080/content - Content partial only")
	fmt.Println("  http://localhost:8080/footer  - Footer partial only")
	fmt.Println("  http://localhost:8080/sidebar - Sidebar partial only")
	
	log.Fatal(router.Run(":8080"))
}