package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const MAXCOUNT int = 10000

var (
	
	baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

	slowScrollScript string = "window.scrollTo({top: document.body.scrollHeight, behavior: 'smooth'});"
	// fastScrollScript string = "window.scrollTo(0, document.body.scrollHeight);"
)

func main() {
	// 크롬 드라이버 실행 옵션 설정
	service:= makeChromeDriver()
	defer service.Stop()

	// 크롬 브라우저 실행 설정
	wd := makeBrowser()
	defer wd.Quit()

	// 마지막 페이지 번호 얻기
	totalPage := getPageMax(wd)
	fmt.Println(totalPage)

	// 전체 페이지 url 획득
	for i := 0; i < totalPage; i++ {
		fmt.Println(i,"번째 페이지 url", getPageURL(i))
	}  

	// 웹 페이지 열기
	openPage(wd, baseURL)

	// JavaScript를 사용하여 페이지 끝까지 스크롤
	scrollDown(wd, slowScrollScript)

	// 페이지의 HTML 소스코드 가져오기
	htmlSrc := getHtml(wd)

	// goquery를 사용하여 HTML 파싱하기
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSrc))
	if err != nil {
			panic(err)
	}

	doc.Find("nav[role='navigation']").Each(func(i int, s *goquery.Selection){
		// s.Find("a[data-testid='pagination-page-next']")
	})

	element, err := wd.FindElement(selenium.ByCSSSelector, "a[data-testid='pagination-page-next']")
	if err != nil {
  panic(err)
	}
	if err := element.Click(); err != nil {
  panic(err)
	}

	// JavaScript를 사용하여 페이지 끝까지 스크롤
	if _, err := wd.ExecuteScript(slowScrollScript, nil); err != nil {
	panic(err)
  }


	// 2초 대기 후 브라우저 종료
	time.Sleep(2 * time.Second)
}

// 브라우저 생성파트

// makeChromeDriver 크롬 드라이버 실행 옵션 설정 
func makeChromeDriver() *selenium.Service{
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("./chromedriver", 8080, opts...)
	if err != nil {
		panic(err)
	}
	return service
}

// makeBrowser 크롬 브라우저 실행 설정
func makeBrowser() selenium.WebDriver{
	// 크롬 브라우저 실행 설정
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
	Args: []string{
		"--disable-notifications", // 알림창 끄기
		"--disable-popup-blocking", // 팝업창 끄기
		// "--headless", // 브라우저를 헤드리스 모드로 실행
	},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 8080))
	if err != nil {
		panic(err)
	}

	return wd
}

// ----------------------------------------------

// 브라우저 기능 파트

// openPage 웹 페이지 열기
func openPage(wd selenium.WebDriver, url string){
	if err := wd.Get(url); err != nil {
		panic(err)
	}
}

// scrollDown JavaScript를 사용하여 페이지 끝까지 스크롤
func scrollDown(wd selenium.WebDriver, script string){
	if _, err := wd.ExecuteScript(script, nil); err != nil {
		panic(err)
		}
}

// getHtml  페이지의 HTML 소스코드 가져오기
func getHtml(wd selenium.WebDriver) string{
	htmlSrc, err := wd.PageSource()
	if err != nil {
		panic(err)
	}
	return htmlSrc
}





// getPageMax 총 페이지 확인 
func getPageMax(wd selenium.WebDriver) int {
	pageURL := baseURL + "&start=" + strconv.Itoa(MAXCOUNT)
	fmt.Println("page URL", pageURL)
	openPage(wd, pageURL)
	
	// 하단으로 내리기
	scrollDown(wd, slowScrollScript)

	// 페이지의 HTML 소스코드 가져오기
	htmlSrc := getHtml(wd)
	
	// goquery를 사용하여 HTML 파싱하기
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSrc))
	if err != nil {
			panic(err)
	}

	lastDiv := doc.Find("nav[role='navigation']").Children().Last()
	buttonText := lastDiv.Find("button").Text()
	lastPageNumber, _ := strconv.Atoi(buttonText)
	return lastPageNumber
}


// getPage 페이지 url얻기
func getPageURL(page int) string {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	return pageURL
}