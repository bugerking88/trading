package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

func main() {
	// 创建一个新的 context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf), // 启用日志
	)
	defer cancel()

	// 定义一个变量来存储 HTML 内容
	var htmlContent string

	// 执行 chromedp 任务
	err := chromedp.Run(ctx,
		// 导航到目标 URL
		chromedp.Navigate("https://www.coinglass.com/zh-TW/FundingRate"),
		// 等待指定元素的可见性，确保页面渲染完成
		chromedp.WaitVisible(".ant-table-row.ant-table-row-level-0", chromedp.ByQuery),
		// 获取渲染后的 HTML 内容
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		log.Fatal(err)
	}

	// 打印渲染后的 HTML 内容
	//fmt.Println("Rendered HTML Content: ", htmlContent)

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(stringToReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}
	// 找交易所名稱
	var symbolName []string
	doc.Find("div.symbol-name").Each(func(i int, s *goquery.Selection) {
		// 获取 <a> 元素的文本内容
		text := s.Text()
		symbolName = append(symbolName, text)
	})

	// 找到所有 <a class="shou"> 元素
	var aList []string
	doc.Find("a.shou").Each(func(i int, s *goquery.Selection) {
		// 获取 <a> 元素的文本内容
		text := s.Text()
		// 获取 ref 属性的值
		ref := s.AttrOr("href", "")
		// 提取 ref 中的最后部分（如 "DOT"）
		if ref != "" {
			parts := strings.Split(ref, "/")
			if len(parts) > 0 {
				aList = append(aList, parts[len(parts)-1]+":"+text)
			}
		}
	})

	// 打印结果
	//fmt.Println("Extracted Texts:")
	//for _, item := range aList {
	//	fmt.Println(item)
	//}
	//fmt.Println("交易所名:")
	//for _, item := range symbolName {
	//	fmt.Println(item)
	//}
	// 用于存储分组数据
	groupedData := make(map[string][]string)

	// 遍历 aList
	for index, item := range aList {
		// 获取每个元素的前缀（例如 BNX 或 PEPE）
		parts := strings.SplitN(item, ":", 2)
		if len(parts) > 1 {
			group := parts[0] // 前缀
			// 如果前缀是 BTC，且是基数行（奇数行），则保存
			if group == "BTC" && index%2 == 0 {
				groupedData[group] = append(groupedData[group], item)
			} else if group != "BTC" {
				// 非 BTC 前缀的数据不受限制，全部保存
				groupedData[group] = append(groupedData[group], item)
			}
		}
	}
	groupedData2 := make(map[string][]string)

	for group, items := range groupedData {
		for i, item3 := range items {
			groupedData2[group] = append(groupedData2[group], symbolName[i]+":"+item3)
		}

	}

	for group, items := range groupedData2 {
		fmt.Printf("Group: %s\n", group)
		for _, item := range items {
			fmt.Println(item)
		}
		fmt.Println() // 每个组之间加空行
	}
}

// 辅助函数，将 string 转换为 reader
func stringToReader(s string) *strings.Reader {
	return strings.NewReader(s)
}
