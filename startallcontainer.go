package main
import(
	"bufio"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"github.com/docker/docker/client"
	"log"
	"os"
)
func main(){
	file, err := os.Open("/home/htunlin/go/src/result.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	ctx:=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err:=cli.ContainerStart(ctx,scanner.Text(),types.ContainerStartOptions{});err!=nil{
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}


}

