package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shpdwx/web/conf"
	"github.com/shpdwx/web/glm"
	"github.com/shpdwx/web/mq"
	"github.com/shpdwx/web/s3"
)

var (
	wg sync.WaitGroup

	stopChan = make(chan struct{}, 1)
	logsChan = make(chan string, 30)
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// yaml conf
	conf, err := conf.NewConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	logId := uuid.New()
	fmt.Printf("request id: %s\n", logId.String())

	// 错误处理
	wg.Add(1)
	go func() {
		defer wg.Done()

		// startup mq
		ch := mq.NewRabbitMQ(conf.RabbitMQ)
		// defer ch.Close()

		for {
			select {
			case v := <-logsChan:
				mq.LogMq(ctx, ch, conf.RabbitMQ, v)

			case <-stopChan:
				fmt.Println("stop chan.")
				return
			}
		}
	}()

	generate(ctx, logId, conf)

	stopChan <- struct{}{}
	wg.Wait()
}

func generate(ctx context.Context, id uuid.UUID, conf conf.Conf) {

	txt := "古风水墨画江山如此多娇"

	rsp, err := glm.Image(conf.CogView, txt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successly generate image.")
	record(id, rsp)

	for _, v := range rsp.Data {

		fmt.Printf("\n处理图片 %s\n\n", v.Url)

		image, err := s3.Fetch(v.Url)
		if err != nil {
			fmt.Printf("下载图片失败 %v", err)
			return
		}
		fmt.Println("Successly download image.")
		image.RequestId = id.String()
		image.Desc = txt
		record(id, image)

		info, err := s3.Put(ctx, conf.Minio, image)
		if err != nil {
			fmt.Printf("上传图片失败 %v", err)
			return
		}
		fmt.Println("Successly put image.")
		record(id, info)

	}
}

func record(id uuid.UUID, rsp any) {
	if b, err := json.Marshal(rsp); err == nil {
		logsChan <- fmt.Sprintf("uuid %s content %s", id.String(), string(b))
	}
}
