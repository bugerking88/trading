package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"strconv"
	"strings"
)

type GroupStats struct {
	Max     float64
	Min     float64
	MaxItem string
	MinItem string
}

func main() {
	// 创建 Gin 实例
	r := gin.Default()

	// 定义一个路由，访问 /fetch 时触发抓取和处理
	r.GET("/fetch", fetchHandler)
	r.GET("/fetch2", fetchHandler2)

	// 启动 Gin Web 服务器，监听 8080 端口
	log.Println("Server started at :8080")
	r.Run(":8080")
}

// fetchHandler 处理抓取 HTML 内容的请求
func fetchHandler(c *gin.Context) {
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
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error: %v", err)})
		return
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(stringToReader(htmlContent))
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error parsing HTML: %v", err)})
		return
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

	// 找每組最大最小
	// 遍历 groupedData2，找出最大和最小的百分比
	groupStatsRe := make(map[string]GroupStats)

	for group, items := range groupedData2 {
		var max float64 = -math.MaxFloat64
		var min float64 = math.MaxFloat64
		var maxItem, minItem string

		for _, item := range items {
			// 提取百分比部分
			parts := strings.Split(item, ":")
			if len(parts) > 2 && parts[2] != "-" {
				percentStr := strings.TrimSuffix(parts[2], "%") // 去掉百分号
				percent, err := strconv.ParseFloat(percentStr, 64)
				if err != nil {
					log.Printf("Error parsing percent: %v", err)
					continue
				}

				// 查找最大值和最小值
				if percent > max {
					max = percent
					maxItem = item
				}
				if percent < min {
					min = percent
					minItem = item
				}
			}
		}
		// 将最大值和最小值存入 map
		groupStatsRe[group] = GroupStats{
			Max:     max,
			Min:     min,
			MaxItem: maxItem,
			MinItem: minItem,
		}
		//// 打印最大值和最小值
		//fmt.Printf("Group: %s\n", group)
		//fmt.Printf("Max: %s -> %.4f%%\n", maxItem, max)
		//fmt.Printf("Min: %s -> %.4f%%\n", minItem, min)
	}

	// 返回结果作为 JSON 响应
	c.JSON(200, groupStatsRe)
}

// 辅助函数，将 string 转换为 reader
func stringToReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

// fetchHandler 处理抓取 HTML 内容的请求
func fetchHandler2(c *gin.Context) {
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
		chromedp.Navigate("https://www.coinglass.com/zh-TW/FrArbitrage"),
		// 等待指定元素的可见性，确保页面渲染完成
		chromedp.WaitVisible(".ant-table-tbody", chromedp.ByQuery),
		// 获取渲染后的 HTML 内容
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error: %v", err)})
		return
	}
	// 打印渲染后的 HTML 内容
	//fmt.Println("Rendered HTML Content: ", htmlContent)

	//// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(stringToReader(htmlContent))
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error parsing HTML: %v", err)})
		return
	}

	// 找價格
	var symbolName []string
	var percentages []string
	// 提取符号和百分比数据
	doc.Find("div.lh30").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		// 判断当前文本是否是百分比（假设百分比总是以 '%' 结尾）
		if len(text) > 0 && text[len(text)-1] == '%' {
			percentages = append(percentages, text)
		} else {
			symbolName = append(symbolName, text)
		}
	})

	// 幣種
	var coins []string
	doc.Find("div.symbol-name").Each(func(i int, s *goquery.Selection) {
		// 获取 <a> 元素的文本内容
		if i%3 == 0 {
			text := s.Text()
			coins = append(coins, text)
		}
	})
	for _, item := range coins {
		fmt.Println(item)
	}

	groupStatsRe := make(map[string][]string)
	// 合并符号和百分比数据并打印
	for i := 0; i < len(symbolName); i++ {
		// 保证每个 symbolName 都有对应的百分比
		if i < len(percentages) {
			groupStatsRe[coins[i/2]] = append(groupStatsRe[coins[i/2]], symbolName[i]+":"+percentages[i])
		}
	}
	var finalResult []string
	// 打印最终的结果，按照 coins 数组的顺序输出
	for _, coin := range coins {
		if items, exists := groupStatsRe[coin]; exists {
			for _, item := range items {
				finalResult = append(finalResult, coin+"/"+item)
			}
		}
	}

	// 返回结果作为 JSON 响应
	c.JSON(200, finalResult)
}
