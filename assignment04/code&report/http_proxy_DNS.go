/*****************************************************************************
 * http_proxy_DNS.go                                                                 
 * Names: 
 * NetIds:
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)
package main

import(
    "bytes"
    "fmt"
    "io"
    "log"
    "net"
    "net/url"
    "os"
    "strings"
    
    "golang.org/x/net/html"
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
    
    
    htmlbytes:=make([]byte,0)
    buf:=bytes.NewBuffer(htmlbytes)
    _,err=io.Copy(buf,server)
    
    go dnsPrefetch(buf.String())
    
    _,err=io.Copy(client,buf)
    if err != nil {
        return
    }
    
}                        
                        
func dnsPrefetch(htmlResp string){
    doc, err := html.Parse(strings.NewReader(htmlResp))
    if err != nil {
        log.Fatal(err)
    }
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" {
                    hostURL,err:=url.Parse(a.Val)
                    _,err=net.LookupHost(hostURL.Host)
                        break
                    }
                }
            }
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                f(c)
            }
        }
        f(doc)
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