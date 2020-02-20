package control
import (
	"fmt"
	"os"
)
func CreateDockerfileforPhp(imageName string)error{
	f,err:=os.Create("/home/htunlin/go/src/control/Dockerfile")
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return err
	}
	d:=[]string{"FROM "+imageName,"RUN docker-php-ext-install mysqli && docker-php-ext-enable mysqli","RUN docker-php-ext-install pdo_mysql"}
	for _,v:=range d{
		fmt.Fprintln(f,v)
		if err!=nil{
			fmt.Println(err)
			return err
		}
	}
	err=f.Close()
	if err!=nil{
		fmt.Println(err)
		return err
	}
	fmt.Println("File written successfully")
	return nil

}