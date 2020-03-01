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
func main(){
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
	panic(err)
}
srcpath:="/home/htunlin/go/src/userproject/exam/"
fmt.Println(srcpath)
imageName:="onlineexamtest"
// out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
// if err!=nil{
// 	panic(err)
// }
// io.Copy(os.Stdout,out)
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:imageName,
},
&container.HostConfig{ 
	Mounts:[]mount.Mount {
		{
			Type: mount.TypeBind,
			Source:srcpath,
			Target:"/var/www/html",
		},
	},
},
nil,"onlineExamtest25")
if err!=nil{
	panic(err)
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	panic(err)
}
}












