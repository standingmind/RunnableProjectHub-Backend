package main

import (
	"control"
	"createContainer"

	"archive/zip"
	"bytes"
	"dbutil"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//Compile templates on start
var uploadTemplate = template.Must(template.ParseFiles("./template/upload.html"))
var projectTemplate = template.Must(template.ParseFiles("./template/projectList.html"))
var homeTemplate = template.Must(template.ParseFiles("./template/index.html"))
var loginTemplate = template.Must(template.ParseFiles("./template/login.html"))
var registerTemplate = template.Must(template.ParseFiles("./template/register.html"))
var profileTemplate = template.Must(template.ParseFiles("./template/profile.html"))

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
func displayUpload(w http.ResponseWriter, tmpl string, data interface{}) {
	fmt.Println("display upload worked.")
	uploadTemplate.ExecuteTemplate(w, tmpl+".html", data)

}

func displayProject(w http.ResponseWriter, tmpl string, data interface{}) {
	fmt.Println("display project worked")
	projectTemplate.ExecuteTemplate(w, tmpl+".html", data)
}

func displayLogin(w http.ResponseWriter, tmpl string, data interface{}) {
	//fmt.Println("display login worked")
	loginTemplate.ExecuteTemplate(w, tmpl+".html", data)
}
func displayRegister(w http.ResponseWriter, tmpl string, data interface{}) {
	//fmt.Println("display Register worked")
	registerTemplate.ExecuteTemplate(w, tmpl+".html", data)
}
func displayProfile(w http.ResponseWriter, tmpl string, data interface{}) {
	//fmt.Println("display profile worked")
	profileTemplate.ExecuteTemplate(w, tmpl+".html", data)
}

func displayHome(w http.ResponseWriter, tmpl string, data interface{}) {
	//fmt.Println("display home worked")
	homeTemplate.ExecuteTemplate(w, tmpl+".html", data)
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]interface{}{}
	message["Username"] = ""
	message["UserID"] = ""
	message["Email"] = ""

	// check if session is active
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	if session != nil {
		message["Username"] = session.Values["username"]
		message["UserID"] = session.Values["userid"]
		message["Email"] = session.Values["email"]
	}
	//get projectList
	//fmt.Print("project handler work")
	switch r.Method {
	case "GET":

		var pattern = ""
		var projectType = ""
		var language = ""
		if len(r.URL.Query().Get("pattern")) != 0 {
			pattern = r.URL.Query().Get("pattern")
		}
		if len(r.URL.Query().Get("projectType")) != 0 {
			projectType = r.URL.Query().Get("projectType")
		}
		if len(r.URL.Query().Get("language")) != 0 {
			language = r.URL.Query().Get("language")
		}

		fmt.Println("get method work")
		if projectType == "" {
			projectType = "All"
		}
		if language == "" {
			language = "All"
		}

		data, _ := getProjectInfo(pattern, projectType, language)
		if message["Email"] != nil && message["UserID"] != nil && message["Username"] != nil {
			data.Email = fmt.Sprintf("%v", message["Email"])
			data.UserID = fmt.Sprintf("%v", message["UserID"])
			data.Username = fmt.Sprintf("%v", message["Username"])
		}

		// data.ContainerExist = func(name string) bool {
		// 	return checkContainerExist(name)
		// }
		// data.ContainerStart = func(name string) {
		// 	startContainer(name)
		// }
		// data.ContainerIP = func(name string) string {
		// 	return getContainerIP(name)
		// }
		// data.ContainerStop = func(name string) {
		// 	stopContainer(name)
		// }

		fmt.Println("All INformation:", data)
		displayProject(w, "projectList", data)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	//displayProject(w, "projectList", "")
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]interface{}{}
	message["Username"] = ""
	message["UserID"] = ""
	message["Email"] = ""

	// check if session is active
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	if session != nil {
		message["Username"] = session.Values["username"]
		message["UserID"] = session.Values["userid"]
		message["Email"] = session.Values["email"]
	}

	//get connection from db
	con, err := dbutil.GetConnection()
	if err != nil {
		log.Println(err)
		return
	}
	var lastId int
	if lastId, err = dbutil.GetLastId(con); err != nil {
		log.Println(err)
		return
	}
	var userId = fmt.Sprintf("%v", message["UserID"])
	fmt.Println("last Id :", lastId)
	var projectId = lastId
	var containerIP = "http://www.google.com"
	var containerName = ""
	var containerPath string
	var passContainerPath string
	var databasePath string
	var projectRootDir = "./userproject/" + userId + "/" + strconv.Itoa(projectId) + "/"

	var zipName string = ""
	var sourceDir string = ""
	var destDir string = ""

	switch r.Method {
	//GET displays the upload form.
	case "GET":
		//fmt.Println("get work")
		dbutil.CloseConnection(con)
		displayUpload(w, "upload", message)

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		//fmt.Println("post work")
		//parse the multipart form in the request

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Panic(err)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		pj := dbutil.Project{}
		pj.UserID, _ = strconv.Atoi(userId)
		pj.Title = m.Value["project_name"][0]
		pj.Description = m.Value["project_description"][0]
		pj.ProjectType = m.Value["project_type"][0]
		if len(m.Value["permission"]) == 0 {
			pj.Permission = "private"
		} else {
			pj.Permission = "public"
		}

		if pj.ProjectType == "Website" {
			//get the *fileheaders
			pj.Language = m.Value["language"][0]
			pj.LanguageVersion = m.Value["language_version"][0]
			pj.Database = m.Value["database"][0]
			pj.DatabaseVersion = m.Value["database_version"][0]
			pj.DatabaseName = m.Value["database_name"][0]

			files := m.File["projectFolder"]

			if len(files) == 0 {
				displayUpload(w, "upload", "please upload your project folder !")
				return
			}
			var i int
			for i = range files {
				//for each fileheader, get a handle to the actual file
				file, err := files[i].Open()

				if err != nil {
					log.Panic(err)
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				//make nested directories
				dir := filepath.Dir(files[i].Filename)
				absoluteDir := filepath.Join(projectRootDir, strings.Replace(dir, "\\", "/", -1))
				//	./userproject/dir/dir
				fmt.Println("absoluteDir :", absoluteDir)
				pathErr := os.MkdirAll(absoluteDir, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				//create destination file making sure the path is writeable.
				constructFile(projectRootDir+files[i].Filename, file) //file path + filename, file content
				file.Close()
			}
			fmt.Println("Finished copying ", i+1, " files")

			//container path name
			fileDir := filepath.Dir(files[0].Filename)
			fileDir = strings.Replace(fileDir, "\\", "/", -1)
			parts := strings.Split(fileDir, "/")
			containerPath = projectRootDir + parts[0] //container path is also zip path
			passContainerPath = projectRootDir + parts[0]

			//zip the folder
			sourceDir = containerPath + "/"
			zipName = parts[0] + ".zip"
			destDir = projectRootDir + zipName

			containerName = strings.Replace(strings.ToLower(parts[0]), " ", "", -1)
			containerName = containerName + "_" + userId + "_" + strconv.Itoa(projectId)
			//containerName = strings.Replace(tempName, " ", "", -1)
			fmt.Println("containerName :" + containerName)
			//container name is (user project root folder name )+ count
			//tempName := strings.ToLower(parts[0]) + containerCount
			//containerName = strings.Replace(tempName, " ", "", -1)

			//display success message.
		} else if pj.ProjectType == "Window" {
			pj.Language = m.Value["language"][0]
			pj.LanguageVersion = m.Value["language_version"][0]
			pj.Database = m.Value["database"][0]
			pj.DatabaseVersion = m.Value["database_version"][0]
			pj.DatabaseName = m.Value["database_name"][0]

			f := m.File["application"]
			if len(f) != 0 {
				applicatoinFile, err := f[0].Open()
				if err != nil {
					log.Println(err)
					return
				}
				//application name as dir name (removing extension) and store under it's dir with the same name
				fileName := f[0].Filename
				parts := strings.Split(fileName, ".")
				containerPath = projectRootDir + parts[0]                      //for creating folder
				passContainerPath = projectRootDir + parts[0] + "/" + fileName //for creating folder

				containerName = strings.Replace(strings.ToLower(parts[0]), " ", "", -1)
				containerName = containerName + "_" + userId + "_" + strconv.Itoa(projectId)
				//containerName = strings.Replace(tempName, " ", "", -1)
				fmt.Println("containerName :" + containerName)
				//container name is (user project root folder name )+ count
				//tempName := strings.ToLower(parts[0]) + containerCount
				//containerName = strings.Replace(tempName, " ", "", -1)

				//container path is dir name since application name is dir name
				pathErr := os.MkdirAll(containerPath, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				fmt.Println("pathName :" + containerPath)
				//create destination file making sure the path is writeable.
				constructFile(containerPath+"/"+f[0].Filename, applicatoinFile) //file path + filename, file content
				applicatoinFile.Close()

				//zip the folder
				sourceDir = containerPath + "/"
				zipName = parts[0] + ".zip"
				destDir = projectRootDir + zipName

			}
		} else if pj.ProjectType == "Android" {
			f := m.File["apk"]
			if len(f) != 0 {
				pj.Device = m.Value["device"][0]
				pj.AndroidVersion = m.Value["android_version"][0]
				pj.APIVersion = m.Value["api_version"][0]
				applicatoinFile, err := f[0].Open()
				if err != nil {
					log.Println(err)
					return
				}

				//application name as dir name (removing extension) and store under it's dir with the same name
				fileName := f[0].Filename
				parts := strings.Split(fileName, ".")
				containerPath = projectRootDir + parts[0] //for creating folder
				passContainerPath = projectRootDir + parts[0] + "/" + fileName

				containerName = strings.Replace(strings.ToLower(parts[0]), " ", "", -1)
				containerName = containerName + "_" + userId + "_" + strconv.Itoa(projectId)
				//containerName = strings.Replace(tempName, " ", "", -1)
				fmt.Println("containerName :" + containerName)
				//container name is (user project root folder name )+ count
				//tempName := strings.ToLower(parts[0]) + containerCount
				//containerName = strings.Replace(tempName, " ", "", -1)

				//container path is dir name since application name is dir name
				pathErr := os.MkdirAll(containerPath, 0777)
				if pathErr != nil {
					log.Panic(pathErr)
				}
				fmt.Println("pathName :" + containerPath)
				//create destination file making sure the path is writeable.
				constructFile(containerPath+"/"+f[0].Filename, applicatoinFile) //file path + filename, file content
				applicatoinFile.Close()

				//zip the folder
				sourceDir = containerPath + "/"
				zipName = parts[0] + ".zip"
				destDir = projectRootDir + zipName
			}
		}
		//get backup file
		if pj.ProjectType != "Android" {
			f := m.File["database_backup"]
			fmt.Println("length:", len(f))
			if len(f) != 0 {

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
				fmt.Println("back work")
				if f[0].Filename != "" {
					constructFile(databasePath+"/"+f[0].Filename, backUpFile)
				}

			}
			//end of constructing backup

		}

		//creating zip
		zipit(sourceDir, destDir)
		passContainer := passContainerPath[1:len(passContainerPath)]
		fmt.Println("passContainerPath :", passContainer)
		err = createContainer.CreateContainer(pj.Database, databasePath, pj.DatabaseVersion, pj.Language, pj.LanguageVersion, passContainer, containerName, pj.DatabaseName, pj.Device, pj.AndroidVersion, pj.APIVersion)
		//if create container success , write container info , increate count
		if err != nil {
			log.Panic(err)
		}
		// data := containerName + ":" + contianerID + "\n"
		// err = filehandler.WriteContainerInfo(containerInfoPath, data)
		// if err != nil {
		// 	log.Panic(err)
		// }

		//get IP and save to database
		t := time.Now()
		ts := t.Format("2006-01-02 15:04:05")
		pj.CreateDate = ts
		pj.LastModifiedDate = ts
		if pj.ProjectType == "Website" {
			fmt.Println("Pj.Language :", pj.Language)
			fmt.Println("container :", containerName)
			containerIP = getIP(containerName, pj.Language)

			pj.DatabaseBackupPath = databasePath
			pj.ZipPath = zipName
			pj.Path = containerPath
			pj.ContainerIP = containerIP
			pj.ContainerName = containerName
			err := dbutil.SaveWebAndWindowProject(con, pj)
			if err != nil {
				log.Println(err)
				return
			}

		} else if pj.ProjectType == "Window" {
			fmt.Println("Pj.Language :", pj.Language)
			fmt.Println("container :", containerName)
			containerIP = getIP(containerName+"bridge", pj.Language)

			pj.DatabaseBackupPath = databasePath
			pj.ZipPath = zipName
			pj.Path = containerPath
			pj.ContainerIP = containerIP
			pj.ContainerName = containerName
			err := dbutil.SaveWebAndWindowProject(con, pj)
			if err != nil {
				log.Println(err)
				return
			}

		} else {

			fmt.Println("Pj.Language :", pj.Language)
			fmt.Println("container :", containerName)

			containerIP = getIP(containerName, pj.Language)
			//fmt.Println("start sleeping ...")
			//time.Sleep(60* time.Second)
			//fmt.Println("sleeping ...")

			pj.DatabaseBackupPath = databasePath
			pj.ZipPath = zipName
			pj.Path = containerPath
			pj.ContainerIP = containerIP
			pj.ContainerName = containerName
			err := dbutil.SaveAndroidProject(con, pj)
			if err != nil {
				log.Println(err)
				return
			}
		}

		fmt.Println("container IP :", containerIP)
		//close connection
		dbutil.CloseConnection(con)

		// err = appendProjectInfo(projectInfoPath, projectInfo)
		if err != nil {
			fmt.Println("upload failed.")
			displayUpload(w, "upload", "upload failed, try again.")
		} else {
			fmt.Println("upload successful")
			displayUpload(w, "upload", "Upload successful.")
		}
		fmt.Println("reach end")
		//w.WriteHeader(200)
		fmt.Fprintln(w, "ok")

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

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/project", projectHandler)
	r.HandleFunc("/startApplication", startApplicationHandler)
	r.HandleFunc("/download", downloadHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/profile ", profileHandler)
	r.HandleFunc("/logout", logoutHandler)
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
		w.WriteHeader(http.StatusOK)
		w.Write(byteValue)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]interface{}{}
	message["Username"] = ""
	message["UserID"] = ""
	message["Email"] = ""

	// check if session is active
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	if session != nil {
		message["Username"] = session.Values["username"]
		message["UserID"] = session.Values["userid"]
		message["Email"] = session.Values["email"]
	}
	switch r.Method {
	case "GET":
		displayHome(w, "index", message)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//read from session
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	// remove the username
	session.Values["username"] = ""
	session.Values["email"] = ""
	session.Values["userid"] = ""
	err := session.Save(r, w)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "http://localhost:8080/upload", http.StatusFound)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]interface{}{}
	message["Username"] = ""
	message["UserID"] = ""
	message["Email"] = ""

	// check if session is active
	session, _ := loggedUserSession.Get(r, "authenticated-user-session")

	if session != nil {
		message["Username"] = session.Values["username"]
		message["UserID"] = session.Values["userid"]
		message["Email"] = session.Values["email"]
	}
	switch r.Method {
	case "GET":
		displayRegister(w, "register", message)

	case "POST":
		fmt.Println("post register work")
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		fmt.Println(username, password, email)
		if username != "" && password != "" && email != "" {
			//save to database
			fmt.Println("registering to db")
			user := dbutil.User{}
			user.Username = username
			user.Password = password
			user.Email = email
			con, err := dbutil.GetConnection()
			if err != nil {
				log.Panic(err)
			}
			flag, err := dbutil.IsEmailExist(email, con)
			if err != nil {
				log.Panic(err)
			}
			if !flag {
				err = dbutil.CreateUser(con, user)
				if err != nil {
					log.Panic(err)
				}
				message["Success"] = "User Creation Success!"
			} else {
				err = dbutil.CloseConnection(con)
				if err != nil {
					log.Panic(err)
				}
				message["Fail"] = "Email Already Exist!"
			}

			displayRegister(w, "register", message)
		}
	}

}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]interface{}{}
	message["Username"] = ""
	message["Email"] = ""
	message["UserID"] = ""

	switch r.Method {
	case "GET":
		displayLogin(w, "login", message)
	case "POST":
		fmt.Println("login work")

		//check if session is active
		session, _ := loggedUserSession.Get(r, "authenticated-user-session")
		if session != nil {
			message["Username"] = session.Values["username"]
			message["Email"] = session.Values["email"]
			message["UserID"] = session.Values["userid"]

		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		//get connection from db
		con, err := dbutil.GetConnection()
		if err != nil {
			log.Println(err)
			return
		}

		if email != "" && password != "" {
			ans, err := dbutil.IsUserExist(email, password, con)
			if err != nil {
				log.Println(err)
			}
			if !ans {
				message["Fail"] = "Email or Password incorrect."
				message["LoginError"] = true
				displayLogin(w, "login", message)
			} else {
				userId, userName, err := dbutil.GetUserByEmail(email, con)
				if err != nil {
					log.Println(err)
					return
				}
				message["Username"] = userName
				message["UserID"] = userId
				message["Email"] = email
				// create a new session and redirect to dashboard
				session, _ := loggedUserSession.New(r, "authenticated-user-session")

				session.Values["username"] = userName
				session.Values["userid"] = userId
				session.Values["email"] = email
				err = session.Save(r, w)

				if err != nil {
					log.Println(err)
				}

				http.Redirect(w, r, "http://localhost:8080/upload   ", http.StatusFound)
			}
		}

	}

}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayProfile(w, "profile", "")
	}

}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	var zipPath = "./userproject/"
	var zipname string = ""
	var projectId = ""
	var userId = ""

	zipname = r.URL.Query().Get("zipName")
	userId = r.URL.Query().Get("userId")
	projectId = r.URL.Query().Get("projectId")
	if zipname != "" {
		data, err := ioutil.ReadFile(zipPath + userId + "/" + projectId + "/" + zipname)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+zipname)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		http.ServeContent(w, r, zipPath+zipname, time.Now(), bytes.NewReader(data))
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
