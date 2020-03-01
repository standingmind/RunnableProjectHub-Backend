package main
import (
// "os"
// "io"
"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"fmt"
"github.com/docker/docker/api/types/mount"
"github.com/docker/docker/client"
"golang.org/x/net/context"
)
func main() {
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
panic(err)
}
imageName:="mysql"
// out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
// if err!=nil{
// panic(err)
// }
srcpath:="/home/htunlin/go/src/userproject/Online_exam_free_mysqli/database/quiz_new.sql/"
// io.Copy(os.Stdout,out)
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:imageName,
Env:[]string{"MYSQL_ROOT_PASSWORD=root","MYSQL_DATABASE=mysqldb"},
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
nil,"mysqlforonlineExamtest2")
if err!=nil{
panic(err)
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	panic(err)
}
fmt.Println(resp.ID)
}












