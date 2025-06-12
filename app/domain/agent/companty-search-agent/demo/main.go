/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}

	url := "https://www.tianyancha.com/nsearch?key=%E9%95%BF%E6%98%A5%E5%8C%97%E7%94%B5%E8%BD%AF%E4%BB%B6"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
	but.Cleanup()
}
