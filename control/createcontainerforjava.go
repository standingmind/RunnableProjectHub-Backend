package control

import (
	//  "os"
	//  "io"
	// "github.com/docker/docker/api/types"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func CreateContainerforJava(jarpath string, CName string, bridgename string, checkdb bool) error {
	fmt.Println("Createing java")
	fmt.Println(bridgename)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
		return err
	}
	imageName := "java:8"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
	fmt.Println("path is ", jarpath)
	fmt.Println(bridgename)
	if checkdb == true {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Env:   []string{"MODE=tcp", "XPRA_HTML=yes", "DISPLAY=:14"},
			Cmd:   []string{"java", "-jar", "UIOutput.jar"},
			Tty:   true,
		},
			&container.HostConfig{
				Links: []string{CName + "sql:localhost"},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: jarpath,
						Target: "/UIOutput.jar",
					},
				},
				VolumesFrom: []string{bridgename},
			}, nil, CName)
		if err != nil {
			return err
		}
		fmt.Println(resp.ID)
	} else {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Env:   []string{"MODE=tcp", "XPRA_HTML=yes", "DISPLAY=:14"},
			Cmd:   []string{"java", "-jar", "UIOutput.jar"},
			Tty:   true,
		},
			&container.HostConfig{
				//Links: []string{CName + "sql:localhost"},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: jarpath,
						Target: "/UIOutput.jar",
					},
				},
				VolumesFrom: []string{bridgename},
			}, nil, CName)
		if err != nil {
			return err
		}
		fmt.Println(resp.ID)

	}
	return nil
}
