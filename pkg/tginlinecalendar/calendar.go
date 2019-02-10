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

type CalendarView int

const (
	MonthView CalendarView = 0
	YearView  CalendarView = 1
)

//CalendarKeyboard holds internal calendar data and generates markup
type CalendarKeyboard struct {
	month int
	year  int
	view  CalendarView
}

const (
	CallbackNextMonth  = "tginlinecalendar-next-month"
	CallbackPrevMonth  = "tginlinecalendar-prev-month"
	CallbackNextYear   = "tginlinecalendar-next-year"
	CallbackPrevYear   = "tginlinecalendar-prev-year"
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

func (ck *CalendarKeyboard) NextYear() {
	ck.year += 1
}

func (ck *CalendarKeyboard) PrevYear() {
	ck.year -= 1
}

func (ck *CalendarKeyboard) SetViewMode(mode CalendarView) {
	ck.view = mode
}

func (ck *CalendarKeyboard) GetReplyMarkup() tgbotapi.InlineKeyboardMarkup {
	switch ck.view {
	case YearView:
		return ck.getYearReplyMarkup()
	default:
		return ck.getMonthReplyMarkup()
	}
}

func (ck *CalendarKeyboard) getMonthReplyMarkup() tgbotapi.InlineKeyboardMarkup {

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

	dayRowsCount := int(math.Ceil(float64(daysInMonth+startDayOfWeek) / 7))
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

func (ck *CalendarKeyboard) getYearReplyMarkup() tgbotapi.InlineKeyboardMarkup {

	date := time.Date(ck.year, 1, 1, 0, 0, 0, 0, time.UTC)

	yearRow := make([]tgbotapi.InlineKeyboardButton, 3)
	yearRow[0] = tgbotapi.NewInlineKeyboardButtonData("<", CallbackPrevYear)
	yearRow[1] = tgbotapi.NewInlineKeyboardButtonData(date.Format("2006"), CallbackEmpty)
	yearRow[2] = tgbotapi.NewInlineKeyboardButtonData(">", CallbackNextYear)

	monthRows := make([][]tgbotapi.InlineKeyboardButton, 4)

	for row := range monthRows {
		monthRows[row] = make([]tgbotapi.InlineKeyboardButton, 3)

		for col := range monthRows[row] {
			index := col + row*3
			monthDate := time.Date(ck.year, time.Month(index+1), 1, 0, 0, 0, 0, time.UTC)

			monthRows[row][col] = tgbotapi.NewInlineKeyboardButtonData(
				monthDate.Format("Jan 2006"),
				fmt.Sprintf("%s%d-%d-%d", CallbackDatePrefix, ck.year, index+1, 1))
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(append(monthRows, yearRow)...)
}

// daysIn returns the number of days in a month for a given year.
func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
