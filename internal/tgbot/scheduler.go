package bot

import (
	"log"
	"time"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"
)

func (tgBot *TgBot) scheduler(sm SotrageManager, vpn VPNManager, tm TgMethods) {

	const op = "tgbot.scheduler"
	tgBot.task(sm, vpn, tm)
	go func() {
		for {
			duration := nextScheduleTime()
			log.Printf("Next task scheduled to run in %v\n", duration)
			timer := time.NewTimer(duration)
			<-timer.C
			tgBot.task(sm, vpn, tm)
		}
	}()
}

func nextScheduleTime() time.Duration {
	now := time.Now()
	// Определяем ближайшее время 12:00 или 00:00
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	if now.Hour() < 12 {
		next = next.Add(time.Hour * time.Duration(12-now.Hour()))
	} else {
		next = next.Add(time.Hour * time.Duration(24-now.Hour()))
	}

	// Если текущее время совпадает с 12:00 или 00:00, добавляем 12 часов
	if now.Equal(next) {
		next = next.Add(12 * time.Hour)
	}

	return next.Sub(now)
}

func (tgBot *TgBot) task(sm SotrageManager, vpn VPNManager, tm TgMethods) {
	const op = "tgbot.scheduler"
	ls := sm.NewLinkStorage()
	keys, err := ls.Find(&db.LinkFilter{TimeEnd: time.Now().Format(layout), State: models.StateAllowed}, &db.OrderLinks{}, 0, 10000)
	if err != nil {
		log.Println(op, err)
	} else {
		if len(keys) == 0 {
			return
		}
		for _, k := range keys {
			k.State = models.StateBanned
			err = vpn.UpdateKey(k)
			if err != nil {
				log.Println(op, err)
			}
			user := models.User{UserId: k.UserId, Command: "Notification"}
			user.SetCommandData(k.Link)
			tgBot.sendMsgToUser(&user, tm.MessageManager(&user))
		}
	}

	// if err != nil {
	// 	log.Fatal(err)
	// }
	tx, _, err := ls.Update(nil, &models.Link{State: models.StateBanned}, &db.LinkFilter{TimeEnd: time.Now().Format(layout), State: models.StateAllowed})
	if err != nil {
		tx.Rollback()
		log.Println(op, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Println(op, err)
	}
}
