package main
import (
	"fmt"
	"os"
)
func main(){
	f,err:=os.Create("/home/htunlin/go/src/userproject/Online_exam_free_mysqli/dbstuff.inc")
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return
	}
	ipaddr:="172.17.0.2"
	d:=[]string{"<?php","$host="+"\""+ipaddr+"\";","$username="+"\""+"root"+"\";","$password="+"\"root\";","$dbname="+"\"mysqldb\";","?>"}
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