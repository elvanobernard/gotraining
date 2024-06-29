package main

import (
	"fmt"
	"log"
	"net/http"

	"wolung/reusable/auth"
	"wolung/reusable/database"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

type paxStruct struct {
	PaxName string `form:"pax_name"`
	PaxAge  int    `form:"pax_age"`
}

func (pax paxStruct) Info() string {
	return pax.PaxName + " " + fmt.Sprint(pax.PaxAge)
}

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

func GetDataC(c *gin.Context) {
	var b StructC
	c.Bind(&b)
	c.JSON(200, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

func GetDataD(c *gin.Context) {
	var b StructD
	c.Bind(&b)
	c.JSON(200, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})
}

type UserDataStruct struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func SignUp(c *gin.Context) {
	var userData UserDataStruct
	c.Bind(&userData)

	db := database.GetDb("root", "", "tcp", "127.0.0.1:3306", "gotest")

	result, err := auth.NewUser(db, userData.UserName, userData.Password)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"result":   result,
		"username": userData.UserName,
		"password": userData.Password + "received",
	})
}

func LogIn(c *gin.Context) {
	var userData UserDataStruct
	c.Bind(&userData)

	db := database.GetDb("root", "", "tcp", "127.0.0.1:3306", "gotest")

	result, err := auth.TryLogin(db, userData.UserName, userData.Password)
	fmt.Println(err)

	if err != nil {
		c.JSON(401, gin.H{
			"message": err,
		})
	} else if !result {
		c.JSON(401, gin.H{
			"result": "incorrect password",
		})
	} else {
		c.JSON(200, gin.H{
			"result": "success log in",
		})
	}

}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	r.GET("/pax", func(c *gin.Context) {
		var pax paxStruct
		c.Bind(&pax)
		c.JSON(200, pax.Info())
	})

	r.GET("/getb", GetDataB)
	r.GET("/getc", GetDataC)
	r.GET("/getd", GetDataD)
	r.POST("/sign-up", SignUp)
	r.POST("/log-in", LogIn)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
