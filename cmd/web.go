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

package cmd

import (
	"fmt"
	"net/http"

	"github.com/Unknwon/log"

	"github.com/go-macaron/i18n"
	"github.com/go-macaron/pongo2"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"

	"github.com/sampx/peach/models"
	"github.com/sampx/peach/pkg/context"
	"github.com/sampx/peach/pkg/setting"
	"github.com/sampx/peach/routes"
)

// Web docs todo
var Web = cli.Command{
	Name:   "web",
	Usage:  "Start Peach web server",
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "custom/app.ini", "Custom configuration file path"),
	},
}

func runWeb(ctx *cli.Context) {

	//AppVer是从peach.go中设置的
	log.Info("PeachDoc %s", setting.AppVer)

	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
		log.Info("加载自定义配置文件：%s", setting.CustomConf)
	}

	//加载配置文件
	setting.NewContext()
	//加载模型文件
	models.NewContext()

	m := macaron.New()
	if !setting.ProdMode {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Statics(macaron.StaticOptions{
		SkipLogging: setting.ProdMode,
	}, "custom/public", "public", models.HTMLRoot))
	m.Use(i18n.I18n(i18n.Options{
		Files:       setting.Docs.Locales,
		DefaultLang: setting.Docs.Langs[0],
	}))
	tplDir := "templates"
	if setting.Page.UseCustomTpl {
		tplDir = "custom/templates"
	}
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory: tplDir,
	}))
	m.Use(context.Contexter())

	m.Get("/", routes.Home)
	m.Get("/docs", routes.Docs)
	m.Get("/docs/images/*", routes.DocsStatic)
	m.Get("/docs/*", routes.Protect, routes.Docs)
	m.Post("/hook", routes.Hook)
	m.Get("/search", routes.Search)
	m.Get("/*", routes.Pages)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)
	log.Info("%s Listen on %s", setting.Site.Name, listenAddr)
	log.Fatal("Fail to start Peach: %v", http.ListenAndServe(listenAddr, m))
}
