package main
import (
	"os/exec"
	"fmt"
	"bytes"
)
func main(){
	cmd :=exec.Command("docker","exec","-i","gg_11_48","adb","install","/gg.apk")
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