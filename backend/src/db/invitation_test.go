package db

import (
	"errors"
	"forum/src/ferror"
	"testing"
)

func setupFollowTestDB(t *testing.T) (int64, int64, int64) {
	t.Helper()
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("InitDB(%q) failed: %v", ":memory:", err)
	}
	var ids []int64
	for _, name := range []string{"alice", "bob", "carol"} {
		user := UserRowType{}
		user.Username = name
		user.Email = name + "@test.com"
		user.Hash = "x"
		if err := user.InsertUserWithHash(); err != nil {
			t.Fatalf("create user %s: %v", name, err)
		}
		ids = append(ids, user.Id)
	}
	return ids[0], ids[1], ids[2]
}

func insertInvitation(from, to int64) (int64, error) {
	inv := &InvitationRowType{}
	inv.FromUserId = from
	inv.ToUserId = to
	if err := inv.Insert(); err != nil {
		return 0, err
	}
	return inv.Id, nil
}

func TestInvitationSelectFollowersAndFollowing(t *testing.T) {
	alice, bob, carol := setupFollowTestDB(t)

	bobToAlice, err := insertInvitation(bob, alice)
	if err != nil {
		t.Fatalf("insert bob->alice: %v", err)
	}
	carolToAlice, err := insertInvitation(carol, alice)
	if err != nil {
		t.Fatalf("insert carol->alice: %v", err)
	}
	// alice follows bob; bob also follows alice (two-way accepted)
	aliceToBob, err := insertInvitation(alice, bob)
	if err != nil {
		t.Fatalf("insert alice->bob: %v", err)
	}

	// These users have no profiles, so invitations are created as 'pending'.
	// Accept them all to test the accepted-based queries.
	for _, id := range []int64{bobToAlice, carolToAlice, aliceToBob} {
		inv := InvitationRowType{}
		inv.Id = id
		inv.Status = "accepted"
		if err := inv.UpdateInvitationStatus(); err != nil {
			t.Fatalf("update status: %v", err)
		}
	}

	followers, err := SelectFollowersByUserId(alice)
	if err != nil {
		t.Fatalf("SelectFollowersByUserId() error: %v", err)
	}
	if len(followers) != 2 {
		t.Errorf("SelectFollowersByUserId(alice) len = %d, want 2", len(followers))
	}
	usernames := map[string]bool{}
	for _, f := range followers {
		usernames[f.Username] = true
		if f.Status != "accepted" || f.ToUserId != alice {
			t.Errorf("follower = %+v, want accepted invite to alice", f)
		}
	}
	if !usernames["bob"] || !usernames["carol"] {
		t.Errorf("followers = %v, want bob and carol", usernames)
	}

	following, err := SelectFollowingByUserId(alice)
	if err != nil {
		t.Fatalf("SelectFollowingByUserId() error: %v", err)
	}
	if len(following) != 1 || following[0].Username != "bob" || following[0].ToUserId != bob {
		t.Errorf("SelectFollowingByUserId(alice) = %+v, want 1 follow of bob", following)
	}
}

func TestInvitationSelectPendingByUserId(t *testing.T) {
	_, bob, carol := setupFollowTestDB(t)

	if _, err := insertInvitation(bob, carol); err != nil {
		t.Fatalf("insert bob->carol: %v", err)
	}
	if _, err := insertInvitation(carol, bob); err != nil {
		t.Fatalf("insert carol->bob: %v", err)
	}

	pending, err := SelectPendingFollowRequestsByUserId(bob)
	if err != nil {
		t.Fatalf("SelectPendingFollowRequestsByUserId() error: %v", err)
	}
	if len(pending) != 2 {
		t.Fatalf("len(pending) = %d, want 2", len(pending))
	}
	fromBob, toBob := 0, 0
	for _, p := range pending {
		if p.Status != "pending" {
			t.Errorf("received status %q, want pending", p.Status)
		}
		if p.FromUserId == bob {
			fromBob++
			if p.ToUserId != carol || p.Username != "carol" {
				t.Errorf("outgoing = %+v, want to carol", p)
			}
		}
		if p.ToUserId == bob {
			toBob++
			if p.FromUserId != carol || p.Username != "carol" {
				t.Errorf("incoming = %+v, want from carol", p)
			}
		}
	}
	if fromBob != 1 || toBob != 1 {
		t.Errorf("fromBob = %d, toBob = %d, want 1 each", fromBob, toBob)
	}
}

func TestInvitationSelectByUsersAndNotFound(t *testing.T) {
	alice, bob, carol := setupFollowTestDB(t)

	id, err := insertInvitation(bob, alice)
	if err != nil {
		t.Fatalf("insert bob->alice: %v", err)
	}

	inv, err := SelectInvitationByUsers(bob, alice)
	if err != nil {
		t.Fatalf("SelectInvitationByUsers() error: %v", err)
	}
	if inv.Id != id || inv.FromUserId != bob || inv.ToUserId != alice {
		t.Errorf("SelectInvitationByUsers() = %+v, want invitation %d bob->alice", inv, id)
	}

	if _, err := SelectInvitationByUsers(bob, carol); err == nil {
		t.Error("SelectInvitationByUsers() should return error when invitation does not exist")
	} else if !errors.Is(err, ferror.ErrorInvitationNotFound) {
		t.Errorf("SelectInvitationByUsers() error = %v, want %v", err, ferror.ErrorInvitationNotFound)
	}
}

func TestInvitationReactivateAfterDeny(t *testing.T) {
	alice, bob, _ := setupFollowTestDB(t)

	id, err := insertInvitation(bob, alice)
	if err != nil {
		t.Fatalf("insert bob->alice: %v", err)
	}

	inv := InvitationRowType{}
	inv.Id = id
	if err := inv.SelectInvitationById(); err != nil {
		t.Fatalf("SelectInvitationById() error: %v", err)
	}
	inv.Status = "denied"
	if err := inv.UpdateInvitationStatus(); err != nil {
		t.Fatalf("deny invitation: %v", err)
	}

	// A new request between the same users reactivates instead of erroring.
	newInv := InvitationRowType{}
	newInv.FromUserId = bob
	newInv.ToUserId = alice
	if err := newInv.Insert(); err != nil {
		t.Fatalf("re-follow after deny: %v", err)
	}
	reloaded := InvitationRowType{}
	reloaded.Id = id
	if err := reloaded.SelectInvitationById(); err != nil {
		t.Fatalf("SelectInvitationById() (reload) error: %v", err)
	}
	if reloaded.Status != "pending" {
		t.Errorf("status after reactivation = %q, want pending", reloaded.Status)
	}

	// Re-following again while pending must fail.
	dup := InvitationRowType{}
	dup.FromUserId = bob
	dup.ToUserId = alice
	if err := dup.Insert(); err == nil {
		t.Error("Insert() should fail when invitation is already pending")
	} else if !errors.Is(err, ferror.ErrorInvitationAlreadyExists) {
		t.Errorf("Insert() error = %v, want %v", err, ferror.ErrorInvitationAlreadyExists)
	}
}

func TestInvitationSelectByIdAndUpdateStatus(t *testing.T) {
	alice, bob, _ := setupFollowTestDB(t)

	id, err := insertInvitation(bob, alice)
	if err != nil {
		t.Fatalf("insert bob->alice: %v", err)
	}

	inv := InvitationRowType{}
	inv.Id = id
	if err := inv.SelectInvitationById(); err != nil {
		t.Fatalf("SelectInvitationById() error: %v", err)
	}
	if inv.Status != "pending" || inv.FromUserId != bob || inv.ToUserId != alice {
		t.Errorf("SelectInvitationById() = %+v, want pending bob->alice", inv)
	}

	inv.Status = "denied"
	if err := inv.UpdateInvitationStatus(); err != nil {
		t.Fatalf("UpdateInvitationStatus() error: %v", err)
	}
	var reloaded InvitationRowType
	reloaded.Id = id
	if err := reloaded.SelectInvitationById(); err != nil {
		t.Fatalf("SelectInvitationById() (reload) error: %v", err)
	}
	if reloaded.Status != "denied" {
		t.Errorf("status after update = %q, want denied", reloaded.Status)
	}

	missing := InvitationRowType{}
	missing.Id = id + 500
	if err := missing.SelectInvitationById(); err == nil {
		t.Error("SelectInvitationById() should return error for nonexistent invitation")
	} else if !errors.Is(err, ferror.ErrorInvitationNotFound) {
		t.Errorf("SelectInvitationById() error = %v, want %v", err, ferror.ErrorInvitationNotFound)
	}
}
