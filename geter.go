package mulgeter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Mulgeter struct {
	length     int
	tasknums   int
	overnums   int
	requestUrl string
}

func NewMulgeter(url string) *Mulgeter {
	return &Mulgeter{
		requestUrl: url,
	}
}

func (m *Mulgeter) GetLength() int {
	headres, _ := http.Head(m.requestUrl)
	longi, _ := strconv.Atoi(headres.Header.Get("Content-length"))
	m.length = longi
	return longi
}
func (m *Mulgeter) Length() int {
	return m.length
}
func (m *Mulgeter) BeginDownload() {
	var wg sync.WaitGroup
	for i := 0; i < m.length; i++ {
		j := i
		wg.Add(1)
		go func() {
			str := fmt.Sprintf("bytes=%d-%d", j, j)
			reqObj, err := http.NewRequest("GET", m.requestUrl, nil)
			if err != nil {
				fmt.Println("错误")
			}
			reqObj.Header.Set("Range", str)
			client := &http.Client{
				Timeout: time.Second * 5,
			}
			resp, err := client.Do(reqObj)
			if err != nil {
				return
			}
			defer func() {
				resp.Body.Close()
				wg.Done()
			}()
			// b, err := ReadAll(repon.Body)
			// 有多少
			f, err := os.OpenFile(fmt.Sprintf("%d.txt", j), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			b, err := io.Copy(f, resp.Body)
			if err != nil {
				return
			}
			fmt.Println("复制长度", b)
		}()
	}
	wg.Wait()
	m.mergeFile()
	return
}
func (m *Mulgeter) mergeFile() {
	fmt.Println("下载完成")
	zongfile, _ := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	for i := 0; i < m.length; i++ {
		f, _ := os.OpenFile(fmt.Sprintf("%d.txt", i), os.O_RDWR, 0666)
		wrted, _ := io.Copy(zongfile, f)
		fmt.Println("写入多少字节", wrted)
	}
	return
}
