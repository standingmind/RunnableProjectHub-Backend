package control
import (
	"os/exec"
	"bytes"
    "fmt"
)
func InstallApk(CName string){
	cmd :=exec.Command("docker","exec","-i",CName,"adb","install","/gg.apk")
	cmdOutput:=&bytes.Buffer{}
	cmdError:=&bytes.Buffer{}
	cmd.Stdout=cmdOutput
	cmd.Stderr=cmdError
	err:=cmd.Run()
	if err!=nil{
		fmt.Println(cmdError.String())
	
	}
	fmt.Println("Result:"+cmdOutput.String())
}