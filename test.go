package main

import (
	/*"io"
	"os"*/
	"github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"fmt"
	"io"
	"os"
)

func main() {
	ctx :=context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
	if err!=nil{
	panic (err)
}
	srcpath:="/home/htunlin/go/src/gg.apk"
	fmt.Println(srcpath)
	imageName:="budtmo/docker-android-x86-8.1"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	resp,err:=cli.ContainerCreate(ctx,&container.Config{
	Image:imageName,
	Env:[]string{"DEVICE=Samsung Galaxy S6"},
	Tty:true,
	},
	&container.HostConfig{ 
		Privileged:true,
		Mounts:[]mount.Mount {
			{
				Type: mount.TypeBind,
				Source:srcpath,
				Target:"/gg.apk",
			},
		},
	},nil,"Androidtest")
	if err!=nil{
		panic(err)
	}
	if err:=cli.ContainerStart(ctx,resp.ID,types.ContainerStartOptions{});err!=nil{
		panic(err)
	}

}
