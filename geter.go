package mulgeter

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const FENPIAN = 2048

type TestInterface interface {
	Read()
}
type Mulgeter struct {
	length     int
	tasknums   int
	overnums   int
	filename   string
	fileext    string
	requestUrl string
	TestInterface
}

func NewMulgeter(url string) *Mulgeter {
	urlslice := strings.Split(url, "/")
	fmt.Printf("%v", urlslice)
	filename := urlslice[len(urlslice)-1]
	fmt.Printf("名字%s", filename)
	return &Mulgeter{
		requestUrl: url,
		filename:   filename,
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
func (m *Mulgeter) Read() {
	fmt.Println("读取")
}
func (m *Mulgeter) CalcTask() {
	c := float64(m.length) / FENPIAN
	tx := math.Ceil(c)
	fmt.Print(tx)
	m.tasknums = int(int64(tx))
	fmt.Println("任务数", m.tasknums)
	return
}
func (m *Mulgeter) BeginDownload() {
	m.Read()
	m.CalcTask()
	var wg sync.WaitGroup
	startt := 0
	endt := 0
	lett := m.length
	for i := 0; i < m.tasknums; i++ {
		j := i
		fmt.Println("任务", j)
		if m.length <= FENPIAN {
			startt = 0
			endt = m.length - 1
		} else {
			if j == 0 {
				startt = 0
				endt = FENPIAN - 1
			} else {
				startt = endt + 1
				if lett <= FENPIAN {
					endt = startt + lett - 1
				} else {
					endt = startt + FENPIAN - 1
				}
			}
			lett = lett - FENPIAN
			fmt.Println("lett", lett)
			fmt.Println("lett", startt)
		}
		wg.Add(1)
		go func() {
			str := fmt.Sprintf("bytes=%d-%d", startt, endt)
			fmt.Println(str)
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
			f, err := os.OpenFile(fmt.Sprintf("%d", j), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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
	zongfile, _ := os.OpenFile(m.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	for i := 0; i < m.tasknums; i++ {
		f, _ := os.OpenFile(fmt.Sprintf("%d", i), os.O_RDWR, 0666)
		wrted, _ := io.Copy(zongfile, f)
		fmt.Println("写入多少字节", wrted)
	}
	return
}
