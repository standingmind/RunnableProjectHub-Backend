package control
import (
	"os/exec"
	"os"
	"bytes"
    "strings"
    "fmt"
)
func IPtest(containerName string,language string) string{
	fmt.Println("containerName :",containerName)
	fmt.Println("language :",language)
	cmd :=exec.Command("docker","inspect","--format","{{.NetworkSettings.IPAddress}}",containerName)
    cmdOutput:=&bytes.Buffer{}
	cmd.Stdout=cmdOutput
	err:=cmd.Run()
	if err!=nil{
		os.Stderr.WriteString(err.Error())
	}
	if strings.ToLower(language)=="j2ee" || strings.ToLower(language)=="html"{
	return string(cmdOutput.Bytes())+":8080/"+containerName
	}else if strings.ToLower(language)=="java"{
	 return string(cmdOutput.Bytes())+":10000"
	}else if strings.ToLower(language)=="php"{
		return string(cmdOutput.Bytes())
	}else{
		return string(cmdOutput.Bytes())+":6080"
	}
}

