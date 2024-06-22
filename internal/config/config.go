package config

import (
	"encoding/json"
	"fmt"
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

func Init() {
	data, err := os.ReadFile("gomx.config.json")
	if err != nil {
		fmt.Println("Error reading config file.")
		log.Println(err)
		fmt.Println("Using default configuration.")
		return
	}
	var c = &DefaultConfig
	err = json.Unmarshal(data, c)
	if err != nil {
		fmt.Println("Error reading config file.")
		log.Println(err)
		fmt.Println("Using default configuration.")
		return
	}

	fmt.Println("Using config:")

	if c.AppRootDir != "" {
		AppRootDir = c.AppRootDir
	}
	AppRootDir = filepath.ToSlash(filepath.Clean(AppRootDir))
	fmt.Printf("\"appRoot\" = %s\n", AppRootDir)
	if c.ApiRootDir != "" {
		ApiRootDir = c.ApiRootDir
	}
	ApiRootDir = filepath.Join(AppRootDir, ApiRootDir)
	ApiRootDir = filepath.ToSlash(filepath.Clean(ApiRootDir))
	fmt.Printf("\"apiRoot\" = %s\n", ApiRootDir)
	if c.RoutesDir != "" {
		RoutesDir = c.RoutesDir
	}
	RoutesDir = filepath.Join(AppRootDir, RoutesDir)
	RoutesDir = filepath.ToSlash(filepath.Clean(RoutesDir))
	fmt.Printf("\"routes\" = %s\n", RoutesDir)
	if c.ReservedDir != "" {
		ReservedDir = c.ReservedDir
	}
	ReservedDir = filepath.Join(RoutesDir, ReservedDir)
	ReservedDir = filepath.ToSlash(filepath.Clean(ReservedDir))
	fmt.Printf("\"reserved\" = %s\n", ReservedDir)
	if c.BaseTemplate != "" {
		BaseTemplate = c.BaseTemplate
	}
	BaseTemplate = filepath.Join(AppRootDir, BaseTemplate)
	BaseTemplate = filepath.ToSlash(filepath.Clean(BaseTemplate))
	fmt.Printf("\"baseTemplate\" = %s\n", BaseTemplate)
}
