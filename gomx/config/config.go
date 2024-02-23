package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var AppRootDir string
var ApiRootDir string
var RoutesDir string
var ReservedDir string
var BaseTemplate string

type Config struct {
	AppRootDir   string `json:"appRoot"`
	ApiRootDir   string `json:"apiRoot"`
	RoutesDir    string `json:"routes"`
	ReservedDir  string `json:"reserved"`
	BaseTemplate string `json:"baseTemplate"`
}

var DefaultConfig = Config{
	AppRootDir:   "./app",
	ApiRootDir:   "./api",
	RoutesDir:    "/routes",
	ReservedDir:  "/_",
	BaseTemplate: "/index.gohtml",
}

func init() {
	data, err := os.ReadFile("gomx.config.json")
	if err != nil {
		log.Println("Error reading config file.")
		log.Println(err)
		log.Println("Using default configuration.")
		return
	}
	var c = &DefaultConfig
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Println("Error reading config file.")
		log.Println(err)
		log.Println("Using default configuration.")
		return
	}

	log.Println("Using config:")

	if c.AppRootDir != "" {
		AppRootDir = c.AppRootDir
	}
	AppRootDir = filepath.ToSlash(filepath.Clean(AppRootDir))
	log.Printf("\"appRoot\" = %s\n", AppRootDir)
	if c.ApiRootDir != "" {
		ApiRootDir = c.ApiRootDir
	}
	ApiRootDir = filepath.Join(AppRootDir, ApiRootDir)
	ApiRootDir = filepath.ToSlash(filepath.Clean(ApiRootDir))
	log.Printf("\"apiRoot\" = %s\n", ApiRootDir)
	if c.RoutesDir != "" {
		RoutesDir = c.RoutesDir
	}
	RoutesDir = filepath.Join(AppRootDir, RoutesDir)
	RoutesDir = filepath.ToSlash(filepath.Clean(RoutesDir))
	log.Printf("\"routes\" = %s\n", RoutesDir)
	if c.ReservedDir != "" {
		ReservedDir = c.ReservedDir
	}
	ReservedDir = filepath.Join(AppRootDir, RoutesDir, ReservedDir)
	ReservedDir = filepath.ToSlash(filepath.Clean(ReservedDir))
	log.Printf("\"reserved\" = %s\n", ReservedDir)
	if c.BaseTemplate != "" {
		BaseTemplate = c.BaseTemplate
	}
	BaseTemplate = filepath.Join(AppRootDir, BaseTemplate)
	BaseTemplate = filepath.ToSlash(filepath.Clean(BaseTemplate))
	log.Printf("\"baseTemplate\" = %s\n", BaseTemplate)
}

func ReadConfigFile() *Config {
	data, err := os.ReadFile("gomx.config.json")
	if err != nil {
		log.Println("Error reading config file.")
		log.Println("Missing 'gomx.config.json' file.")
		log.Println("Using default configuration.")
		return &DefaultConfig
	}
	var c *Config
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Println("Error reading config file.")
		log.Println("Missing 'gomx.config.json' file.")
		log.Println("Using default configuration.")
		return &DefaultConfig
	}

	log.Println("Using config:")

	if c.AppRootDir == "" {
		c.AppRootDir = DefaultConfig.AppRootDir
	}
	c.AppRootDir = filepath.ToSlash(filepath.Clean(c.AppRootDir))
	log.Printf("\"appRoot\" = %s\n", c.AppRootDir)
	if c.ApiRootDir == "" {
		c.ApiRootDir = DefaultConfig.ApiRootDir
	}
	c.ApiRootDir = filepath.Join(c.AppRootDir, c.ApiRootDir)
	c.ApiRootDir = filepath.ToSlash(filepath.Clean(c.ApiRootDir))
	log.Printf("\"apiRoot\" = %s\n", c.ApiRootDir)
	if c.RoutesDir == "" {
		c.RoutesDir = DefaultConfig.RoutesDir
	}
	c.RoutesDir = filepath.Join(c.AppRootDir, c.RoutesDir)
	c.RoutesDir = filepath.ToSlash(filepath.Clean(c.RoutesDir))
	log.Printf("\"routes\" = %s\n", c.RoutesDir)
	if c.ReservedDir == "" {
		c.ReservedDir = DefaultConfig.ReservedDir
	}
	c.ReservedDir = filepath.Join(c.AppRootDir, c.RoutesDir, c.ReservedDir)
	c.ReservedDir = filepath.ToSlash(filepath.Clean(c.ReservedDir))
	log.Printf("\"reserved\" = %s\n", c.ReservedDir)
	if c.BaseTemplate == "" {
		c.BaseTemplate = DefaultConfig.BaseTemplate
	}
	c.BaseTemplate = filepath.Join(c.AppRootDir, c.BaseTemplate)
	c.BaseTemplate = filepath.ToSlash(filepath.Clean(c.BaseTemplate))
	log.Printf("\"baseTemplate\" = %s\n", c.BaseTemplate)
	return c
}
