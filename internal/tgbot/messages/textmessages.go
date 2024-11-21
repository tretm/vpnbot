package messages

import "html"

// common variables
var (
	BtmBackText    string
	BtmBackValue   string
	BtmCancelText  string = html.UnescapeString("&#8617;") + " Отмена"
	BtmOkText      string = html.UnescapeString("&#127383;")
	BtmCancelValue string = "Cancel"
)

// start meny
var (
	MESSAGESTART string = "Доступ к VPN осуществляется через приложение <b>V2Ray VPN</b>\n\n" +
		"Наш сервис предоставляет:\n" +
		html.UnescapeString("&#128526;") + " Полную анонимность\n" +
		html.UnescapeString("&#128640;") + " Высокую скорость и стабильное подключение\n" +
		html.UnescapeString("&#128737;") + " Защиту от рекламы и вредоносных программ\n" +
		html.UnescapeString("&#128241;&#128187;") + " Работа на всех устройствах\n\n" +
		"<b>Мы не ограничиваем скорость подключения и трафик</b>" + html.UnescapeString("&#10071;&#10071;&#10071;")

	BtmVpnTrialText    string = html.UnescapeString("&#128273;") + " Попробовать бесплатно"
	BtmVpnTrialValue   string = "Trial"
	BtmVpnKeyText      string = html.UnescapeString("&#128477;") + " Мои VPN ключи"
	BtmVpnKeyValue     string = "VPNKey"
	BtmBalancAddText   string = html.UnescapeString("&#128176;") + " Пополнить баланс"
	BtmBalancAddValue  string = "BalancAdd"
	BtmHistoryPayText  string = html.UnescapeString("&#128467;") + " История платежей"
	BtmHistoryPayValue string = "HistoryPay"
	BtmHelpText        string = html.UnescapeString("&#128214;") + " Помощь"
	BtmHelpValue       string = "Help"
	BtmReferalText     string = "Реферальная программа"
	BtmReferalValue    string = "Referal"
)

// trial  menu
var (
	MessageTrial string = "Для активации пробного ключа от VPN Вам необходимо:\n\n" +
		html.UnescapeString("1&#8419;") + " Установить приложение <b>V2Ray VPN</b> для вашей операционной системы;\n" +
		html.UnescapeString("2&#8419;") + " Нажать кнопку <b>Попробовать</b>;\n" +
		html.UnescapeString("3&#8419;") + " Скопировать полученный VPN ключ и добавить в ваше приложение <b>V2Ray VPN</b>;\n\n" +
		html.UnescapeString("4&#8419;") + " Наслаждаться доступом к Вашим ресурсам!"
	BtmGenerateKeyText  string = "Попробовать"
	BtmGenerateKeyValue string = "TryVPN"

	MessageTrialGood string = "Срок действия ключа (в сутках): %d\n\n" +
		"Скопируйте и добавьте следующий ключ в приложение <b>V2Ray VPN</b>\n\n<code>%s</code>"
	MessageTrialExist      string = "У Вас уже есть ключ от нашего VPN: %s\n%s"
	MessageTrialIsFinished string = "Пробный период закончен"

	BtmHowToUseText  string = html.UnescapeString("&#128214;") + " Инструкция по подключению"
	BtmHowToUseValue string = "Instruction"
)

// pay menu
var (
	BtmPayMenuText  string = "Выберите сумму, на которую хотите пополнить баланс:"
	BtmPayMenuValue string = "PayAmount"

	MsgTitelPay       string = "Пополнение баланса vpnGiga_bot"
	MsgDescriptionPay string = "Это однократный платеж в Телеграм (не подписка) за сервис vpnGiga_bot."

	MsgPayAdded string = "Вы пополнили баланс на: <b>%d</b> ₽\nВаш текущий баланс <b>%d</b> ₽ "

	MsgPayUrlText string = "Вам необходимо оплатить <b>%d</b> ₽\n" +
		"Нажмите кнопку <b>Оплатить</b>\n" +
		"Оплата производится через сервис <b>yoomoney.ru</b>\n\n<b>Оплатить можно банковской картой или из кошелька yoomoney</b>"
	BtmPayUrlText    string = html.UnescapeString("&#128717;") + " Оплатить"
	BtmPayCheckText  string = "Я оплатил"
	BtmPayCheckValue string = "checkpay"

	MsgPayNotFoundText string = "Ваш платеж не найден, возможно, вы не произвели оплату или платеж еще не зачислен, попробуйте проверить состояние платежа позже"
	MsgPayFoundText    string = "Ваш платеж успешно зачислен вам на баланс, вернитесь в главное меню нажав кнопку Отмена"
)

// pay history
var (
	MsgPayHistoryEmpty string = "Вы не производили зачислений денежных средств в наш сервис"
	MsgPayHistory      string = "Детализация по Вашему счету:\n%s"
	BtmNextPageText    string = html.UnescapeString("&#10145;")
	BtmNextPageValue   string = "nextpage"
	BtmPrevPageText    string = html.UnescapeString("&#11013;")
	BtmPrevPageValue   string = "prevpage"
)

// vpn key menu
var (
	MsgKeysMenu string = "Ваш баланс <b>%d</b> ₽\n\n" +
		"В данном меню будут отображаься все ваши VPN ключи (действующие и не действующие).\n" +
		"Вы можете преобрести новый VPN ключ или продлить срок действия имеющегося не активного VPN ключа."
	BtmAddKeyText        string = html.UnescapeString("&#127381;") + " Получить новый ключ"
	BtmAddKeyValue       string = "AddKeyMenu"
	BtmVpnKeyDetaleText  string = "%s %d. Ключ-%s"
	BtmVpnKeyDetaleValue string = "keydetale"
)

// price menu
var (
	MsgPriceMeny string = "Ваш баланс: <b>%d</b>  ₽\n\n" +
		"Выберите срок на который хотите приобрести  VPN ключ:"

	MsgPriceExtendMenu string = "Ваш баланс: <b>%d</b>  ₽\n\n" +
		"Выберите срок на который хотите продлить действие VPN ключа:"
	BtmPriceText       string = "%s - %d  ₽"
	BtmPriceValue      string = "priceValue" // Первый параметр период на который покупают, второй - цена покупки
	BtmPiceExtendValue string = "priceExtendValue"
)

// Buy text
var (
	MsgBuyOkText  string = "Ваш ключ для VPN:\n\n<code>%s</code>\n\nКлюч будет действовать до <b>%s</b>"
	MsgBuyErrText string = "У вас недостаточно средств для преобретения VPN ключа на выбранный вами период"
)

// Help menu
var (
	MsgHelpText string = "<b><a href='https://t.me/VPNGiga_bot'>@VPNGiga_bot</a></b> работает на основе протокола <b>vless</>, данный протокол обеспечивает безопасную связь между клиентом и сервером Xray\n" +
		"<b>Установите одно из приложений поддерживающих <b>vless</b> для вашей OS</b>" + html.UnescapeString("&#128071;") +
		"\n\n" +
		html.UnescapeString("&#127823;") +
		" <b>iOs: <a href='https://apps.apple.com/ru/app/streisand/id6450534064'>Streisand</a>, <a href='https://apps.apple.com/ru/app/foxray/id6448898396'>FoXray</a></b>" +
		"\n\n" +
		html.UnescapeString("&#129302;") +
		" <b>Android: <a href='https://play.google.com/store/apps/details?id=com.v2ray.ang'>v2rayNG</a>, <a href='https://play.google.com/store/apps/details?id=moe.nb4a'>NekoBox</a></b>" +
		"\n\n" +
		html.UnescapeString("&#128421;") +
		" <b>Windows: <a href='https://apps.apple.com/ru/app/streisand/id6450534064'>Furious</a>, <a href='https://github.com/InvisibleManVPN/InvisibleMan-XRayClient'>InvisibleMan-XRayClient</a></b>" +
		"\n\n" +
		html.UnescapeString("&#128039;") +
		" <b>Linux: <a href='https://github.com/MatsuriDayo/nekoray'>Nekoray</a></b>\n\n" +
		"<b>Подключите полученные VPN ключи в приложение:</b>\n" +
		"Вам нужно лишь сделать их импорт в приложение. Во всех приложениях это называется примерно “Add config from clipboard”. Приложение само берёт из буфера обмена строку с VPN ключем, либо вам нужно его вставить в поле для ввода ключа." +
		"После импорта у вас появится подключение. Обычно его нужно выделить и нажать кнопку снизу. Либо в десктопных приложениях правой кнопкой на подключение и в контекстном меню выбрать “Start”. Подключений может быть несколько. И между ними можно легко переключаться."

	BtmOSValue      string = "selectos"
	BtmAndroidText  string = html.UnescapeString("&#129302;") + " Android"
	BtmAndroidValue string = "Android"
	BtmIOSText      string = html.UnescapeString("&#127823;") + " iOs"
	BtmIOSValue     string = "iOs"
	BtmWindowsText  string = html.UnescapeString("&#128421;") + " Windows"
	BtmWindowsValue string = "Windows"
	BtmLinuxText    string = html.UnescapeString("&#128039;") + " Linux"
	BtmLinuxValue   string = "Linux"
	BtmSupportText  string = html.UnescapeString("&#127384;") + " Техподдержка"

	MsgAndroidText string = "Android инстркуция по подключению"
	MsgIOSText     string = "iOs инстркуция по подключению"
	MsgWindowsText string = "Windows инстркуция по подключению"
	MsgLinuxText   string = "Linux инстркуция по подключению"
)

// Detalisation key
var (
	MsgKeyDetalisationOk string = html.UnescapeString("&#9989;") + " Ваш VPN ключ:\n\n" +
		"<code>%s</code>\n\n" +
		"Действует до: <b>%s</b>\n" +
		"Осталось дней: %d"
	MsgKeyDetalisationBad string = html.UnescapeString("&#10060;") + " Ваш VPN ключ:\n\n" +
		"<code>%s</code>\n\n" +
		"Срок действия истек: <b>%s</b>\n"
	BtmExtendText string = html.UnescapeString("&#128260;") +
		" Продлить ключ"
	BtmExtndValue string = "extend"
)

// Confirm buy
var (
	MsgConfirmBuyText string = "Вы приобретаете новый VPN ключ\n" +
		"- Срок действия (дней): %d\n" +
		"- За: %s ₽"
	MsgConfirmExtendText string = "Вы продлеваете действие VPN ключа\n" +
		"- Ключ: %s\n" +
		"- Срок действия (дней): %d\n" +
		"- За: %s ₽"
	BtmConfirmYesText  string = html.UnescapeString("&#9989;") + "ДА"
	BtmConfirmYesValue string = "ConfirmYes"
	BtmConfirmNoText   string = html.UnescapeString("&#10060;") + "НЕТ"
	BtmConfirmNoValue  string = "ConfirmNo"
)

var (
	MsgNotificationText  string = html.UnescapeString("&#9888;") + "Истек срок действия VPN ключа:\n\n<b>%s</b>\n\nДля его продления вы можете перейти в раздел главного меню\n<b>Мои VPN ключи</b>"
	BtmNotificationValue string = "Notification"
)
