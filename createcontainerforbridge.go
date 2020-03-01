package main
import (
 "os"
 "io"
"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"fmt"
"github.com/docker/docker/client"
"golang.org/x/net/context"
)
func main() {
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
panic(err)
}
imageName:="jare/x11-bridge"
 out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
 if err!=nil{
 panic(err)
 }
 io.Copy(os.Stdout,out)
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:imageName,
Env:[]string{"MODE=tcp","XPRA_HTML=yes","DISPLAY=:14"},
Tty:true,
},
nil,nil,"x11-bridged4")
if err!=nil{
panic(err)
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	panic(err)
}
fmt.Println(resp.ID)
}












