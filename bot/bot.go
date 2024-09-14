package bot

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"time"
	"wxbot/getdata"
)

var laststring string
var lasttime time.Time

func Bot() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())

	dispatcher := openwechat.NewMessageMatchDispatcher()

	//只处理消息类型为文本类型的消息
	dispatcher.OnGroup(func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		if !msg.IsText() {
			return
		}
		fmt.Println("Text: ", msg.Content)
		var matchingCoins []string
		var matchingfutureCoins []string
		qbool, qtext := getdata.ContainsDexWithEnglishString(msg.Content)

		if qbool {
			sendtext := getdata.GetDexPrice(qtext)
			if sendtext != "" {
				msg.ReplyText(sendtext)
			}
		} else {
			matchingCoins = getdata.FindMatchingCoins(msg.Content, getdata.Coins)
			matchingfutureCoins = getdata.FindMatchingCoins(msg.Content, getdata.Futurecoins)
		}

		if len(matchingCoins) > 0 || len(matchingfutureCoins) > 0 {
			if time.Since(lasttime) < 2*time.Second && laststring == matchingCoins[len(matchingCoins)-1] {
				return
			}

			laststring = matchingCoins[len(matchingCoins)-1]
			lasttime = time.Now()
			fmt.Println("现货:", matchingCoins)
			fmt.Println("合约:", matchingfutureCoins)
			sendtext := ""
			for _, item := range matchingCoins {
				pricestring := getdata.GetPairPrice(item)
				sendtext = sendtext + pricestring + "\n"
			}
			for _, item := range matchingfutureCoins {
				pricestring := getdata.GetFuturePairPrice(item)
				sendtext = sendtext + pricestring + "\n"
			}

			if sendtext != "" {
				msg.ReplyText(sendtext)
			}
		}
	})

	//只处理消息类型为文本类型的消息
	dispatcher.OnFriend(func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		if !msg.IsText() {
			return
		}
		fmt.Println("Text: ", msg.Content)
		var matchingCoins []string
		var matchingfutureCoins []string
		qbool, qtext := getdata.ContainsDexWithEnglishString(msg.Content)

		if qbool {
			sendtext := getdata.GetDexPrice(qtext)
			if sendtext != "" {
				msg.ReplyText(sendtext)
			}
		} else {
			matchingCoins = getdata.FindMatchingCoins(msg.Content, getdata.Coins)
			matchingfutureCoins = getdata.FindMatchingCoins(msg.Content, getdata.Futurecoins)
		}

		if len(matchingCoins) > 0 || len(matchingfutureCoins) > 0 {
			if time.Since(lasttime) < 2*time.Second && laststring == matchingCoins[len(matchingCoins)-1] {
				return
			}

			laststring = matchingCoins[len(matchingCoins)-1]
			lasttime = time.Now()
			fmt.Println("现货:", matchingCoins)
			fmt.Println("合约:", matchingfutureCoins)
			sendtext := ""
			for _, item := range matchingCoins {
				pricestring := getdata.GetPairPrice(item)
				sendtext = sendtext + pricestring + "\n"
			}
			for _, item := range matchingfutureCoins {
				pricestring := getdata.GetFuturePairPrice(item)
				sendtext = sendtext + pricestring + "\n"
			}

			if sendtext != "" {
				msg.ReplyText(sendtext)
			}
		}
	})

	// 注册消息回调函数
	bot.MessageHandler = dispatcher.AsMessageHandler()
	//// 注册登陆二维码回调
	//bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	//
	//// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
