# Broadcast To Telegram Chat with Telegram API

Referensi :
https://techthoughts.info/how-to-create-a-telegram-bot-and-send-messages-via-api/

Contoh :
https://api.telegram.org/bot429153606:AAH1BoL-eUMSKaYTCi9JzBojlRDQlS2otMY/sendMessage?chat_id=85390240&text=Hello+World

Dapetin chat id:
https://telegram.me/userinfobot


func BroadcastTelegram(text string) (gorequest.Response,string,[]error) {
	token := "429153606:AAH1BoL-eUMSKaYTCi9JzBojlRDQlS2otMY"
	chatId := "85390240"
	request := gorequest.New()
	return request.Get("https://api.telegram.org/bot"+token+"/sendMessage?chat_id="+chatId+"&text="+text).
		End()
}