package main
import(
"fmt"
"github.com/docker/docker/api/types"
"golang.org/x/net/context"
"github.com/docker/docker/client"
)
func main(){
	fmt.Println("CName=")
	ctx:=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
	   panic(err)
	}
	if _,err:=cli.ImageRemove(ctx, "jess/gimp", types.ImageRemoveOptions{});err!=nil{
		panic(err)
	}
	   
	fmt.Println("Removed image")
}

