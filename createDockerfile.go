package main
import (
	"fmt"
	"os"
)
func main(){
	f,err:=os.Create("Dockerfile")
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return
	}
	d:=[]string{"FROM php:7.3-apache","RUN docker-php-ext-install mysqli && docker-php-ext-enable mysqli"}
	for _,v:=range d{
		fmt.Fprintln(f,v)
		if err!=nil{
			fmt.Println(err)
			return
		}
	}
	err=f.Close()
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("File written successfully")

}