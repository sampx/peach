// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package models

import (
	"os"

	"github.com/Unknwon/com"
	"github.com/Unknwon/log"

	"github.com/sampx/peach/pkg/setting"
)

func initLangDocs(tocs map[string]*Toc, localRoot, lang string) {
	toc := tocs[lang]

	for _, dir := range toc.Nodes {
		if !com.IsFile(dir.FileName) {
			continue
		}

		//解析md文件为js
		if err := dir.ReloadContent(); err != nil {
			log.Error("Fail to load doc file: %v", err)
			continue
		}

		for _, file := range dir.Nodes {
			if !com.IsFile(file.FileName) {
				continue
			}

			if err := file.ReloadContent(); err != nil {
				log.Error("Fail to load doc file: %v", err)
				continue
			}
		}
	}

	//解析toc目录下的文件
	for _, page := range toc.Pages {
		if !com.IsFile(page.FileName) {
			continue
		}

		if err := page.ReloadContent(); err != nil {
			log.Error("Fail to load doc file: %v", err)
			continue
		}
	}
}

func initDocs(tocs map[string]*Toc, localRoot string) {
	for _, lang := range setting.Docs.Langs {
		initLangDocs(tocs, localRoot, lang)
	}
}

//NewContext todo doc
func NewContext() {
	//如果存在html目录，则先删除
	if com.IsExist(HTMLRoot) {
		if err := os.RemoveAll(HTMLRoot); err != nil {
			log.Fatal("Fail to clean up HTMLRoot: %v", err)
		}
	}

	//重新加载Docs
	if err := ReloadDocs(); err != nil {
		log.Fatal("Fail to init docs: %v", err)
	}
}
