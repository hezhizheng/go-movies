package tmpl

import (
	"embed"
	"html/template"
)

// main 函数中定义的全局变量 在其他包中无法调用 ，重新定义新包实现全局变量
var (
	//go:embed *.html
	embedTmpl embed.FS
	GoTpl     = template.Must(template.ParseFS(embedTmpl, "*.html")) // 利用 air 监听文件变动 实时重新加载。修改html无须手动重启服务
	//GoTpl = template.Must(template.ParseGlob("./views/tmpl/*.html"))
)
