package worship

import (
	"context"
	"errors"
	"testing"

	"github.com/pisondev/mosque-api/internal/testutil"
)

func TestRepositoryIntegrationWorship(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	tenant := testutil.SeedTenant(t, db)
	defer testutil.CleanupTenant(t, db, tenant.ID)

	repo := NewRepository(db)
	ctx := context.Background()

	setReq := PrayerTimeSettingsRequest{
		Timezone:      "Asia/Jakarta",
		LocationMode:  "city",
		CityName:      ptrString("Bandung"),
		AdjSubuhMin:   1,
		AdjDzuhurMin:  0,
		AdjAsharMin:   1,
		AdjMaghribMin: 0,
		AdjIsyaMin:    0,
	}
	setRes, err := repo.UpsertPrayerTimeSettings(ctx, tenant.ID, setReq)
	if err != nil {
		t.Fatalf("UpsertPrayerTimeSettings failed: %v", err)
	}
	if setRes.LocationMode != "city" {
		t.Fatalf("unexpected location mode: %s", setRes.LocationMode)
	}

	gotSet, err := repo.GetPrayerTimeSettings(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("GetPrayerTimeSettings failed: %v", err)
	}
	if gotSet.Timezone != "Asia/Jakarta" {
		t.Fatalf("unexpected timezone: %s", gotSet.Timezone)
	}

	dailyReq := PrayerTimesDailyPayload{
		DayDate:     "2030-01-02",
		SubuhTime:   "04:30:00",
		DzuhurTime:  "11:59:00",
		AsharTime:   "15:10:00",
		MaghribTime: "18:01:00",
		IsyaTime:    "19:11:00",
		SourceLabel: ptrString("integration"),
	}
	dailyRes, err := repo.CreatePrayerTimesDaily(ctx, tenant.ID, dailyReq)
	if err != nil {
		t.Fatalf("CreatePrayerTimesDaily failed: %v", err)
	}
	if dailyRes.ID == 0 {
		t.Fatal("CreatePrayerTimesDaily returned empty id")
	}

	dailyList, totalDaily, err := repo.ListPrayerTimesDaily(ctx, tenant.ID, PrayerTimesDailyQuery{
		From:  "2030-01-01",
		To:    "2030-01-10",
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListPrayerTimesDaily failed: %v", err)
	}
	if totalDaily < 1 || len(dailyList) < 1 {
		t.Fatalf("ListPrayerTimesDaily expected data, total=%d len=%d", totalDaily, len(dailyList))
	}

	gotDaily, err := repo.GetPrayerTimesDaily(ctx, tenant.ID, dailyRes.ID)
	if err != nil {
		t.Fatalf("GetPrayerTimesDaily failed: %v", err)
	}
	if gotDaily.DayDate != "2030-01-02" {
		t.Fatalf("unexpected day date: %s", gotDaily.DayDate)
	}

	updateDaily := dailyReq
	updateDaily.SourceLabel = ptrString("updated")
	if err := repo.UpdatePrayerTimesDaily(ctx, tenant.ID, dailyRes.ID, updateDaily); err != nil {
		t.Fatalf("UpdatePrayerTimesDaily failed: %v", err)
	}

	if err := repo.DeletePrayerTimesDaily(ctx, tenant.ID, dailyRes.ID); err != nil {
		t.Fatalf("DeletePrayerTimesDaily failed: %v", err)
	}
	_, err = repo.GetPrayerTimesDaily(ctx, tenant.ID, dailyRes.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetPrayerTimesDaily after delete expected ErrNotFound, got %v", err)
	}

	dutyRes, err := repo.CreatePrayerDuty(ctx, tenant.ID, PrayerDutyPayload{
		Category:    "fardhu",
		DutyDate:    "2030-01-03",
		Prayer:      ptrString("maghrib"),
		ImamName:    ptrString("Imam A"),
		MuadzinName: ptrString("Muadzin A"),
	})
	if err != nil {
		t.Fatalf("CreatePrayerDuty failed: %v", err)
	}
	if dutyRes.ID == 0 {
		t.Fatal("CreatePrayerDuty returned empty id")
	}

	specialRes, err := repo.CreateSpecialDay(ctx, tenant.ID, SpecialDayPayload{
		Kind:    "other",
		Title:   "Hari Khusus Integrasi",
		DayDate: "2030-01-04",
	})
	if err != nil {
		t.Fatalf("CreateSpecialDay failed: %v", err)
	}
	if specialRes.ID == 0 {
		t.Fatal("CreateSpecialDay returned empty id")
	}

	cal, err := repo.GetPrayerCalendar(ctx, tenant.ID, "2030-01-01", "2030-01-10")
	if err != nil {
		t.Fatalf("GetPrayerCalendar failed: %v", err)
	}
	if _, ok := cal["prayer_duties"]; !ok {
		t.Fatal("GetPrayerCalendar missing prayer_duties")
	}
	if _, ok := cal["special_days"]; !ok {
		t.Fatal("GetPrayerCalendar missing special_days")
	}
}

func ptrString(v string) *string {
	return &v
}
