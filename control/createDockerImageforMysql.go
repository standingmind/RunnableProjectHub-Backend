package control

import (
    "archive/tar"
    "bytes"
    "context"
    "io"
    "io/ioutil"
    "log"
    "os"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
)

func CreateImageforMysql(iname string)error {
    dir,err:=os.Getwd()
    if err!=nil{
        log.Fatal(err)
    }
    ctx := context.Background()
	cli,err:=client.NewClientWithOpts(client.FromEnv,client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }

    buf := new(bytes.Buffer)
    tw := tar.NewWriter(buf)
    defer tw.Close()

    dockerFile := "myDockerfile"
    dockerFileReader, err := os.Open(dir+"/control/Dockerfile")
    if err != nil {
        panic(err)
    }
    readDockerFile, err := ioutil.ReadAll(dockerFileReader)
    if err != nil {
        panic(err)
    }

    tarHeader := &tar.Header{
        Name: dockerFile,
        Size: int64(len(readDockerFile)),
    }
    err = tw.WriteHeader(tarHeader)
    if err != nil {
        panic(err)
    }
    _, err = tw.Write(readDockerFile)
    if err != nil {
        panic(err)
    }
    dockerFileTarReader := bytes.NewReader(buf.Bytes())

    imageBuildResponse, err := cli.ImageBuild(
        ctx,
        dockerFileTarReader,
        types.ImageBuildOptions{
            Context:    dockerFileTarReader,
            Tags:       []string{iname},
            Dockerfile: dockerFile,
            Remove:     true})
    if err != nil {
        panic(err)
    }
    defer imageBuildResponse.Body.Close()
    _, err = io.Copy(os.Stdout, imageBuildResponse.Body)
    if err != nil {
        panic(err)
	}
	return nil
}