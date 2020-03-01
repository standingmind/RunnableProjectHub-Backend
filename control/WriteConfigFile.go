package control
import (
	"fmt"
	"os"
	"strings"
)
func WriteConfigFile(path string,ipaddr string,dbname string) error{
	f,err:=os.Create("/home/htunlin/go/src"+path+"/dbstuff.inc")
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return err
	}
	ip:=strings.TrimSpace(ipaddr)
	d:=[]string{"<?php","$host="+"\""+ip+"\";","$user="+"\""+"root"+"\";","$password="+"\"root\";","$dbname="+"\""+dbname+"\";","?>"}
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