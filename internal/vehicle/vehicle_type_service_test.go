package vehicle

import (
    "testing"
)

func TestVehicleTypeService_CreateAndFind(t *testing.T) {
    repo := newMockRepo()
    svc := NewVehicleTypeService(repo)

    user := "test@example.com"
    vt, err := svc.Create("Sedan", true, user)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    if vt.ID == 0 {
        t.Fatalf("expected ID to be set")
    }
    if vt.Name != "Sedan" || vt.Status != true || vt.CreatedBy != user {
        t.Fatalf("unexpected vehicle fields: %+v", vt)
    }
    // FindByID
    fetched, err := svc.FindByID(vt.ID)
    if err != nil {
        t.Fatalf("FindByID error: %v", err)
    }
    if fetched == nil || fetched.Name != "Sedan" {
        t.Fatalf("FindByID returned wrong result: %+v", fetched)
    }
}

func TestVehicleTypeService_Update(t *testing.T) {
    repo := newMockRepo()
    svc := NewVehicleTypeService(repo)
    user := "creator@example.com"
    vt, _ := svc.Create("SUV", false, user)

    // Update fields
    updatedUser := "updater@example.com"
    updated, err := svc.Update(vt.ID, "Coupe", true, updatedUser)
    if err != nil {
        t.Fatalf("Update error: %v", err)
    }
    if updated.Name != "Coupe" || updated.Status != true {
        t.Fatalf("Update did not change fields: %+v", updated)
    }
    if updated.UpdatedBy == nil || *updated.UpdatedBy != updatedUser {
        t.Fatalf("UpdatedBy not set correctly: %+v", updated.UpdatedBy)
    }
    if updated.UpdatedAt == nil {
        t.Fatalf("UpdatedAt not set")
    }
}

func TestVehicleTypeService_ListPagedAndCount(t *testing.T) {
    repo := newMockRepo()
    svc := NewVehicleTypeService(repo)
    // create 3 items
    svc.Create("A", true, "a@example.com")
    svc.Create("B", false, "b@example.com")
    svc.Create("C", true, "c@example.com")

    list, total, err := svc.ListPaged(1, 2, "id", "asc")
    if err != nil {
        t.Fatalf("ListPaged error: %v", err)
    }
    if total != 3 {
        t.Fatalf("expected total 3, got %d", total)
    }
    if len(list) != 2 {
        t.Fatalf("expected 2 items, got %d", len(list))
    }
    // Ensure ordering by ID ascending (mock preserves insertion order, but we can just check IDs are increasing)
    if list[0].ID >= list[1].ID {
        t.Fatalf("expected ascending IDs, got %d then %d", list[0].ID, list[1].ID)
    }
}

func TestVehicleTypeService_Delete(t *testing.T) {
    repo := newMockRepo()
    svc := NewVehicleTypeService(repo)
    vt, _ := svc.Create("Truck", true, "t@example.com")
    if err := svc.Delete(vt.ID); err != nil {
        t.Fatalf("Delete error: %v", err)
    }
    fetched, _ := svc.FindByID(vt.ID)
    if fetched != nil {
        t.Fatalf("expected record to be deleted, got %+v", fetched)
    }
}

func TestVehicleTypeService_ListAll(t *testing.T) {
    repo := newMockRepo()
    svc := NewVehicleTypeService(repo)
    svc.Create("A", true, "a@example.com")
    svc.Create("B", false, "b@example.com")

    list, err := svc.ListAll()
    if err != nil {
        t.Fatalf("ListAll error: %v", err)
    }
    if len(list) != 2 {
        t.Fatalf("expected 2 items, got %d", len(list))
    }
}

func TestVehicleTypeService_Errors(t *testing.T) {
    t.Run("ListAllError", func(t *testing.T) {
        repo := newMockRepo()
        repo.getAllErr = true
        svc := NewVehicleTypeService(repo)
        _, err := svc.ListAll()
        if err == nil {
            t.Fatal("expected error from ListAll")
        }
    })

    t.Run("ListPagedDefaults", func(t *testing.T) {
        repo := newMockRepo()
        svc := NewVehicleTypeService(repo)
        svc.Create("X", true, "x@example.com")
        // page=0 and pageSize=0 should be normalised to 1 and 10
        list, total, err := svc.ListPaged(0, 0, "id", "asc")
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if total != 1 || len(list) != 1 {
            t.Fatalf("expected total 1 and 1 item, got %d/%d", total, len(list))
        }
    })

    t.Run("ListPagedGetPagedError", func(t *testing.T) {
        repo := newMockRepo()
        repo.getPagedErr = true
        svc := NewVehicleTypeService(repo)
        _, _, err := svc.ListPaged(1, 10, "id", "asc")
        if err == nil {
            t.Fatal("expected error from GetPaged")
        }
    })

    t.Run("ListPagedCountError", func(t *testing.T) {
        repo := newMockRepo()
        repo.countErr = true
        svc := NewVehicleTypeService(repo)
        _, _, err := svc.ListPaged(1, 10, "id", "asc")
        if err == nil {
            t.Fatal("expected error from Count")
        }
    })

    t.Run("CreateError", func(t *testing.T) {
        repo := newMockRepo()
        repo.createErr = true
        svc := NewVehicleTypeService(repo)
        _, err := svc.Create("Fail", true, "u@example.com")
        if err == nil {
            t.Fatal("expected error from Create")
        }
    })

    t.Run("UpdateGetByIDError", func(t *testing.T) {
        repo := newMockRepo()
        repo.getByIDErr = true
        svc := NewVehicleTypeService(repo)
        _, err := svc.Update(1, "Name", true, "u@example.com")
        if err == nil {
            t.Fatal("expected error from Update when GetByID fails")
        }
    })

    t.Run("UpdateNotFound", func(t *testing.T) {
        repo := newMockRepo()
        svc := NewVehicleTypeService(repo)
        result, err := svc.Update(999, "Name", true, "u@example.com")
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if result != nil {
            t.Fatal("expected nil result for non-existent record")
        }
    })

    t.Run("UpdateRepoError", func(t *testing.T) {
        repo := newMockRepo()
        repo.updateErr = true
        svc := NewVehicleTypeService(repo)
        svc.Create("Existing", true, "u@example.com")
        // ID 1 exists; repo.Update will fail
        _, err := svc.Update(1, "New", false, "u@example.com")
        if err == nil {
            t.Fatal("expected error from Update when repo.Update fails")
        }
    })
}
