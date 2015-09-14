package main

import (
        _ "bufio"
        "bytes"
        "fmt"
        "github.com/gorilla/websocket"
        "io/ioutil"
        "net"
        "net/http"
        _ "net/http/pprof"
        "net/url"
        "os"
        "sync"
        "time"
)

var C = 1000

//var sout = bufio.NewWriter(os.Stdout)
var m = &sync.Mutex{}

func main() {
        wg := &sync.WaitGroup{}

        for i := 0; i < C; i++ {
                wg.Add(1)
                go connect(i, wg)
        }
        wg.Wait()
}

func connect(cid int, wg *sync.WaitGroup) {
        conn, err := net.Dial("tcp", "localhost:4000")
        if err != nil {
                fmt.Println(err)
                return
        }
        u, err := url.Parse("http://localhost:4000/foobar")
        if err != nil {
                fmt.Println(err)
                return
        }
        h := &http.Header{}
        h.Add("Origin", "localhost:4000")
        sc, _, err := websocket.NewClient(conn, u, *h, 0, 0)
        defer sc.Close()
        if err != nil {
                fmt.Println(err)
                return
        }
        var mid = 0
        for {
                time.Sleep(1 * time.Nanosecond)
                if err = sc.WriteMessage(1, []byte("Hello from client")); err != nil {
                        fmt.Println(err)
                        return
                }
                _, p, err := sc.ReadMessage()
                if err != nil {
                        fmt.Println(err)
                        return
                }
                //m.Lock()
                sout := bytes.NewBuffer(make([]byte, 0))
                fmt.Fprintf(sout, "[c-%v][m-%v] RCVD:", cid, mid)
                fmt.Fprintf(sout, "%v\n", string(p))
                //m.Unlock()
                go ioutil.WriteFile(fmt.Sprintf("results/%v.%v", cid, mid), sout.Bytes(), os.ModePerm)
                mid++
                if mid == 100 {
                        return
                }
        }
        time.Sleep(5 * time.Second)
        // sout.Flush()
        time.Sleep(5 * time.Second)
        wg.Done()
}
