package main

import (
	"bytes"
	"fmt"
	"github.com/truexf/gobuffer"
	"github.com/truexf/goutil"
	"io/ioutil"
	"os"
	"time"
)

func test1() {
	fn := fmt.Sprintf("/tmp/gobuffer%d.buf", time.Now().Unix())
	fmt.Println(fn)
	fd, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	buf := new(bytes.Buffer)
	for i := 0; i < 10240; i++ {
		buf.WriteString(time.Now().String())
		buf.WriteString("\n")
	}
	bts := buf.Bytes()
	md5, _ := goutil.BytesMd5(bts)
	fmt.Println(len(bts))
	fmt.Println("md5:", md5)
	goBuf, _ := gobuffer.NewGoBuffer(256, fd, 1)
	goBuf.Start()
	pos := 0
	for {
		x := len(bts) - pos
		if x > 160 {
			x = 160
		}
		n, e := goBuf.Write(bts[pos : pos+x])
		pos += n
		// fmt.Printf("len(bts) %d,pos %d\n", len(bts),pos)
		if pos >= len(bts) {
			break
		}
		if e != nil {
			fmt.Println(e.Error())
			return
		}
	}
	goBuf.Flush()
	<-time.After(time.Second * 2)
	fd.Close()
	fmd5, _ := goutil.FileMd5(fn)
	fmt.Printf("%s:%s\n", md5, fmd5)
	fmt.Printf("size: %d\n", len(bts))
	fbts, _ := ioutil.ReadFile(fn)
	s, _ := goutil.BytesMd5(fbts)
	fmt.Println(s)
}

func test2() {
	fn := fmt.Sprintf("/tmp/gobuffer%d.buf", time.Now().Unix())
	fd, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0666)
	buf := new(bytes.Buffer)
	sLen := 0
	for i := 0; i < 10240; i++ {
		s := fmt.Sprintf("YYYY-MM-DD hh:NN:ss:ZZZ hello worldYYYY-MM-DD hh:NN:ss:ZZZ hello world %d\n", i)
		buf.WriteString(s)
		sLen = len(s)
	}
	bts := buf.Bytes()
	md5, _ := goutil.BytesMd5(bts)
	fmt.Println(len(bts))
	fmt.Println(md5)
	goBuf, _ := gobuffer.NewGoBuffer(256, fd, 1)
	goBuf.Start()
	count := 10
	ch := make(chan int, count)
	for i := 0; i < count; i++ {
		fmt.Printf("i: %d\n", i)
		go func() {
			pos := 0
			for {
				x := len(bts) - pos
				if x > sLen {
					x = sLen
				}
				n, e := goBuf.Write(bts[pos : pos+x])
				pos += n
				// fmt.Printf("pos %d\n", pos)
				// fmt.Printf("len(bts) %d,pos %d\n", len(bts),pos)
				if pos >= len(bts) {
					fmt.Printf("pos 1 %d\n", pos)
					break
				}
				if e != nil {
					fmt.Println(e.Error())
					return
				}

			}
			fmt.Printf("posend %d\n", pos)
			goBuf.Flush()
			<-time.After(time.Second * 2)
			ch <- 1
		}()
	}
	for i := 0; i < count; i++ {
		<-ch
	}
	<-time.After(time.Second * 2)
	fd.Close()
	fdata, _ := ioutil.ReadFile(fn)
	fmt.Printf("file size: %d\n", len(fdata))
}

func main() {
	test1()
	//go test2()

	// for {
	// 	<-time.After(time.Second)
	// 	fmt.Println("fuck")
	// }

	return
}
