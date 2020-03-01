package model

import (
	"time"
)

type Project struct {
	ProjectType        string `json : "projectType"`
	ProjectName        string `json : "projectName"`
	Language           string `json : "language"`
	LanguageVersion    string `json:"languageVersion"`
	IsDbUsed           bool   `json:"dbUsed"`
	DbName             string `json:"dbName"`
	DbVersion          string `json:"dbVersion"`
	ZipPath            string `json:"zipPath"`
	UserName           string `json:"username"`
	DownloadPermission bool   `json :"downloadPermission"`
	Description        string `json:"description"`
	ContainerName      string `json : "containerName"`
	ContainerIP        string `json :"containerIP"`
	Device             string `json : "device"`
	AndroidVersion     string `json:"androidVersion"`
	ApiVersion         string `json:"apiVersion"`

	CreatedDate time.Time `json:"createdDate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

type WebProject struct {
	ProjectType        string `json : "projectType"`
	ProjectName        string `json : "projectName"`
	Language           string `json : "language"`
	LanguageVersion    string `json:"languageVersion"`
	IsDbUsed           bool   `json:"dbUsed"`
	DbName             string `json:"dbName"`
	DbVersion          string `json:"dbVersion"`
	ProjectPath        string `json:"projectPath"`
	ZipPath            string `json:"zipPath"`
	DbBackupName       string `json:"dbBackupName"`
	DbBackupPath       string `json:"dbBackupPath"`
	UserName           string `json:"username"`
	DownloadPermission bool   `json :"downloadPermission"`
	Description        string `json:"description"`
	ContainerName      string `json : "containerName"`
	ContainerIP        string `json :"containerIP"`

	CreatedDate time.Time `json:"createdDate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

type WindowProject struct {
	ProjectType        string `json : "projectType"`
	ProjectName        string `json : "projectName"`
	Language           string `json : "language"`
	LanguageVersion    string `json:"languageVersion"`
	IsDbUsed           bool   `json:"dbUsed"`
	DbName             string `json:"dbName"`
	DbVersion          string `json:"dbVersion"`
	ProjectPath        string `json:"projectPath"`
	ZipPath            string `json:"zipPath"`
	DbBackupName       string `json:"dbBackupName"`
	DbBackupPath       string `json:"dbBackupPath"`
	UserName           string `json:"username"`
	DownloadPermission bool   `json :"downloadPermission"`
	Description        string `json:"description"`
	ContainerName      string `json : "containerName"`
	ContainerIP        string `json :"containerIP"`

	CreatedDate time.Time `json:"createdDate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

type AndroidProject struct {
	ProjectType        string `json : "projectType"`
	ProjectName        string `json : "projectName"`
	Device             string `json : "device"`
	AndroidVersion     string `json:"androidVersion"`
	ApiVersion         string `json:"apiVersion"`
	ProjectPath        string `json:"projectPath"`
	ZipPath            string `json:"zipPath"`
	UserName           string `json:"username"`
	DownloadPermission bool   `json :"downloadPermission"`
	Description        string `json:"description"`
	ContainerName      string `json : "containerName"`
	ContainerIP        string `json :"containerIP"`

	CreatedDate time.Time `json:"createdDate"`
	LastUpdate  time.Time `json:"lastUpdate"`
}
