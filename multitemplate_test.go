package multitemplate

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createFromFile() Render {
	r := New()
	r.AddFromFiles("index", "tests/base.html", "tests/article.html")

	return r
}

func createFromGlob() Render {
	r := New()
	r.AddFromGlob("index", "tests/global/*")

	return r
}

func createFromFS() Render {
	r := New()
	r.AddFromFS("index", os.DirFS("."), "tests/base.html", "tests/article.html")

	return r
}

func createFromString() Render {
	r := New()
	r.AddFromString("index", "Welcome to {{ .name }} template")

	return r
}

func createFromStringsWithFuncs() Render {
	r := New()
	r.AddFromStringsFuncs(
		"index",
		template.FuncMap{},
		`Welcome to {{ .name }} {{template "content"}}`, `{{define "content"}}template{{end}}`,
	)

	return r
}

func createFromFilesWithFuncs() Render {
	r := New()
	r.AddFromFilesFuncs("index", template.FuncMap{}, "tests/welcome.html", "tests/content.html")

	return r
}

func TestMissingTemplateOrName(t *testing.T) {
	r := New()
	tmpl := template.Must(template.New("test").Parse("Welcome to {{ .name }} template"))
	assert.Panics(t, func() {
		r.Add("", tmpl)
	}, "template name cannot be empty")

	assert.Panics(t, func() {
		r.Add("test", nil)
	}, "template can not be nil")
}

func TestAddFromFiles(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromFile()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Test Multiple Template",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>Test Multiple Template</p>\nHi, this is article template\n", w.Body.String())
}

func TestAddFromGlob(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromGlob()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Test Multiple Template",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>Test Multiple Template</p>\nHi, this is login template\n", w.Body.String())
}

func TestAddFromFS(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromFS()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Test Multiple Template",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>Test Multiple Template</p>\nHi, this is article template\n", w.Body.String())
}

func TestAddFromString(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromString()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"name": "index",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome to index template", w.Body.String())
}

func TestAddFromStringsFruncs(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromStringsWithFuncs()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"name": "index",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome to index template", w.Body.String())
}

func TestAddFromFilesFruncs(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromFilesWithFuncs()
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"name": "index",
		})
	})

	w := performRequest(router)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Welcome to index template\n", w.Body.String())
}

func TestDuplicateTemplate(t *testing.T) {
	assert.Panics(t, func() {
		r := New()
		r.AddFromString("index", "Welcome to {{ .name }} template")
		r.AddFromString("index", "Welcome to {{ .name }} template")
	})
}

func createFromPartialFiles() Render {
	r := New()
	r.AddFromFiles("partials", "tests/partial-base.html", "tests/partials.html")
	return r
}

func TestPartialRendering(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromPartialFiles()
	
	// Test rendering just the header partial
	router.GET("/header", func(c *gin.Context) {
		c.HTML(200, "partials#header", gin.H{
			"name": "TestUser",
		})
	})
	
	// Test rendering just the content partial  
	router.GET("/content", func(c *gin.Context) {
		c.HTML(200, "partials#content", gin.H{
			"name": "TestUser",
		})
	})
	
	// Test rendering just the footer partial
	router.GET("/footer", func(c *gin.Context) {
		c.HTML(200, "partials#footer", gin.H{
			"name": "TestUser",
		})
	})
	
	// Test rendering sidebar partial that's not used in main template
	router.GET("/sidebar", func(c *gin.Context) {
		c.HTML(200, "partials#sidebar", gin.H{
			"name": "TestUser",
		})
	})

	// Test header partial
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/header", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<header>Welcome TestUser</header>", w.Body.String())

	// Test content partial
	req, _ = http.NewRequestWithContext(context.Background(), "GET", "/content", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>This is the main content for TestUser</p>", w.Body.String())

	// Test footer partial
	req, _ = http.NewRequestWithContext(context.Background(), "GET", "/footer", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<footer>© 2024 TestUser</footer>", w.Body.String())

	// Test sidebar partial
	req, _ = http.NewRequestWithContext(context.Background(), "GET", "/sidebar", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<aside>Sidebar content</aside>", w.Body.String())
}

func TestFullTemplateStillWorks(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromPartialFiles()
	
	router.GET("/full", func(c *gin.Context) {
		c.HTML(200, "partials", gin.H{
			"title": "Test Page",
			"name": "TestUser",
		})
	})

	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/full", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	expected := "<!DOCTYPE html>\n<html>\n<head><title>Test Page</title></head>\n<body>\n<h1>Test Page</h1>\n<header>Welcome TestUser</header>\n<main><p>This is the main content for TestUser</p></main>\n<footer>© 2024 TestUser</footer>\n</body>\n</html>"
	assert.Equal(t, expected, w.Body.String())
}

func TestPartialEdgeCases(t *testing.T) {
	router := gin.New()
	router.HTMLRender = createFromPartialFiles()
	
	// Test non-existent partial
	router.GET("/nonexistent", func(c *gin.Context) {
		c.HTML(200, "partials#nonexistent", gin.H{
			"name": "TestUser",
		})
	})
	
	// Test empty partial name
	router.GET("/empty", func(c *gin.Context) {
		c.HTML(200, "partials#", gin.H{
			"name": "TestUser",
		})
	})

	// Test non-existent partial - should return error
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Check if it contains an error response (it should fail during template execution)
	t.Logf("Non-existent partial response code: %d, body: %s", w.Code, w.Body.String())
	
	// Test empty partial name - should render full template
	req, _ = http.NewRequestWithContext(context.Background(), "GET", "/empty", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	expected := "<!DOCTYPE html>\n<html>\n<head><title></title></head>\n<body>\n<h1></h1>\n<header>Welcome TestUser</header>\n<main><p>This is the main content for TestUser</p></main>\n<footer>© 2024 TestUser</footer>\n</body>\n</html>"
	assert.Equal(t, expected, w.Body.String())
}
