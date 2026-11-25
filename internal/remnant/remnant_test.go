package remnant

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/edanko/moses/internal/domain"
)

func profileType(t *testing.T) domain.ProfileDesignation {
	t.Helper()
	parsed, err := domain.ParseProfileDesignation("RP160X8")
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func TestLoadMissingFileIsEmpty(t *testing.T) {
	store, err := Load(filepath.Join(t.TempDir(), "missing.csv"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(store.Items()) != 0 {
		t.Fatalf("items = %d, want 0", len(store.Items()))
	}
}

func TestSaveLoadAddAndUse(t *testing.T) {
	path := filepath.Join(t.TempDir(), "remnants.csv")
	store, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	added := store.Add(domain.Remnant{
		Length:      5000,
		Project:     "project",
		Source:      "section",
		Designation: profileType(t),
		Grade:       "A40",
		Marking:     "nest-L5000",
	})
	if added.ID != 1 {
		t.Fatalf("id = %d, want 1", added.ID)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), header+"\n") {
		t.Fatalf("missing header: %q", string(data))
	}

	reloaded, err := Load(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	available := reloaded.Available("project", profileType(t), "A40")
	if len(available) != 1 {
		t.Fatalf("available = %d, want 1", len(available))
	}
	if available[0].Marking != "nest-L5000" || available[0].Length != 5000 {
		t.Fatalf("unexpected remnant: %+v", available[0])
	}

	reloaded.MarkUsed(1, "nest-1")
	if available := reloaded.Available("project", profileType(t), "A40"); len(available) != 0 {
		t.Fatalf("used remnant is still available: %+v", available)
	}
	if err := reloaded.Save(); err != nil {
		t.Fatal(err)
	}

	again, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	items := again.Items()
	if len(items) != 1 || !items[0].Used || items[0].ConsumedBy != "nest-1" {
		t.Fatalf("unexpected saved item: %+v", items)
	}
}

func TestAvailableFilters(t *testing.T) {
	store, err := Load(filepath.Join(t.TempDir(), "remnants.csv"))
	if err != nil {
		t.Fatal(err)
	}
	other, err := domain.ParseProfileDesignation("FB80x8")
	if err != nil {
		t.Fatal(err)
	}
	store.Add(domain.Remnant{Length: 1000, Project: "a", Source: "s", Designation: profileType(t), Grade: "A40", Marking: "m1"})
	store.Add(domain.Remnant{Length: 1000, Project: "b", Source: "s", Designation: profileType(t), Grade: "A40", Marking: "m2"})
	store.Add(domain.Remnant{Length: 1000, Project: "a", Source: "s", Designation: other, Grade: "A40", Marking: "m3"})

	available := store.Available("a", profileType(t), "A40")
	if len(available) != 1 || available[0].Marking != "m1" {
		t.Fatalf("unexpected available: %+v", available)
	}
}
