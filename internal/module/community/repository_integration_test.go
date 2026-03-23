package community

import (
	"context"
	"testing"

	"github.com/pisondev/mosque-api/internal/testutil"
)

func TestRepositoryIntegrationCommunity(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	tenant := testutil.SeedTenant(t, db)
	defer testutil.CleanupTenant(t, db, tenant.ID)

	repo := NewRepository(db)
	ctx := context.Background()

	eventReq := EventPayload{
		Title:          "Kajian Integrasi",
		Category:       "kajian_rutin",
		StartDate:      "2031-01-02",
		TimeMode:       "exact_time",
		StartTime:      ptrStr("07:00:00"),
		Status:         "upcoming",
		FeeAmount:      ptrStr("0"),
		RepeatWeekdays: []int16{},
	}
	eventRes, err := repo.CreateEvent(ctx, tenant.ID, eventReq)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	if eventRes.ID == 0 {
		t.Fatal("CreateEvent returned empty id")
	}

	events, totalEvents, err := repo.ListEvents(ctx, tenant.ID, EventListQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if totalEvents < 1 || len(events) < 1 {
		t.Fatalf("ListEvents expected data, total=%d len=%d", totalEvents, len(events))
	}

	gotEvent, err := repo.GetEvent(ctx, tenant.ID, eventRes.ID)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}
	if gotEvent.Title != "Kajian Integrasi" {
		t.Fatalf("unexpected event title: %s", gotEvent.Title)
	}

	eventReq.Title = "Kajian Integrasi Updated"
	if err := repo.UpdateEvent(ctx, tenant.ID, eventRes.ID, eventReq); err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}
	if err := repo.UpdateEventStatus(ctx, tenant.ID, eventRes.ID, "upcoming"); err != nil {
		t.Fatalf("UpdateEventStatus failed: %v", err)
	}

	albumRes, err := repo.CreateGalleryAlbum(ctx, tenant.ID, GalleryAlbumPayload{
		Title:     "Album Integrasi",
		MediaKind: "photo",
	})
	if err != nil {
		t.Fatalf("CreateGalleryAlbum failed: %v", err)
	}
	if albumRes.ID == 0 {
		t.Fatal("CreateGalleryAlbum returned empty id")
	}

	itemRes, err := repo.CreateGalleryItem(ctx, tenant.ID, GalleryItemPayload{
		AlbumID:     &albumRes.ID,
		MediaType:   "image",
		MediaURL:    "https://cdn.example.com/integration.jpg",
		SortOrder:   1,
		IsHighlight: true,
	})
	if err != nil {
		t.Fatalf("CreateGalleryItem failed: %v", err)
	}
	if itemRes.ID == 0 {
		t.Fatal("CreateGalleryItem returned empty id")
	}

	mmRes, err := repo.CreateManagementMember(ctx, tenant.ID, ManagementMemberPayload{
		FullName:   "Pengurus Integrasi",
		RoleTitle:  "Ketua",
		ShowPublic: true,
		SortOrder:  1,
	})
	if err != nil {
		t.Fatalf("CreateManagementMember failed: %v", err)
	}
	if mmRes.ID == 0 {
		t.Fatal("CreateManagementMember returned empty id")
	}

	publicEvents, _, err := repo.ListPublicEvents(ctx, tenant.Hostname, 1, 10)
	if err != nil {
		t.Fatalf("ListPublicEvents failed: %v", err)
	}
	if len(publicEvents) < 1 {
		t.Fatal("ListPublicEvents expected data")
	}

	publicAlbums, _, err := repo.ListPublicGalleryAlbums(ctx, tenant.Hostname, 1, 10)
	if err != nil {
		t.Fatalf("ListPublicGalleryAlbums failed: %v", err)
	}
	if len(publicAlbums) < 1 {
		t.Fatal("ListPublicGalleryAlbums expected data")
	}

	publicItems, _, err := repo.ListPublicGalleryItems(ctx, tenant.Hostname, 1, 10)
	if err != nil {
		t.Fatalf("ListPublicGalleryItems failed: %v", err)
	}
	if len(publicItems) < 1 {
		t.Fatal("ListPublicGalleryItems expected data")
	}

	publicMembers, _, err := repo.ListPublicManagementMembers(ctx, tenant.Hostname, 1, 10)
	if err != nil {
		t.Fatalf("ListPublicManagementMembers failed: %v", err)
	}
	if len(publicMembers) < 1 {
		t.Fatal("ListPublicManagementMembers expected data")
	}

	if err := repo.DeleteGalleryItem(ctx, tenant.ID, itemRes.ID); err != nil {
		t.Fatalf("DeleteGalleryItem failed: %v", err)
	}
	if err := repo.DeleteGalleryAlbum(ctx, tenant.ID, albumRes.ID); err != nil {
		t.Fatalf("DeleteGalleryAlbum failed: %v", err)
	}
	if err := repo.DeleteManagementMember(ctx, tenant.ID, mmRes.ID); err != nil {
		t.Fatalf("DeleteManagementMember failed: %v", err)
	}
	if err := repo.DeleteEvent(ctx, tenant.ID, eventRes.ID); err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}
}

func ptrStr(v string) *string {
	return &v
}
