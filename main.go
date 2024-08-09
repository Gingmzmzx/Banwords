package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"github.com/zheng-ji/goAcAutoMachine"
)

type Server struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	User string `json:"user"`
	Key  string `json:"key"`
}

type Data struct {
	Path string `json:"path"`
}

type Config struct {
	Server Server `json:"server"`
	Data   Data   `json:"data"`
}

var logger = go_logger.NewLogger()

var ac = goAcAutoMachine.NewAcAutoMachine()

func loadConfig(configFile string) Config {
	logger.Info("Loading config from " + configFile)
	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	return config
}

func setupAcAM(dataPath string) {
	logger.Info("Loading data from " + dataPath)
	f, err := os.Open(dataPath)
	if err != nil {
		logger.Error(err.Error())
	}
	defer f.Close()

	logger.Info("Adding Patterns...")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// logger.Info("ACAutoMachine AddPattern: " + scanner.Text())
		ac.AddPattern(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err.Error())
	}

	ac.Build()
}

func setupRouter(user string, key string) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authorized := r.Group("/")
	authorized.Use(gin.BasicAuth(gin.Accounts{
		user: key,
	}))

	authorized.GET("check/:str", func(c *gin.Context) {
		str := c.Params.ByName("str")
		results := ac.Query(str)

		logger.Info("Run check str: " + str + " result: " + strings.Join(results, ", "))

		if len(results) == 0 {
			c.JSON(http.StatusOK, gin.H{"result": false})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": true, "text": results})
		}
	})

	return r
}

func main() {
	configFile := "config.json"
	if len(os.Args) >= 3 {
		configFile = os.Args[2]
	}
	config := loadConfig(configFile)
	user := config.Server.User
	key := config.Server.Key
	addr := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	logger.Info("Running server with user: " + user + " key: " + key + " on: " + addr)

	setupAcAM(config.Data.Path)
	r := setupRouter(user, key)
	r.Run(addr)
}
