package types

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestTime_MarshalJSON(t *testing.T) {
	now := time.Now()
	nowJson, err := json.Marshal(now)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nowJson:", string(nowJson[:]))

	tim := DateTime(now)
	timJson, err := json.Marshal(tim)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("timJson:", string(timJson[:]))

	dat := Date(now)
	datJson, err := json.Marshal(dat)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("datJson:", string(datJson[:]))
}

func TestTime_GetDays(t *testing.T) {
	now := time.Now()

	today := DateTime(now)
	days := today.GetDays(now)
	if days != 0 {
		log.Fatal("expect 0; actual ", days)
	}

	yesterday := DateTime(now.Add(-time.Hour * 24))
	days = yesterday.GetDays(now)
	if days != -1 {
		log.Fatal("expect -1; actual ", days)
	}

	tomorrow := DateTime(now.Add(time.Hour * 24))
	days = tomorrow.GetDays(now)
	if days != 1 {
		log.Fatal("expect 1; actual ", days)
	}

	after := DateTime(now.Add(time.Hour * 24 * 7))
	days = after.GetDays(now)
	if days != 7 {
		log.Fatal("expect 7; actual ", days)
	}

	before := DateTime(now.Add(-time.Hour * 24 * 5))
	days = before.GetDays(now)
	if days != -5 {
		log.Fatal("expect -5; actual ", days)
	}
}
