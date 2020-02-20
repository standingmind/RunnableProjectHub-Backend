package control
import (
	"fmt"
	"os"
)
func CreateDockerfileforMysql(dbversion string)error{
	f,err:=os.Create("/home/htunlin/go/src/control/Dockerfile")
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return nil
	}
	d:=[]string{"FROM mysql:"+dbversion}
	for _,v:=range d{
		fmt.Fprintln(f,v)
		if err!=nil{
			fmt.Println(err)
			return nil
		}
	}
	err=f.Close()
	if err!=nil{
		fmt.Println(err)
		return nil
	}
	return nil

}