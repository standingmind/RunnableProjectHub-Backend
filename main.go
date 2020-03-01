package main

import (
	"control"
	"createContainer"
	"helper"
	"model"
	"strconv"

	"archive/zip"
	"bytes"
	"dbutil"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//Compile templates on start

//containerNameCount file path
var databaseDirName = "sql_script"
var passContainerPath = ""

var encryptionKey = "runnable-project-hub"
var loggedUserSession = sessions.NewCookieStore([]byte(encryptionKey))

func init() {

	loggedUserSession.Options = &sessions.Options{
		// change domain to match your machine. Can be localhost
		// IF the Domain name doesn't match, your session will be EMPTY!
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 3, // 3 hours
		HttpOnly: true,
	}
}

//Display the named template
func projectHandler(w http.ResponseWriter, r *http.Request) {

	helper.GetProjects(w, r)

}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	//get username to use as folder
	var containerIP = "http://www.google.com"
	var containerName = ""
	var containerPath string
	var passContainerPath string
	var databasePath string

	var zipName string = ""
	var sourceDir string = ""
	var destDir string = ""

	switch r.Method {
	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Panic(err)
			return
		}
		//get a ref to the parsed multipart form
		fmt.Println("upload handler work")
		m := r.MultipartForm
		var web model.WebProject
		var win model.WindowProject
		var android model.AndroidProject

		if m.Value["projectType"][0] == "Website" {
			fmt.Println("website worked")
			fmt.Println(m.Value)
			//get the *fileheaders
			fmt.Println("website work")
			web.UserName = m.Value["userName"][0]
			web.ProjectType = m.Value["projectType"][0]
			web.ProjectName = m.Value["projectName"][0]
			web.Language = m.Value["language"][0]
			web.LanguageVersion = m.Value["languageVersion"][0]
			web.IsDbUsed, _ = strconv.ParseBool(m.Value["isDbUsed"][0])
			web.DbName = m.Value["dbName"][0]
			web.DbVersion = m.Value["dbVersion"][0]
			web.DbBackupName = m.Value["dbBackupName"][0]
			web.Description = m.Value["description"][0]
			web.DownloadPermission, _ = strconv.ParseBool(m.Value["downloadPermission"][0])

			projectBasePath := "./userproject/" + web.UserName + "/" + web.ProjectName + "/"
			files := m.File["projectFiles"]

			if len(files) == 0 {
				//no files return
				log.Println("no file")
				return
			}
			var i int
			for i = range files {
				file, err := files[i].Open()
				if err != nil {
					log.Panic(err)
					return
				}
				//constructing directory and nested directories
				fmt.Println("file name:", files[i].Filename)
				dir := filepath.Dir(files[i].Filename)
				fmt.Println("dir :", dir)
				absoluteDir := filepath.Join(projectBasePath, strings.Replace(dir, "\\", "/", -1))
				//	./userproject/userName/projectName/dir
				fmt.Println("absoluteDir :", absoluteDir)
				pathErr := os.MkdirAll(absoluteDir, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				//create destination file making sure the path is writeable.
				constructFile(projectBasePath+files[i].Filename, file) //file path + filename, file content
				file.Close()
			}
			fmt.Println("Finished copying ", i+1, " files")

			//container path name
			fileDir := filepath.Dir(files[0].Filename)
			fileDir = strings.Replace(fileDir, "\\", "/", -1)
			parts := strings.Split(fileDir, "/")
			containerPath = projectBasePath + parts[0] //container path is also zip path
			passContainerPath = projectBasePath + parts[0]
			//zip the folder
			sourceDir = containerPath + "/"
			zipName = parts[0] + ".zip"
			destDir = projectBasePath + zipName

			containerName = web.UserName + "_" + strings.Replace(strings.ToLower(web.ProjectName), " ", "", -1)
			//containerName = strings.Replace(tempName, " ", "", -1)
			fmt.Println("containerName :" + containerName)

			web.ContainerName = containerName
			web.ZipPath = zipName

			//display success message.
		} else if m.Value["projectType"][0] == "Window Application" {
			fmt.Println("window worked")
			win.UserName = m.Value["userName"][0]
			win.ProjectType = m.Value["projectType"][0]
			win.ProjectName = m.Value["projectName"][0]
			win.Language = m.Value["language"][0]
			win.LanguageVersion = m.Value["languageVersion"][0]
			win.IsDbUsed, _ = strconv.ParseBool(m.Value["isDbUsed"][0])
			win.DbName = m.Value["dbName"][0]
			win.DbVersion = m.Value["dbVersion"][0]
			win.DbBackupName = m.Value["dbBackupName"][0]
			win.Description = m.Value["description"][0]
			win.DownloadPermission, _ = strconv.ParseBool(m.Value["downloadPermission"][0])

			projectBasePath := "./userproject/" + win.UserName + "/" + win.ProjectName + "/"
			files := m.File["projectFiles"]
			fmt.Println("length of files :", len(files))
			if len(files) != 0 {
				applicationFile, err := files[0].Open()
				if err != nil {
					log.Println(err)
					return
				}
				//application name as dir name (removing extension) and store under it's dir with the same name
				fileName := files[0].Filename
				parts := strings.Split(fileName, ".")

				containerPath = projectBasePath + parts[0]                      //for creating folder
				passContainerPath = projectBasePath + parts[0] + "/" + fileName //for creating folder

				containerName = win.UserName + "_" + strings.Replace(strings.ToLower(win.ProjectName), " ", "", -1)

				fmt.Println("containerName :" + containerName)
				fmt.Println("container Path :", containerPath)
				//container path is application folder path
				pathErr := os.MkdirAll(containerPath, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				fmt.Println("pathName :" + containerPath)
				//create destination file making sure the path is writeable.
				constructFile(containerPath+"/"+files[0].Filename, applicationFile) //file path + filename, file content
				applicationFile.Close()

				//zip the folder
				sourceDir = containerPath + "/"
				zipName = parts[0] + ".zip"
				destDir = projectBasePath + zipName

				win.ContainerName = containerName
				win.ZipPath = zipName
			}
		} else if m.Value["projectType"][0] == "Android Application" {
			fmt.Println("android worked")
			android.UserName = m.Value["userName"][0]
			android.ProjectType = m.Value["projectType"][0]
			android.ProjectName = m.Value["projectName"][0]
			android.Device = m.Value["device"][0]
			android.AndroidVersion = m.Value["androidVersion"][0]
			android.ApiVersion = m.Value["apiVersion"][0]
			android.DownloadPermission, _ = strconv.ParseBool(m.Value["downloadPermission"][0])
			android.Description = m.Value["description"][0]
			projectBasePath := "./userproject/" + android.UserName + "/" + android.ProjectName + "/"

			files := m.File["projectFiles"]
			fmt.Println("length of files :", len(files))
			if len(files) != 0 {
				applicationFile, err := files[0].Open()
				if err != nil {
					log.Println(err)
					return
				}
				//application name as dir name (removing extension) and store under it's dir with the same name
				fileName := files[0].Filename
				parts := strings.Split(fileName, ".")

				containerPath = projectBasePath + parts[0]                      //for creating folder
				passContainerPath = projectBasePath + parts[0] + "/" + fileName //for creating folder

				containerName = android.UserName + "_" + strings.Replace(strings.ToLower(win.ProjectName), " ", "", -1)

				fmt.Println("containerName :" + containerName)

				//container path is dir name since application name is dir name
				pathErr := os.MkdirAll(containerPath, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				fmt.Println("pathName :" + containerPath)
				//create destination file making sure the path is writeable.
				constructFile(containerPath+"/"+files[0].Filename, applicationFile) //file path + filename, file content
				applicationFile.Close()

				//zip the folder
				sourceDir = containerPath + "/"
				zipName = parts[0] + ".zip"
				destDir = projectBasePath + zipName

				android.ContainerName = containerName
				android.ZipPath = zipName
			}
		}
		//get backup file
		if m.Value["projectType"][0] != "Android Application" {
			f := m.File["dbBackupFile"]
			fmt.Println("length:", len(f))
			if len(f) != 0 {
				//only one file so use f[0] to take only one
				backUpFile, err := f[0].Open()
				defer backUpFile.Close()
				if err != nil {
					log.Panic(err)
				}
				databasePath = containerPath + "/" + databaseDirName
				err = os.MkdirAll(databasePath, 0777)
				if err != nil {
					log.Panic(err)
				}
				if f[0].Filename != "" {
					constructFile(databasePath+"/"+f[0].Filename, backUpFile)
				}
			}
			//end of constructing backup
		}
		//create Container
		fmt.Println("pass Container path :", passContainerPath)
		passContainer := passContainerPath[1:len(passContainerPath)]
		t := time.Now()
		if m.Value["projectType"][0] == "Website" {

			err = createContainer.CreateContainer(web.IsDbUsed, web.DbName, web.DbBackupPath, web.DbVersion, web.Language, web.LanguageVersion, passContainer, containerName, web.DbBackupName, "", "", "")
			//if create container success , write container info , increate count
			if err != nil {
				log.Panic(err)
			}
			containerIP = getIP(containerName, web.Language)

			web.ContainerIP = containerIP
			web.DbBackupPath = databasePath
			web.CreatedDate = t
			web.LastUpdate = t

			zipit(sourceDir, destDir)
			//save json to database
			helper.CreateProject(w, web)

		} else if m.Value["projectType"][0] == "Window Application" {

			err = createContainer.CreateContainer(win.IsDbUsed, win.DbName, win.DbBackupPath, win.DbVersion, win.Language, win.LanguageVersion, passContainer, containerName, win.DbBackupName, "", "", "")
			//if create container success , write container info , increate count
			if err != nil {
				log.Panic(err)
			}
			containerIP = getIP(containerName+"bridge", win.Language)

			win.ContainerIP = containerIP
			win.DbBackupPath = databasePath
			win.CreatedDate = t
			win.LastUpdate = t

			zipit(sourceDir, destDir)
			//save json to database
			helper.CreateProject(w, win)
		} else if m.Value["projectType"][0] == "Android Application" {

			err = createContainer.CreateContainer(false, "", "", "", "", "", passContainer, containerName, "", android.Device, android.AndroidVersion, android.ApiVersion)
			//if create container success , write container info , increate count
			if err != nil {
				log.Panic(err)
			}
			containerIP = getIP(containerName, "")

			android.ContainerIP = containerIP
			android.CreatedDate = t
			android.LastUpdate = t

			zipit(sourceDir, destDir)
			//save json to database
			helper.CreateProject(w, android)
		}

		fmt.Println("container IP :", containerIP)
		fmt.Println("upload reach end")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func constructFile(filepath string, file multipart.File) {
	dst, err := os.Create(filepath)
	defer dst.Close()
	if err != nil {
		log.Panic(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		log.Panic(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("copied " + filepath)
}

func getIP(name string, language string) string {
	return control.IPtest(name, language)
}

func main() {
	r := mux.NewRouter()

	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	// // r.HandleFunc("/", homeHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/project", projectHandler)
	r.HandleFunc("/startApplication", startApplicationHandler)
	r.HandleFunc("/download", downloadHandler)
	// r.HandleFunc("/register", registerHandler)
	// r.HandleFunc("/login", loginHandler)
	// r.HandleFunc("/profile ", profileHandler)
	// r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/supportLanguage", supportLanguageHandler)

	http.ListenAndServe(":8080", r)
	fmt.Print("server starting...")
}

func supportLanguageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonFile, err := os.Open("./support/upload.json")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	}
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	message := map[string]interface{}{}
// 	message["Username"] = ""
// 	message["UserID"] = ""
// 	message["Email"] = ""

// 	// check if session is active
// 	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

// 	if session != nil {
// 		message["Username"] = session.Values["username"]
// 		message["UserID"] = session.Values["userid"]
// 		message["Email"] = session.Values["email"]
// 	}
// 	switch r.Method {
// 	case "GET":

// 	}
// }

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
// 	//read from session
// 	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

// 	// remove the username
// 	session.Values["username"] = ""
// 	session.Values["email"] = ""
// 	session.Values["userid"] = ""
// 	err := session.Save(r, w)

// 	if err != nil {
// 		log.Println(err)
// 	}
// 	http.Redirect(w, r, "http://localhost:8080/upload", http.StatusFound)

// }

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	message := map[string]interface{}{}
// 	message["Username"] = ""
// 	message["UserID"] = ""
// 	message["Email"] = ""

// 	// check if session is active
// 	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

// 	if session != nil {
// 		message["Username"] = session.Values["username"]
// 		message["UserID"] = session.Values["userid"]
// 		message["Email"] = session.Values["email"]
// 	}
// 	switch r.Method {
// 	case "GET":

// 	case "POST":
// 		fmt.Println("post register work")
// 		username := r.FormValue("username")
// 		password := r.FormValue("password")
// 		email := r.FormValue("email")
// 		fmt.Println(username, password, email)
// 		if username != "" && password != "" && email != "" {
// 			//save to database
// 			fmt.Println("registering to db")
// 			user := dbutil.User{}
// 			user.Username = username
// 			user.Password = password
// 			user.Email = email
// 			con, err := dbutil.GetConnection()
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			flag, err := dbutil.IsEmailExist(email, con)
// 			if err != nil {
// 				log.Panic(err)
// 			}
// 			if !flag {
// 				err = dbutil.CreateUser(con, user)
// 				if err != nil {
// 					log.Panic(err)
// 				}
// 				message["Success"] = "User Creation Success!"
// 			} else {
// 				err = dbutil.CloseConnection(con)
// 				if err != nil {
// 					log.Panic(err)
// 				}
// 				message["Fail"] = "Email Already Exist!"
// 			}

// 			displayRegister(w, "register", message)
// 		}
// 	}

// }
// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	message := map[string]interface{}{}
// 	message["Username"] = ""
// 	message["Email"] = ""
// 	message["UserID"] = ""

// 	switch r.Method {
// 	case "GET":
// 		displayLogin(w, "login", message)
// 	case "POST":
// 		fmt.Println("login work")

// 		//check if session is active
// 		session, _ := loggedUserSession.Get(r, "authenticated-user-session")
// 		if session != nil {
// 			message["Username"] = session.Values["username"]
// 			message["Email"] = session.Values["email"]
// 			message["UserID"] = session.Values["userid"]

// 		}

// 		email := r.FormValue("email")
// 		password := r.FormValue("password")
// 		//get connection from db
// 		con, err := dbutil.GetConnection()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		if email != "" && password != "" {
// 			ans, err := dbutil.IsUserExist(email, password, con)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			if !ans {
// 				message["Fail"] = "Email or Password incorrect."
// 				message["LoginError"] = true
// 				displayLogin(w, "login", message)
// 			} else {
// 				userId, userName, err := dbutil.GetUserByEmail(email, con)
// 				if err != nil {
// 					log.Println(err)
// 					return
// 				}
// 				message["Username"] = userName
// 				message["UserID"] = userId
// 				message["Email"] = email
// 				// create a new session and redirect to dashboard
// 				session, _ := loggedUserSession.New(r, "authenticated-user-session")

// 				session.Values["username"] = userName
// 				session.Values["userid"] = userId
// 				session.Values["email"] = email
// 				err = session.Save(r, w)

// 				if err != nil {
// 					log.Println(err)
// 				}

// 				http.Redirect(w, r, "http://localhost:8080/upload   ", http.StatusFound)
// 			}
// 		}

// 	}

// }

// func profileHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		displayProfile(w, "profile", "")
// 	}

// }

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	var zipPath = "./userproject/"
	var zipName string = ""
	var projectName = ""
	var userName = ""
	fmt.Println("enter handler")

	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.MultipartForm)
	m := r.MultipartForm
	zipName = m.Value["zipName"][0]
	userName = m.Value["userName"][0]
	projectName = m.Value["projectName"][0]
	fmt.Println("i worked :", zipName, userName, projectName)
	if zipName != "" {
		data, err := ioutil.ReadFile(zipPath + userName + "/" + projectName + "/" + zipName)
		if err != nil {
			log.Println(err)
			helper.GetError(err, w)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+zipName)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeContent(w, r, zipPath+zipName, time.Now(), bytes.NewReader(data))
	}

}

func startApplicationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r :")
	fmt.Println("start application work")
	var containername string
	var projecttype string
	if r.Method == "POST" {

		containername = r.FormValue("containerName")
		projecttype = r.FormValue("projectType")

		if projecttype == "Android" {
			fmt.Println("Installing apk ...")
			control.InstallApk(containername)
			fmt.Println("Finishing installing ...")
		} else if projecttype == "Website" {

			control.StopContainer(containername)
			fmt.Println("stop container...")
			control.StartContainer(containername)
			fmt.Println("start container...")
		} else {
			control.StartContainer(containername)
			fmt.Println("layering window container ...")
		}

	}
	w.WriteHeader(200)
	fmt.Println("containerName :" + containername + ", projectType :" + projecttype)
	fmt.Fprintln(w, "containerName :"+containername+", projectType :"+projecttype)
}

func getProjectInfo(pattern string, project_type string, language string) (dbutil.ProjectInfos, error) {

	projectInfos := dbutil.ProjectInfos{}
	con, err := dbutil.GetConnection()
	if err != nil {
		log.Println(err)
		return projectInfos, err
	}

	infos, err := dbutil.GetProjectInfo(con, pattern, project_type, language)
	if err != nil {
		log.Println(err)
		return projectInfos, err
	}
	dbutil.CloseConnection(con)
	projectInfos.Infos = infos
	return projectInfos, nil
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
