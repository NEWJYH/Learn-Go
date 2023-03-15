package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)


type extractedJob struct {
	id              string
	companyName     string
	title           string
	companyLocation string
	salary          string
	summary         string
}

// createExtractedJob 구조체 생성
func createExtractedJob(id, companyName, title, companyLocation, salary, summary string) extractedJob {
	return extractedJob{id:id, companyName:companyName, title:title, companyLocation:companyLocation, salary:salary, summary:summary}
}


var (
	baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"
	slowScrollScript string = "window.scrollTo({top: document.body.scrollHeight, behavior: 'smooth'});"
)


func main() {

	// jobs := 데이터 배열 타겟
	var jobs []extractedJob

	// 크롬 드라이버 실행 옵션 설정
	service:= makeChromeDriver()
	defer service.Stop()

	// 크롬 브라우저 실행 설정
	wd := makeBrowser()
	defer wd.Quit()

	// 마지막 페이지 번호 얻기
	totalPage := getMaxPageNumber(wd)
	fmt.Println("total pages count :", totalPage)
	// 0.1초
	time.Sleep(100 * time.Millisecond) // Millisecond 1/1000 초

	// 전체 페이지 url 획득
	for i := 0; i < totalPage; i++ {
		fmt.Println(i+1,"번째 페이지 url")
		time.Sleep(2 * time.Second)
		jobs = append(jobs, makeData(wd, getPageURL(i))...)
	}  
	
	// 2초 대기 후 브라우저 종료
	time.Sleep(2 * time.Second)
	writeCSV(jobs)
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
func getMaxPageNumber(wd selenium.WebDriver) int {
	// openPage 웹 페이지 열기
	openPage(wd, baseURL)
	
	// 페이지의 HTML 소스코드 가져오기
	htmlSrc := getHtml(wd)
	
	// goquery를 사용하여 HTML 파싱하기
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSrc))
	if err != nil {
			panic(err)
	}
	Div := doc.Find(".jobsearch-JobCountAndSortPane-jobCount")
	spanText := Div.Find("span").First().Text()
	
	// 정규표현식으로 문자열에서 숫자 추출
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(spanText, -1)
	count, _ := strconv.Atoi(strings.TrimSpace(strings.Join(matches,""))) 


	lastPageNumber := int(math.Ceil(float64(count) / 50))
	return lastPageNumber
}

// getPage 페이지 url얻기
func getPageURL(page int) string {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	return pageURL
}

// extractJob 실질적인 뜯어내는 부분
func extractJob(s *goquery.Selection) extractedJob {
	a := s.Find(".jobTitle a:first-child")
	id, _:= a.Attr("data-jk")
	companyName := cleanString(s.Find(".companyName").Text())
	companyLocation := cleanString(s.Find(".companyLocation").Text())
	title := cleanString(a.Find("span").Text())
	salary := cleanString(s.Find(".salary-snippet-container > div").Text())
	summary := cleanString(s.Find(".result-footer > .job-snippet").Text())
	return createExtractedJob(id, companyName, title, companyLocation, salary, summary)
}

// cleanString 문자열 공백제거
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")	
}


// makeData url 넣으면 브라우저 이동하여 데이터 뜯어옴
func makeData(wd selenium.WebDriver, URL string) []extractedJob {
	var jobs []extractedJob
	// 웹 페이지 열기
	openPage(wd, URL)
	
	time.Sleep(100 * time.Millisecond) // Millisecond 1/1000 초

	// JavaScript를 사용하여 페이지 끝까지 스크롤
	scrollDown(wd, slowScrollScript)

	// 페이지의 HTML 소스코드 가져오기
	htmlSrc := getHtml(wd)

	// goquery를 사용하여 HTML 파싱하기
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSrc))
	if err != nil {
		panic(err)
	}
	
	// 
	searchCards := doc.Find(".jobsearch-ResultsList > li")

	//
	searchCards.Each(func(i int, s *goquery.Selection){
		job := extractJob(s)
		if job.id != ""{
			jobs = append(jobs, job)
		}
	})
	return jobs
}



// writeCSV jobs 를 csv 파일로 저장
func writeCSV(jobs []extractedJob){
	// https://pkg.go.dev/encoding/csv	

	// os create File
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string {"ID","CompanyName","Title","CompanyLocation","Salary","Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk="+job.id, job.companyName, job.title, job.companyLocation, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}





// 클릭이벤트
// doc.Find("nav[role='navigation']").Each(func(i int, s *goquery.Selection){
	// 	// s.Find("a[data-testid='pagination-page-next']")
	// })

	// element, err := wd.FindElement(selenium.ByCSSSelector, "a[data-testid='pagination-page-next']")
	// if err != nil {
  // panic(err)
	// }
	// if err := element.Click(); err != nil {
  // panic(err)
	// }


	func checkErr(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}
	
	// func checkCode(res *http.Response) {
	// 	if res.StatusCode != 200 {
	// 		log.Fatalln("Request failed with Status:", res.StatusCode)
	// 	}
	// }