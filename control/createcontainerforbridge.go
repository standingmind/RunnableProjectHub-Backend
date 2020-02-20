package control
import (
 "os"
 "io"
"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"github.com/docker/docker/client"
"golang.org/x/net/context"
"fmt"
)
func CreateContainerforBridge(CName string) error{
ctx :=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil{
return err
}
imageName:="htetaungkhant/ownui"
 out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
 if err!=nil{
	return err
 }
 io.Copy(os.Stdout,out)
resp,err:=cli.ContainerCreate(ctx,&container.Config{
Image:imageName,
Env:[]string{"MODE=tcp","XPRA_HTML=yes","DISPLAY=:14"},
Tty:true,
},
nil,nil,CName)
if err!=nil{
	return err
}
if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
	return err
}
fmt.Println("started container")
return nil
}












