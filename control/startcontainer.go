package control
import(
"fmt"
"github.com/docker/docker/api/types"
"golang.org/x/net/context"
"github.com/docker/docker/client"
)
func StartContainer(CName string){
	fmt.Println("CName="+CName)
	ctx:=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
	   panic(err)
	}
	if err:=cli.ContainerStart(ctx,CName,types.ContainerStartOptions{});err!=nil{
		panic(err)
	}
	   
	fmt.Println("started container")
}

