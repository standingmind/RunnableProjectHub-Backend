package control
import(
"fmt"
"github.com/docker/docker/api/types"
"golang.org/x/net/context"
"github.com/docker/docker/client"
)
func RemoveImage(CName string){
	fmt.Println("CName="+CName)
	ctx:=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
	   panic(err)
	}
	if _,err:=cli.ImageRemove(ctx, CName, types.ImageRemoveOptions{});err!=nil{
		panic(err)
	}
	   
	fmt.Println("Removed image")
}

