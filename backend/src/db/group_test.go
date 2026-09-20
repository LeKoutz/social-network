package db

import (
	"forum/src/utils"
	"testing"
)

func TestGroupRowIsMember(t *testing.T) {
	if err := InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}

	owner := UserRowType{}
	owner.Username = "gowner"
	owner.Email = "gowner@test.com"
	owner.Hash = "x"
	if err := owner.InsertUserWithHash(); err != nil {
		t.Fatalf("create owner: %v", err)
	}
	member := UserRowType{}
	member.Username = "gmember"
	member.Email = "gmember@test.com"
	member.Hash = "x"
	if err := member.InsertUserWithHash(); err != nil {
		t.Fatalf("create member: %v", err)
	}
	stranger := UserRowType{}
	stranger.Username = "gstranger"
	stranger.Email = "gstranger@test.com"
	stranger.Hash = "x"
	if err := stranger.InsertUserWithHash(); err != nil {
		t.Fatalf("create stranger: %v", err)
	}

	group := GroupRowType{}
	group.Title = "Closed Club"
	group.OwnerUserId = owner.Id
	if err := group.InsertGroup(); err != nil {
		t.Fatalf("create group: %v", err)
	}

	ok, err := group.IsMember(owner.Id)
	if err != nil || !ok {
		t.Fatalf("IsMember(owner) = %v, %v; want true, nil", ok, err)
	}

	_, err = db.Exec(
		`INSERT INTO group_invitations (timestamp, status, group_id, to_user_id, invited_by)
		VALUES (?, 'accepted', ?, ?, ?)`,
		utils.GetCurrentTimestamp(),
		group.Id,
		member.Id,
		owner.Id,
	)
	if err != nil {
		t.Fatalf("create accepted invitation: %v", err)
	}
	ok, err = group.IsMember(member.Id)
	if err != nil || !ok {
		t.Fatalf("IsMember(accepted invitee) = %v, %v; want true, nil", ok, err)
	}

	ok, err = group.IsMember(stranger.Id)
	if err != nil || ok {
		t.Fatalf("IsMember(stranger) = %v, %v; want false, nil", ok, err)
	}
}