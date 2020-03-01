package main

import (
	//  "os"
	//  "io"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	imageName := "java:8"
	//  out,err:=cli.ImagePull(ctx,imageName,types.ImagePullOptions{})
	//  if err!=nil{
	//  panic(err)
	//  }
	//  io.Copy(os.Stdout,out)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   []string{"MODE=tcp", "XPRA_HTML=yes", "DISPLAY=:14"},
		Cmd:   []string{"java", "-jar", "UIOutput.jar"},
		Tty:   true,
	},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/home/htunlin/go/src/UIOutput.jar",
					Target: "/UIOutput.jar",
				},
			},
			VolumesFrom: []string{"koko_sfefrrtkjhjkbridge"},
		}, nil, "javatest14")
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println(resp.ID)
}
