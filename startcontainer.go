package main
import(
"github.com/docker/docker/api/types"
"golang.org/x/net/context"
"github.com/docker/docker/client"
)
func main(){
	ctx:=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
		panic(err)
	}
	if err:=cli.ContainerStart(ctx,"javatest6",types.ContainerStartOptions{});err!=nil{
		panic(err)
	}
}

