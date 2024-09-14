package getdata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CurrencyPairs struct {
	Id              string `json:"id"`
	Base            string `json:"base"`
	Quote           string `json:"quote"`
	Fee             string `json:"fee"`
	MinBaseAmount   string `json:"min_base_amount"`
	MinQuoteAmount  string `json:"min_quote_amount"`
	MaxQuoteAmount  string `json:"max_quote_amount,omitempty"`
	AmountPrecision int    `json:"amount_precision"`
	Precision       int    `json:"precision"`
	TradeStatus     string `json:"trade_status"`
	SellStart       int    `json:"sell_start"`
	BuyStart        int    `json:"buy_start"`
	MaxBaseAmount   string `json:"max_base_amount,omitempty"`
}
type FuturePairs struct {
	Last                  string `json:"last"`
	Low24H                string `json:"low_24h"`
	High24H               string `json:"high_24h"`
	Volume24H             string `json:"volume_24h"`
	ChangePercentage      string `json:"change_percentage"`
	FundingRateIndicative string `json:"funding_rate_indicative"`
	IndexPrice            string `json:"index_price"`
	Volume24HBase         string `json:"volume_24h_base"`
	Volume24HQuote        string `json:"volume_24h_quote"`
	Contract              string `json:"contract"`
	Volume24HSettle       string `json:"volume_24h_settle"`
	FundingRate           string `json:"funding_rate"`
	MarkPrice             string `json:"mark_price"`
	TotalSize             string `json:"total_size"`
	HighestBid            string `json:"highest_bid"`
	LowestAsk             string `json:"lowest_ask"`
	QuantoMultiplier      string `json:"quanto_multiplier"`
}

var Coins []string
var Futurecoins []string

func init() {
	host := "https://api.gateio.ws"
	prefix := "/api/v4"
	url := "/spot/currency_pairs"

	// 创建 http client
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", host+prefix+url, nil)
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// 解析 JSON 数据
	var pair []CurrencyPairs
	err = json.Unmarshal(body, &pair)
	if err != nil {
		panic(err)
	}
	//// 将解析后的数据存储到切片中
	//解析现货
	for _, currencyPair := range pair {
		if currencyPair.Quote == "USDT" {
			Coins = append(Coins, currencyPair.Base)
		}
	}
	log.Println("现货:", Coins)

	// 创建请求
	futurereq, err := http.NewRequest("GET", "https://api.gateio.ws/api/v4/futures/usdt/tickers", nil)
	if err != nil {
		panic(err)
	}
	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	futureresp, err := client.Do(futurereq)
	if err != nil {
		panic(err)
	}
	defer futureresp.Body.Close()

	// 读取响应
	futurebody, err := io.ReadAll(futureresp.Body)
	if err != nil {
		panic(err)
	}
	// 解析 JSON 数据
	var futurepairs []FuturePairs
	err = json.Unmarshal(futurebody, &futurepairs)
	if err != nil {
		panic(err)
	}

	//解析合约币种
	for _, futurepair := range futurepairs {
		parts := strings.Split(futurepair.Contract, "_")
		Futurecoins = append(Futurecoins, parts[0])
	}
	log.Println("合约:", Futurecoins)

}

// 定义函数，根据输入的币种查找交易对
func FindCurrencyPair(currency string, pairs []CurrencyPairs) (*CurrencyPairs, error) {
	currency = strings.ToUpper(currency) // 转换为大写
	for _, pair := range pairs {
		if pair.Base == currency || pair.Quote == currency {
			return &pair, nil
		}
	}
	return nil, fmt.Errorf("未找到匹配的交易对")
}

func FindMatchingCoins(text string, coins []string) []string {
	// 检测文本中是否包含中文字符或 HTTP 链接
	if containsChinese(text) || containsHTTP(text) {
		return []string{} // 直接返回空切片
	}
	// 将所有币种转换为大写，并用 "|" 连接成正则表达式
	pattern := fmt.Sprintf(`\b(%s)\b`, strings.Join(coins, "|"))
	re := regexp.MustCompile(pattern)
	// 查找所有匹配的币种
	matches := re.FindAllString(strings.ToUpper(text), -1)
	// 去重
	uniqueMatches := make(map[string]bool)
	for _, match := range matches {
		uniqueMatches[match] = true
	}

	// 将去重后的结果转换为切片
	var result []string
	for coin := range uniqueMatches {
		result = append(result, coin)
	}
	return result
}

// 检测文本中是否包含中文字符
func containsChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF { // 中文字符范围
			return true
		}
	}
	return false
}

// 检测文本中是否包含 HTTP 链接
func containsHTTP(text string) bool {
	re := regexp.MustCompile(`https?://`)
	return re.MatchString(text)
}

// 检测文本中是否包含 "dex" 后面跟着任何字符串的固定字段，并且 "dex" 必须在最前面
func ContainsDexWithEnglishString(text string) (bool, string) {
	re := regexp.MustCompile(`^dex\s+(.+)`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		log.Println(match[1])
		return true, match[1]
	}
	return false, ""
}

type PairPrice struct {
	CurrencyPair     string `json:"currency_pair"`
	Last             string `json:"last"`
	LowestAsk        string `json:"lowest_ask"`
	HighestBid       string `json:"highest_bid"`
	ChangePercentage string `json:"change_percentage"`
	BaseVolume       string `json:"base_volume"`
	QuoteVolume      string `json:"quote_volume"`
	High24H          string `json:"high_24h"`
	Low24H           string `json:"low_24h"`
	EtfNetValue      string `json:"etf_net_value"`
	EtfPreNetValue   string `json:"etf_pre_net_value"`
	EtfPreTimestamp  int    `json:"etf_pre_timestamp"`
	EtfLeverage      string `json:"etf_leverage"`
}

func GetPairPrice(base string) string {
	host := "https://api.gateio.ws"
	prefix := "/api/v4"
	url := "/spot/tickers"
	queryParam := "?currency_pair=" + base + "_USDT"
	//url = '/futures/usdt/contracts/BTC_USDT'
	// 创建 http client
	client := &http.Client{}
	// 创建请求
	//fmt.Println(host + prefix + url + queryParam)
	req, err := http.NewRequest("GET", host+prefix+url+queryParam, nil)
	if err != nil {
		return ""
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//log.Println(string(body))
	// 解析 JSON 数据
	var price []PairPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		return ""
	}
	if len(price) > 0 {
		log.Println(base, "现货Price:", "$", price[0].Last, "(", price[0].ChangePercentage, "%)")
		return fmt.Sprint(base, " Price: ", "$", price[0].Last, " (", price[0].ChangePercentage, "%)")
		//fmt.Printf("%s Price: $%.4f (%.2f%%)\n", base, price[0].Last, price[0].ChangePercentage)
	}
	return ""
}

type FuturePairPrice struct {
	Last                  string `json:"last"`
	Low24H                string `json:"low_24h"`
	High24H               string `json:"high_24h"`
	Volume24H             string `json:"volume_24h"`
	ChangePercentage      string `json:"change_percentage"`
	FundingRateIndicative string `json:"funding_rate_indicative"`
	IndexPrice            string `json:"index_price"`
	Volume24HBase         string `json:"volume_24h_base"`
	Volume24HQuote        string `json:"volume_24h_quote"`
	Contract              string `json:"contract"`
	Volume24HSettle       string `json:"volume_24h_settle"`
	FundingRate           string `json:"funding_rate"`
	MarkPrice             string `json:"mark_price"`
	TotalSize             string `json:"total_size"`
	HighestBid            string `json:"highest_bid"`
	LowestAsk             string `json:"lowest_ask"`
	QuantoMultiplier      string `json:"quanto_multiplier"`
}

func GetFuturePairPrice(base string) string {
	host := "https://api.gateio.ws"
	prefix := "/api/v4"
	url := "/futures/usdt/tickers"
	queryParam := "?contract=" + base + "_USDT"
	// 创建 http client
	client := &http.Client{}
	// 创建请求
	//fmt.Println(host + prefix + url + queryParam)
	req, err := http.NewRequest("GET", host+prefix+url+queryParam, nil)
	if err != nil {
		return ""
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//log.Println(string(body))
	// 解析 JSON 数据
	var price []FuturePairPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		return ""
	}
	if len(price) > 0 {
		log.Println(base, "合约Price:", "$", price[0].Last, "(", price[0].ChangePercentage, "%)")
		return fmt.Sprint(base, " 合约: ", "$", price[0].Last, " (", price[0].ChangePercentage, "%)")
		//fmt.Printf("%s Price: $%.4f (%.2f%%)\n", base, price[0].Last, price[0].ChangePercentage)
	}
	return ""
}

type DexPair struct {
	SchemaVersion string `json:"schemaVersion"`
	Pairs         []struct {
		ChainId     string   `json:"chainId"`
		DexId       string   `json:"dexId"`
		Url         string   `json:"url"`
		PairAddress string   `json:"pairAddress"`
		Labels      []string `json:"labels,omitempty"`
		BaseToken   struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"baseToken"`
		QuoteToken struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"quoteToken"`
		PriceNative string `json:"priceNative"`
		PriceUsd    string `json:"priceUsd"`
		Txns        struct {
			M5 struct {
				Buys  int `json:"buys"`
				Sells int `json:"sells"`
			} `json:"m5"`
			H1 struct {
				Buys  int `json:"buys"`
				Sells int `json:"sells"`
			} `json:"h1"`
			H6 struct {
				Buys  int `json:"buys"`
				Sells int `json:"sells"`
			} `json:"h6"`
			H24 struct {
				Buys  int `json:"buys"`
				Sells int `json:"sells"`
			} `json:"h24"`
		} `json:"txns"`
		Volume struct {
			H24 float64 `json:"h24"`
			H6  float64 `json:"h6"`
			H1  float64 `json:"h1"`
			M5  float64 `json:"m5"`
		} `json:"volume"`
		PriceChange struct {
			M5  float64 `json:"m5"`
			H1  float64 `json:"h1"`
			H6  float64 `json:"h6"`
			H24 float64 `json:"h24"`
		} `json:"priceChange"`
		Liquidity struct {
			Usd   float64 `json:"usd"`
			Base  float64 `json:"base"`
			Quote float64 `json:"quote"`
		} `json:"liquidity"`
		Fdv           int   `json:"fdv"`
		PairCreatedAt int64 `json:"pairCreatedAt,omitempty"`
		Info          struct {
			ImageUrl string        `json:"imageUrl"`
			Websites []interface{} `json:"websites"`
			Socials  []struct {
				Type string `json:"type"`
				Url  string `json:"url"`
			} `json:"socials"`
		} `json:"info,omitempty"`
	} `json:"pairs"`
}

func GetDexPrice(base string) string {
	url := "https://api.dexscreener.com/latest/dex/search/?q="
	// 创建 http client
	client := &http.Client{}
	// 创建请求
	//fmt.Println(host + prefix + url + queryParam)
	req, err := http.NewRequest("GET", url+base, nil)
	if err != nil {
		return ""
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//log.Println(string(body))
	// 解析 JSON 数据
	var pair DexPair
	err = json.Unmarshal(body, &pair)
	if err != nil {
		log.Println(err)
		return ""
	}
	sendtext := ""
	if len(pair.Pairs) > 0 {
		sendtext = sendtext + "代币名称:" + pair.Pairs[0].BaseToken.Name + "\n"
		sendtext = sendtext + "来源:" + pair.Pairs[0].ChainId + "-" + pair.Pairs[0].DexId + "\n"
		sendtext = sendtext + "代币地址:" + pair.Pairs[0].BaseToken.Address + "\n"
		sendtext = sendtext + "当前价格:" + "$" + pair.Pairs[0].PriceUsd + "\n"
		sendtext = sendtext + "底池数量:" + fmt.Sprint(pair.Pairs[0].Liquidity.Quote) + pair.Pairs[0].QuoteToken.Symbol + "($" + fmt.Sprint(pair.Pairs[0].Liquidity.Usd) + ")" + "\n"
		sendtext = sendtext + "5分钟:" + fmt.Sprint(pair.Pairs[0].PriceChange.M5) + "%" + "\n"
		sendtext = sendtext + "1小时:" + fmt.Sprint(pair.Pairs[0].PriceChange.H1) + "%" + "\n"
		sendtext = sendtext + "6小时:" + fmt.Sprint(pair.Pairs[0].PriceChange.H6) + "%" + "\n"
		sendtext = sendtext + "24小时:" + fmt.Sprint(pair.Pairs[0].PriceChange.H24) + "%" + "\n"
		sendtext = sendtext + "流动性:" + fmt.Sprint(pair.Pairs[0].Fdv) + "\n"
		sendtext = sendtext + "24小时交易量:" + "$" + fmt.Sprint(pair.Pairs[0].Volume.H24) + "\n"
		//log.Println(sendtext)
		return sendtext
	}
	return ""
}
