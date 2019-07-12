package curlToRequestGolang

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func dealwithCookies(r string) {
	cookie := regexp.MustCompile(`cookie.*?(?:\' )|Cookie.*?(?:\' )|cookie.*?(?:\')`)
	first := cookie.FindString(r)
	if first != " " {
		second := strings.TrimSpace(first)
		third := strings.Split(strings.Split(second[:len(second)-1], ": ")[1], "; ")
		fmt.Println("    req.Cookies = map[string]string{")
		for _, i := range third {
			cookiebody := strings.SplitN(i, "=", 2)
			cookiebodys := "    `" + strings.Join(cookiebody, "`: `") + "`,"
			fmt.Println(cookiebodys)
		}
		fmt.Println("    }")
	}
}

func dealwithParameter(parameter string, r string) string {
	parameterLower := strings.ToLower(parameter)
	pa := regexp.MustCompile(parameterLower + ".*?(?:' )+")
	para := strings.Replace(strings.Trim(pa.FindString(r), "' "), ": ", "`: `", 2)
	if para == "" {
		pa = regexp.MustCompile(parameter + ".*?(?:' )+")
		para = strings.Replace(strings.Trim(pa.FindString(r), "' "), ": ", "`: `", 2)
		if para == "" {
			return ""
		} else {
			return "    `" + para + "`,\n"
		}
	}
	return "    `" + para + "`,\n"
}

func dealwithHeaders(r string) {
	//分别调用关键字函数最后进行汇总打印
	AcceptEncoding := dealwithParameter("Accept_Encoding", r)
	Origin := dealwithParameter("Origin", r)
	AcceptLanguage := dealwithParameter("Accept-Language", r)
	UpgradeInsecureRequests := dealwithParameter("Upgrade-Insecure-Requests", r)
	UserAgent := dealwithParameter("User-Agent", r)
	ContentType := dealwithParameter("Content-Type", r)
	Accept := dealwithParameter("Accept:", r)
	CacheControl := dealwithParameter("Cache-Control", r)
	Referer := dealwithParameter("Referer", r)
	Connection := dealwithParameter("Connection", r)
	authority := dealwithParameter("Authority", r)
	xrequestedwith := dealwithParameter("x-requested-with", r)
	headers := "    req.Headers = map[string]string{\n" + Origin + AcceptEncoding + AcceptLanguage + UpgradeInsecureRequests + UserAgent + ContentType + Accept + authority + xrequestedwith + CacheControl + Referer + Connection + "    }" + "\n"
	fmt.Println(headers)
}

func dealUrl(r string) string {
	url := regexp.MustCompile(`curl.*?(?:\' )`).FindString(r)
	urls := strings.Trim(strings.Split(url, " ")[1], "'")
	return urls
}
func main() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("请输入curl(bash):")
	r, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(`
	package main
	func main() {
	c := &http.Client{}
	req := request.NewRequest(c)`)
	dealwithCookies(r)
	dealwithHeaders(r)
	fmt.Printf("    resp, err := req.Get(`"+"%s"+"`)\n"+
		"    if err != nil{\n"+
		"        panic(err)\n"+
		"    }\n"+
		"    defer resp.Body.Close()\n"+
		"    j, err := resp.Text()\n"+
		"    fmt.Println(j)\n", dealUrl(r))
	fmt.Println("    }")
}
