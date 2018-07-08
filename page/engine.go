package page

import (
	"github.com/fedesog/webdriver"
	"goBoss/utils"
	"fmt"
	"log"
	"errors"
	cf "goBoss/config"
)

type Engine struct {
	Dr *webdriver.ChromeDriver
	Ss *webdriver.Session
}

type Engineer interface {
	MaxWindow() error
	Start()
	OpenBrowser()
	SetWindow(width, height int) error
	GetElement(root, name string) *Element
	Driver() *webdriver.ChromeDriver
	Session() *webdriver.Session
	Screen() ([]byte, error)
	GetUrl() (string, error)
}

func (w *Engine) Session() *webdriver.Session {
	return w.Ss
}

func (w *Engine) Driver() *webdriver.ChromeDriver {
	return w.Dr
}

func (w *Engine) GetElement(root, name string) *Element {
	ele, ok := Page[root][name]
	if !ok {
		log.Panicf("page/element.json未找到root: [%s] key: [%s]", root, name)
	}
	return &ele
}

func (w *Engine) MaxWindow() error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s"}`, w.Ss.Id)
	req := utils.Request{
		Data:   p,
		Method: "POST",
		Url:    fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/maximize", w.Dr.Port, w.Ss.Id),
	}

	res := req.Http()
	if !res["status"].(bool) {
		log.Printf("response: %+v", res)
		return errors.New(fmt.Sprintf(`最大化窗口失败, 请检查!%+v`, res["msg"]))
	}
	return nil
}

func (w *Engine) Start() {
	var err error
	w.Dr.Start()
	args := make([]string, 0)
	if cf.Config.Headless {
		args = append(args, "--headless")
	}
	desired := webdriver.Capabilities{
		"Platform":           "Mac",
		"goog:chromeOptions": map[string][]string{"args": args, "extensions": []string{}},
		"browserName":        "chrome",
		"version":            "",
		"platform":           "ANY",
	}
	required := webdriver.Capabilities{}
	w.Ss, err = w.Dr.NewSession(desired, required)
	if err != nil {
		log.Printf("open browser failed: %s", err.Error())
	}

}

func (w *Engine) OpenBrowser() {
	w.Ss.Url(cf.Config.LoginURL)
	err := w.SetWindow(1600, 900)
	if err != nil {
		log.Panicf("最大化浏览器失败!!!Msg: %s", err.Error())
	}
	w.Ss.SetTimeoutsImplicitWait(cf.Config.WebTimeout)
}

func (w *Engine) SetWindow(width, height int) error {
	p := fmt.Sprintf(`{"windowHandle": "current", "sessionId": "%s", "height": %d, "width": %d}`, w.Ss.Id, height, width)
	req := utils.Request{
		Data:   p,
		Method: "POST",
		Url:    fmt.Sprintf("http://127.0.0.1:%d/session/%s/window/current/size", w.Dr.Port, w.Ss.Id),
	}

	res := req.Http()
	if !res["status"].(bool) {
		log.Printf("response: %+v", res)
		return errors.New(fmt.Sprintf(`设置浏览器窗口失败, 请检查!%+v`, res["msg"]))
	}
	return nil
}

func (w *Engine) Close() {
	w.Ss.CloseCurrentWindow()
	w.Dr.Stop()

}

func (w *Engine) Screen() ([]byte, error) {
	return w.Ss.Screenshot()
}

func (w *Engine) GetUrl() (string, error) {
	return w.Ss.GetUrl()
}