package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type DownLoadTask struct {
	URL      string
	Savepath string
}

func DownLoader(taskchan chan DownLoadTask, wg *sync.WaitGroup, resultchan chan string) {
	defer wg.Done()
	//请求网址，得到响应体
	for task := range taskchan {
		resp, err := http.Get(task.URL)
		if err != nil {
			resultchan <- fmt.Sprintf("下载 %s 失败！错误：%v", task.URL, err)
			continue
		}
		//关闭响应体
		defer resp.Body.Close()
		//检查状态码
		if resp.StatusCode != http.StatusOK {
			resultchan <- fmt.Sprintf("%s 响应失败！错误：%d", task.URL, resp.StatusCode)
			continue
		}
		//创建文件
		file, err := os.Create(task.Savepath)
		if err != nil {
			resultchan <- fmt.Sprintf("创建 %s 文件失败！错误：%v", task.URL, err)
			continue
		}
		//关闭文件
		defer file.Close()
		//写入文件
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			resultchan <- fmt.Sprintf("写入 %s 文件失败！错误：%v", task.URL, err)
			continue
		}
		resultchan <- fmt.Sprintf("%s下载操作成功~,储存在——%s", task.URL, task.Savepath)
	}
}

func main() {
	//要爬取的网址以及要储存的地址
	tasks := []DownLoadTask{
		{URL: "https://mis.bjtu.edu.cn/home/", Savepath: "D:\\Go爬虫\\爬取结果\\bjtu.html"},
		{URL: "https://www.liwenzhou.com/", Savepath: "D:\\Go爬虫\\爬取结果\\teacher.html"},
		{URL: "https://pkg.go.dev/", Savepath: "D:\\Go爬虫\\爬取结果\\go.html"},
	}

	var wg sync.WaitGroup
	//初始化通道
	taskchan := make(chan DownLoadTask, len(tasks))
	resultchan := make(chan string, len(tasks))
	Nums := len(tasks)
	//开启并发
	for i := 0; i < Nums; i++ {
		wg.Add(1)
		go DownLoader(taskchan, &wg, resultchan)
	}
	//把任务分发下去
	for _, task := range tasks {
		taskchan <- task
	}
	close(taskchan)
	//程序等待，全部完成后关闭接收通道
	go func() {
		wg.Wait()
		close(resultchan)
	}()
	//打印结果
	fmt.Println("最终结果是：")
	for result := range resultchan {
		fmt.Println(result)
	}
	fmt.Println("——————————————————所有任务已完成！————————————————")
}
