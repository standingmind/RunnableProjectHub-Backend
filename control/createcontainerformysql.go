package control
import (
// "os"
// "io"
"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"github.com/docker/docker/api/types/mount"
"github.com/docker/docker/client"
"golang.org/x/net/context"
)
func CreateContainerforMysql(sqlpath string,iname string,dbname string) error {
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
return err
}
imageName:=iname
// out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
// if err!=nil{
// panic(err)
// }
srcpath:="/home/htunlin/go/src/"+sqlpath+"/"
// io.Copy(os.Stdout,out)
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:imageName,
Env:[]string{"MYSQL_ROOT_PASSWORD=root","MYSQL_DATABASE="+dbname},
Cmd:[]string{"--default-authentication-plugin=mysql_native_password"},
Tty:true,
},
&container.HostConfig{ 
	Mounts:[]mount.Mount {
		{
			Type: mount.TypeBind,
			Source:srcpath,
			Target:"/docker-entrypoint-initdb.d/",
		},
	},
},
nil,iname)
if err!=nil{
return err
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	return err
}
return nil
}