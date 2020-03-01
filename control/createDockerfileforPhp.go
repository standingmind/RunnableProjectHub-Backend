package control
import (
	"fmt"
	"log"
	"os"
)
func CreateDockerfileforPhp(imageName string)error{
	dir,err:=os.Getwd()
	if err!=nil{
		log.Fatal(err)
	}
	f,err:=os.Create(dir+"/control/Dockerfile")
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