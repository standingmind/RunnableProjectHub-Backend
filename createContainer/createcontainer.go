package createContainer
import (
//"os"
//"io"
"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"fmt"
"github.com/docker/docker/api/types/mount"
"github.com/docker/docker/client"
"golang.org/x/net/context"
"strings"
"control"
)
func CreateContainer(db string,sqlpath string,dbversion string,language string,version string,path string,CName string,dbname string,Device string,Androidversion string,APIversion string) error{
fmt.Println(language)
fmt.Println(path)
fmt.Println(CName)
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
	return err
}
var ipaddr string
srcpath:="/home/htunlin/go/src"+path
fmt.Println(srcpath)
if (strings.ToLower(db)=="mysql"){
	err:=control.CreateDockerfileforMysql(dbversion)
	if err!=nil{
		panic(err)
	}
	iname:=CName+"sql"
	err=control.CreateImageforMysql(iname)
	if err!=nil{
		panic(err)
	}
	err=control.CreateContainerforMysql(sqlpath,iname,dbname)
	if err!=nil{
		panic(err)
	}
	ipaddr=control.IPtest(iname,language)
	fmt.Println(ipaddr)
}
if (strings.ToLower(language)=="php") {
fmt.Println(CName)
imageName:=strings.ToLower(language)+":"+version+"-apache"
err:=control.CreateDockerfileforPhp(imageName) 
if err!=nil{
	panic(err)
}
err=control.CreateImageforMysql(CName)
if err!=nil{
	panic(err)
}
 //err=control.WriteConfigFile(path,ipaddr,dbname)
 //if err!=nil{
 //	panic(err)
 //}
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:CName,
},
&container.HostConfig{ 
	Links:[]string{CName+"sql:localhost"},
	Mounts:[]mount.Mount {
		{
			Type: mount.TypeBind,
			Source:srcpath,
			Target:"/var/www/html",
		},
	},
},
nil,CName)
if err!=nil{
	return err
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	return err
}
}else if (strings.ToLower(language)=="j2ee")||(strings.ToLower(language)=="html"){
	imageName:="tomcat:9.0.1-jre8-alpine"
	trpath:="/usr/local/tomcat/webapps/"+CName
	/*out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
	if err!=nil{
		return err
	}
	io.Copy(os.Stdout,out)*/
	resp,err:=cli.ContainerCreate(ctx,&container.Config{
	Image:imageName,
	Cmd:[]string{"catalina.sh", "run"},
	},
	&container.HostConfig{ 
		Links:[]string{CName+"sql:localhost"},
		Mounts:[]mount.Mount {
			{
				Type: mount.TypeBind,
				Source:srcpath,
				Target:trpath,
			},
		},
	},
	nil,CName)
	if err!=nil{
		return err
	}
	if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
		return err
	}

}else if (strings.ToLower(language)=="java"){
	fmt.Println("Java")
	bname:=CName+"bridge"
	err:=control.CreateContainerforBridge(bname)
	if err!=nil{
		panic(err)
	}
	err=control.CreateContainerforJava(srcpath,CName,bname)
	if err!=nil{
		panic(err)
	}
	
}else {
	//srcpath:="/home/htetaungkhant/htetaungkhant/Projects/Proj1/src"+"/gg.apk"
	fmt.Println("Android")
	imageName:="budtmo/docker-android-x86-8.1"
	resp,err:=cli.ContainerCreate(ctx,&container.Config{
	Image:imageName,
	Env:[]string{"DEVICE=Samsung Galaxy S6"},
	},
	&container.HostConfig{ 
		Privileged:true,
		Mounts:[]mount.Mount {
			{
				Type: mount.TypeBind,
				Source:srcpath,
				Target:"/gg.apk",
			},
		},
	},
	nil,CName)
	if err!=nil{
		return err
	}
	if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
		return err
	}
	}
fmt.Println(Androidversion,APIversion)
return nil
}