package telegram

import (
	"context"
	"errors"
	"fmt"
	gmodels "github.com/b3liv3r/tgbot-for-gym/modules/gym/models"
	"github.com/b3liv3r/tgbot-for-gym/modules/subscription/models"
	tmodels "github.com/b3liv3r/tgbot-for-gym/modules/trainers/models"
	wmodels "github.com/b3liv3r/tgbot-for-gym/modules/wallet/models"
	"strings"
	"time"
)

func (r *TelegramRouter) UserIsExist(userID int) bool {
	_, err := r.Services.Userer.Profile(context.Background(), userID)
	if err != nil {
		return false
	}
	return true
}

func (r *TelegramRouter) Deposit(userID int, amount float64, description string) error {
	resp, balance, err := r.Services.Walleter.Transaction(context.Background(), wmodels.Transaction{
		UserId:      userID,
		Amount:      amount,
		Type:        "deposit",
		Description: description,
		Date:        time.Now(),
	})
	if err != nil {
		r.sendMessage(int64(userID), err.Error())
		return err
	}
	r.sendMessage(int64(userID), resp)
	r.sendMessage(int64(userID), fmt.Sprintf(descNewBalance, balance))
	return nil
}

func (r *TelegramRouter) Payment(userID int, amount float64, description string) error {
	resp, balance, err := r.Services.Walleter.Transaction(context.Background(), wmodels.Transaction{
		UserId:      userID,
		Amount:      amount,
		Type:        "payment",
		Description: description,
		Date:        time.Now(),
	})
	if err != nil {
		if err.Error() == "insufficient funds" {
			r.sendMessage(int64(userID), msgInsufficientFunds)
			return err
		}
		r.sendMessage(int64(userID), err.Error())
		return err
	}
	r.sendMessage(int64(userID), resp)
	r.sendMessage(int64(userID), fmt.Sprintf(descNewBalance, balance))
	return nil
}

func (r *TelegramRouter) Rollback(userID int, amount float64, description string) error {
	resp, balance, err := r.Services.Walleter.Transaction(context.Background(), wmodels.Transaction{
		UserId:      userID,
		Amount:      -amount,
		Type:        "rollback",
		Description: "rollback " + description,
		Date:        time.Now(),
	})
	if err != nil {
		r.sendMessage(int64(userID), err.Error())
		return err
	}
	r.sendMessage(int64(userID), resp)
	r.sendMessage(int64(userID), fmt.Sprintf(descNewBalance, balance))
	return nil
}

func parseTimeString(timeStr string) (time.Time, error) {
	layout := "02.01.2006_15:04"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, errors.New("incorrect time format example: 20.12.2004_15:04")
	}
	return parsedTime, nil
}

func isHistoryType(tsxType string) bool {
	switch tsxType {
	case "rollback":
		return true
	case "payment":
		return true
	case "deposit":
		return true
	default:
		return false
	}
}

func formatTransactionList(transactions []wmodels.Transaction) string {
	var builder strings.Builder

	// Перебор транзакций
	for _, txn := range transactions {
		// Формирование строки для каждой транзакции
		txnStr := fmt.Sprintf("Date: %s\n", txn.Date.Format("2006-01-02 15:04:05"))
		txnStr += fmt.Sprintf("Type: %s\n", txn.Type)
		txnStr += fmt.Sprintf("Description: %s\n", txn.Description)
		txnStr += fmt.Sprintf("Ammount: %.2f\n", txn.Amount)

		// Добавление строки транзакции к общей строке
		builder.WriteString(txnStr)
		builder.WriteString("\n") // Разделитель между транзакциями
	}

	return builder.String()
}

func subLvlPrice(lvl int) float64 {
	switch lvl {
	case 0:
		return costSubLvl0
	case 1:
		return costSubLvl1
	case 2:
		return costSubLvl2
	case 3:
		return costSubLvl3
	default:
		return -1
	}
}

func formatSubscriptionMessage(sub models.Subscription) string {
	startDateStr := sub.StartDate.Format("2006-01-02")
	endDateStr := sub.EndDate.Format("2006-01-02")
	// Формируем сообщение о подписке
	message := fmt.Sprintf("Subscription lvl: %d\n", sub.SubscriptionType)
	message += fmt.Sprintf("Start date: %s\n", startDateStr)
	message += fmt.Sprintf("End date: %s\n", endDateStr)

	return message
}

func formatGymList(gyms []gmodels.Gym) string {
	var builder strings.Builder

	for _, gym := range gyms {
		gymStr := fmt.Sprintf("ID: %d\n", gym.Id)
		gymStr += fmt.Sprintf("Address: %s\n", gym.Address)
		gymStr += fmt.Sprintf("Subscription lvl for access: %d\n", gym.SubLvl)

		builder.WriteString(gymStr)
		builder.WriteString("\n")
	}

	return builder.String()
}

func formatSchedulesList(schedules []gmodels.Schedules) string {
	var builder strings.Builder

	for _, sched := range schedules {
		startTimeStr := sched.StartTime.Format("15:04")
		endTimeStr := sched.EndTime.Format("15:04")
		gymStr := fmt.Sprintf("%s\n", sched.DayOfWeek)
		gymStr += fmt.Sprintf("Time open: %s\n", startTimeStr)
		gymStr += fmt.Sprintf("Time closed: %s\n", endTimeStr)

		builder.WriteString(gymStr)
		builder.WriteString("\n")
	}

	return builder.String()
}

func formatTrainersList(trainers []tmodels.Trainer) string {
	var builder strings.Builder

	for _, trainer := range trainers {
		trainerStr := fmt.Sprintf("Name: %s\n", trainer.Name)
		trainerStr += fmt.Sprintf("Speciality: %s\n", trainer.Speciality)

		builder.WriteString(trainerStr)
		builder.WriteString("\n")
	}

	return builder.String()
}

func formatABookingList(bookings []tmodels.Booking) string {
	var builder strings.Builder

	for _, booking := range bookings {
		startTimeStr := booking.StartTime.Format("2006.01.02 15:04")
		endTimeStr := booking.EndTime.Format("2006.01.02 15:04")
		bookingStr := fmt.Sprintf("ID: %d\n", booking.ID)
		bookingStr += fmt.Sprintf("Activity: %s\n", booking.Activity)
		bookingStr += fmt.Sprintf("Start time: %s\n", startTimeStr)
		bookingStr += fmt.Sprintf("End time: %s\n", endTimeStr)

		builder.WriteString(bookingStr)
		builder.WriteString("\n")
	}

	return builder.String()
}
