//tginlinecalendar is a helper package for telergam bots to create calendar in inline
//keyboard using telergam-bot-api framework
package tginlinecalendar

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

//CalendarKeyboard holds internal calendar data and generates markup
type CalendarKeyboard struct {
	month int
	year  int
}

const (
	CallbackNextMonth  = "tginlinecalendar-next-month"
	CallbackPrevMonth  = "tginlinecalendar-prev-month"
	CallbackEmpty      = "tginlinecalendar-emtpy"
	CallbackDatePrefix = "tginlinecalendar-date-"
)

func NewCalendarKeyboard(month int, year int) *CalendarKeyboard {
	return &CalendarKeyboard{month: month, year: year}
}

func ExtractDate(queryData string) (time.Time, error) {
	if !strings.HasPrefix(queryData, CallbackDatePrefix) {
		return time.Now(), &time.ParseError{Layout: CallbackDatePrefix, Value: queryData}
	}
	dateString := queryData[len(CallbackDatePrefix):]
	return time.Parse("2006-1-2", dateString)
}

func (ck *CalendarKeyboard) NextMonth() {
	ck.month += 1
	if ck.month > 12 {
		ck.month = 1
		ck.year += 1
	}
}

func (ck *CalendarKeyboard) PrevMonth() {
	ck.month -= 1
	if ck.month < 1 {
		ck.month = 12
		ck.year -= 1
	}
}

func (ck *CalendarKeyboard) GetReplyMarkup() tgbotapi.InlineKeyboardMarkup {

	date := time.Date(ck.year, time.Month(ck.month), 1, 0, 0, 0, 0, time.UTC)

	fmt.Println(ck.month, date.Month())

	monthRow := make([]tgbotapi.InlineKeyboardButton, 3)
	monthRow[0] = tgbotapi.NewInlineKeyboardButtonData("<", CallbackPrevMonth)
	monthRow[1] = tgbotapi.NewInlineKeyboardButtonData(date.Format("Jan 2006"), CallbackEmpty)
	monthRow[2] = tgbotapi.NewInlineKeyboardButtonData(">", CallbackNextMonth)

	fmt.Println(date.Weekday())

	startDayOfWeek := int(date.Weekday()) - 1
	if startDayOfWeek < 0 {
		startDayOfWeek = 6
	}

	daysInMonth := daysIn(date.Month(), date.Year())

	fmt.Println(daysInMonth)

	dayRowsCount := int(math.Ceil(float64(daysInMonth) / 7))
	dayRows := make([][]tgbotapi.InlineKeyboardButton, dayRowsCount+1)

	for row := range dayRows {
		dayRows[row] = make([]tgbotapi.InlineKeyboardButton, 7)
		if row == 0 {
			dayRows[row][0] = tgbotapi.NewInlineKeyboardButtonData("П", CallbackEmpty)
			dayRows[row][1] = tgbotapi.NewInlineKeyboardButtonData("В", CallbackEmpty)
			dayRows[row][2] = tgbotapi.NewInlineKeyboardButtonData("С", CallbackEmpty)
			dayRows[row][3] = tgbotapi.NewInlineKeyboardButtonData("Ч", CallbackEmpty)
			dayRows[row][4] = tgbotapi.NewInlineKeyboardButtonData("П", CallbackEmpty)
			dayRows[row][5] = tgbotapi.NewInlineKeyboardButtonData("С", CallbackEmpty)
			dayRows[row][6] = tgbotapi.NewInlineKeyboardButtonData("В", CallbackEmpty)

			continue
		}

		for col := range dayRows[row] {
			index := col + (row-1)*7
			if index >= startDayOfWeek && index < daysInMonth+startDayOfWeek {
				day := index - startDayOfWeek + 1
				dayRows[row][col] = tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprint(day), fmt.Sprintf("%s%d-%d-%d", CallbackDatePrefix, ck.year, ck.month, day))
			} else {
				dayRows[row][col] = tgbotapi.NewInlineKeyboardButtonData(" ", CallbackEmpty)
			}
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(append(dayRows, monthRow)...)
}

// daysIn returns the number of days in a month for a given year.
func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
