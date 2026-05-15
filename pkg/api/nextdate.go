package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату задачи согласно правилу повторения.
//
// now    — точка отсчёта; результат должен быть строго позже неё
// dstart — исходная дата задачи в формате 20060102
// repeat — правило повторения (d N | y | w 1-7,… | m 1-31/-1/-2,… [1-12,…])
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart %q: %w", dstart, err)
	}

	// Нормализуем now до начала дня, чтобы сравнения были по дате, а не времени
	now = truncateToDay(now)

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		return nextDateD(now, start, parts)
	case "y":
		return nextDateY(now, start)
	case "w":
		return nextDateW(now, start, parts)
	case "m":
		return nextDateM(now, start, parts)
	default:
		return "", fmt.Errorf("unknown repeat rule %q", parts[0])
	}
}

// nextDateHandler обрабатывает GET /api/nextdate?now=<20060102>&date=<20060102>&repeat=<rule>
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: fmt.Sprintf("invalid now: %v", err)})
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, next)
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// afterBoth — cur строго позже now И строго позже start.
// результат должен быть в будущем относительно now,
// отсчёт повторений идёт от start, поэтому cur должен превышать оба значения.
func afterBoth(cur, now, start time.Time) bool {
	return cur.After(now) && cur.After(start)
}

// nextDateD — правило "d N": каждые N дней (1 ≤ N ≤ 400)
func nextDateD(now, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("d: interval not specified")
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil || n < 1 || n > 400 {
		return "", fmt.Errorf("d: invalid interval %q (must be 1..400)", parts[1])
	}
	cur := start
	for !afterBoth(cur, now, start) {
		cur = cur.AddDate(0, 0, n)
	}
	return cur.Format(DateFormat), nil
}

// nextDateY — правило "y": раз в год (тот же день/месяц)
func nextDateY(now, start time.Time) (string, error) {
	cur := start
	for !afterBoth(cur, now, start) {
		cur = cur.AddDate(1, 0, 0)
	}
	return cur.Format(DateFormat), nil
}

// nextDateW — правило "w 1,2,…":  дни недели (1=пн  7=вс)
func nextDateW(now, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("w: weekdays not specified")
	}

	var weekday [8]bool // индексы 1..7
	for _, tok := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(strings.TrimSpace(tok))
		if err != nil || d < 1 || d > 7 {
			return "", fmt.Errorf("w: invalid weekday %q (must be 1..7)", tok)
		}
		weekday[d] = true
	}

	// Sunday=0, Monday=1  Saturday=6
	// нужно: Monday=1  Sunday=7
	goToOur := func(wd time.Weekday) int {
		if wd == time.Sunday {
			return 7
		}
		return int(wd)
	}

	base := now
	if start.After(now) {
		base = start
	}
	cur := base.AddDate(0, 0, 1)

	for i := 0; i < 8; i++ { // максимум 7 итераций до нужного дня
		if weekday[goToOur(cur.Weekday())] {
			return cur.Format(DateFormat), nil
		}
		cur = cur.AddDate(0, 0, 1)
	}

	return "", errors.New("w: no matching weekday found (internal error)")
}

// nextDateM — правило "m days [months]": указанные дни (и месяцы)
//
// Дни: 1..31, -1 (последний), -2 (предпоследний)
// Месяцы: 1..12
func nextDateM(now, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("m: days not specified")
	}

	//разбираем дни
	// Разделяем на «дни» и, если есть, «месяцы» — они идут вторым токеном через пробел
	dayTokens := strings.Split(parts[1], ",")
	var monthTokens []string
	if len(parts) >= 3 {
		monthTokens = strings.Split(parts[2], ",")
	}

	// Флаги дней: индексы 1..31 для обычных, 32 = -2, 33 = -1
	const idxLastMinus2 = 32
	const idxLast = 33
	var day [34]bool

	for _, tok := range dayTokens {
		tok = strings.TrimSpace(tok)
		d, err := strconv.Atoi(tok)
		if err != nil {
			return "", fmt.Errorf("m: invalid day %q", tok)
		}
		switch {
		case d >= 1 && d <= 31:
			day[d] = true
		case d == -1:
			day[idxLast] = true
		case d == -2:
			day[idxLastMinus2] = true
		default:
			return "", fmt.Errorf("m: day %d out of range (1..31, -1, -2)", d)
		}
	}

	// Флаги месяцев: индексы 1..12; пусто = все месяцы
	var month [13]bool
	allMonths := len(monthTokens) == 0
	if !allMonths {
		for _, tok := range monthTokens {
			tok = strings.TrimSpace(tok)
			mo, err := strconv.Atoi(tok)
			if err != nil || mo < 1 || mo > 12 {
				return "", fmt.Errorf("m: invalid month %q (must be 1..12)", tok)
			}
			month[mo] = true
		}
	}

	// Начинаем перебор с дня, следующего за max
	base := now
	if start.After(now) {
		base = start
	}
	cur := base.AddDate(0, 0, 1)

	// Ищем подходящий день (не более ~800 итераций +-= 2 года)
	for i := 0; i < 800; i++ {
		mo := int(cur.Month())
		if allMonths || month[mo] {
			if matchDay(cur, day, idxLast, idxLastMinus2) {
				return cur.Format(DateFormat), nil
			}
		}
		cur = cur.AddDate(0, 0, 1)
	}

	return "", errors.New("m: no matching date found within 2 years")
}

// matchDay проверяет, подходит ли дата cur под маску дней day.
func matchDay(cur time.Time, day [34]bool, idxLast, idxLastMinus2 int) bool {
	d := cur.Day()

	// Обычный день
	if d <= 31 && day[d] {
		return true
	}

	// Последний день месяца
	last := lastDayOfMonth(cur)
	if day[idxLast] && d == last {
		return true
	}

	// Предпоследний день месяца
	if day[idxLastMinus2] && d == last-1 {
		return true
	}

	return false
}

// lastDayOfMonth возвращает число последнего дня месяца для даты t.
func lastDayOfMonth(t time.Time) int {
	// Первый день следующего месяца минус один день
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
