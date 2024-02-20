/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: 
 * NetIds:
 *****************************************************************************/

 // TODO: implement an HTTP proxy

package main

import(
    "fmt"
    "io"
    "net"
    "net/url"
    "os"
    "strings"
)

func handleClientRequest(client net.Conn){
    if client==nil{
        return
    }
    defer client.Close()
    
    buffer := make([]byte, 1024)
    _, err := client.Read(buffer)
    if err != nil {
        return
    }
    
    reqs:=strings.Split(string(buffer)," ")
    if len(reqs)<=2{
        return
    }
    host:=reqs[1]
    hostURL,err:=url.Parse(host)
    if err != nil {
        return
    }
    
    var address string
    if hostURL.Opaque=="443"{
        address=hostURL.Scheme+":443"
    }else{
        if strings.Index(hostURL.Host,":")==-1{
            address=hostURL.Host+":80"
        }else{
            address=hostURL.Host
        }
    }
    
    server, err := net.Dial("tcp", address)
    if err != nil {
        return
    }
    
    req:=fmt.Sprintf("GET / HTTP/1.1\r\nHost:%s\r\nConnection: close\r\n\r\n",hostURL.Host)
    _, err = server.Write([]byte(req))
    if err != nil {
        return
    }
    
    _,err=io.Copy(client,server)
    if err != nil {
        return
    }
}                        
                        
                        
func main(){
    
    if len(os.Args)!=2{
        return
    }
    
    port:=os.Args[1]
    
    ln, err := net.Listen("tcp",fmt.Sprintf(":%s",port))
    if err != nil {
        return
    }
    
    defer ln.Close()
    
    for {
        conn, err := ln.Accept()
        if err != nil {
            return
        }
        go handleClientRequest(conn)
    } 
}
    
    
    