package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"fortune/internal/config"
	"fortune/internal/fortune"
	"fortune/internal/handler"
	"fortune/internal/payment"
	"fortune/internal/store"
	"fortune/web"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	db, err := store.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	if err := store.Migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	if err := store.SeedAdmin(db, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		log.Fatalf("初始化管理员失败: %v", err)
	}

	prov, err := payment.New(payment.Config{
		Provider: cfg.Payment.Provider,
		BaseURL:  cfg.App.BaseURL,
		Wechat: payment.WechatConfig{
			AppID:           cfg.Payment.Wechat.AppID,
			MchID:           cfg.Payment.Wechat.MchID,
			MchSerialNo:     cfg.Payment.Wechat.MchSerialNo,
			PrivateKeyPath:  cfg.Payment.Wechat.MchPrivateKey,
			PlatformCertPath: cfg.Payment.Wechat.PlatformCertPath,
			APIv3Key:        cfg.Payment.Wechat.APIv3Key,
			NotifyURL:       cfg.Payment.Wechat.NotifyURL,
			PayType:         cfg.Payment.Wechat.PayType,
		},
	})
	if err != nil {
		log.Fatalf("初始化支付失败: %v", err)
	}

	h := handler.New(db, fortune.NewEngine(), prov, cfg)

	distFS, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		log.Fatalf("加载前端资源失败: %v", err)
	}
	fileServer := http.FileServer(http.FS(distFS))

	// 静态资源缓存策略：入口 html 不缓存（保证前端更新即时生效），
	// 带哈希的 /assets/* 长期缓存。
	static := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}
		fileServer.ServeHTTP(w, r)
	})

	root := chi.NewRouter()
	root.Mount("/api", h.Routes())

	// SEO：robots.txt 与 sitemap.xml（供搜索引擎与 AI Agent 抓取）
	root.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /api/\n\nSitemap: %s/sitemap.xml\n", cfg.App.BaseURL)
	})
	root.Get("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n<url><loc>%s/</loc><changefreq>daily</changefreq><priority>1.0</priority></url>\n</urlset>\n", cfg.App.BaseURL)
	})

	root.Handle("/*", static)

	if cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "" {
		log.Printf("%s 已启动，监听 %s（HTTPS，支付通道：%s）", cfg.App.Name, cfg.Server.Addr, cfg.Payment.Provider)
		if err := http.ListenAndServeTLS(cfg.Server.Addr, cfg.Server.TLSCert, cfg.Server.TLSKey, root); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s 已启动，监听 %s（HTTP，支付通道：%s）", cfg.App.Name, cfg.Server.Addr, cfg.Payment.Provider)
		if err := http.ListenAndServe(cfg.Server.Addr, root); err != nil {
			log.Fatal(err)
		}
	}
}