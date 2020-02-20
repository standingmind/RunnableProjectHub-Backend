package control
import (
"github.com/docker/docker/client"
"golang.org/x/net/context"
)
func StopContainer(CName string) error {

ctx:=context.Background()
cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
if err!=nil {
panic(err)
}
if err:=cli.ContainerStop(ctx,CName,nil);err!=nil{
	return err
}
return nil
}
