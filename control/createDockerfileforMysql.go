package control
import (
	"fmt"
	"log"
	"os"
)
func CreateDockerfileforMysql(dbversion string)error{
	dir,err:=os.Getwd()
	if err!=nil{
		log.Fatal(err)
	}
	f,err:=os.Create(dir+"/control/Dockerfile")
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