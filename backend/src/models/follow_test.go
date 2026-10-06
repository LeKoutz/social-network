package models

import (
	"errors"
	"forum/src/db"
	"forum/src/ferror"
	"testing"
)

func setupFollowModelTestDB(t *testing.T) (int64, int64, int64) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("InitDB(%q) failed: %v", ":memory:", err)
	}
	create := func(name string) int64 {
		user := db.UserRowType{}
		user.Username = name
		user.Email = name + "@test.com"
		user.Hash = "x"
		if err := user.InsertUserWithHash(); err != nil {
			t.Fatalf("create user %s: %v", name, err)
		}
		return user.Id
	}
	return create("alice"), create("bob"), create("carol")
}

// requestFollow creates a follow request. Users created without a profile
// default to 'pending' status.
func requestFollow(from, to int64) (FollowRequestType, error) {
	var rq FollowRequestType
	rq.FromUserId = from
	rq.ToUserId = to
	if err := rq.AddFollowInvite(); err != nil {
		return rq, err
	}
	return reload(rq.Id)
}

func acceptByRecipient(rq FollowRequestType, recipient int64) error {
	var acceptor FollowRequestType
	acceptor.InvitationRowType.Id = rq.Id
	return acceptor.Accept(recipient)
}

func reload(rqId int64) (FollowRequestType, error) {
	var r FollowRequestType
	r.InvitationRowType.Id = rqId
	err := r.InvitationRowType.SelectInvitationById()
	return r, err
}

func TestFollowRequestAccept(t *testing.T) {
	alice, bob, _ := setupFollowModelTestDB(t)

	rq, err := requestFollow(bob, alice)
	if err != nil {
		t.Fatalf("request follow: %v", err)
	}
	if rq.Status != FollowStatusPending {
		t.Fatalf("status = %q, want %q", rq.Status, FollowStatusPending)
	}

	if err := acceptByRecipient(rq, alice); err != nil {
		t.Fatalf("Accept() error: %v", err)
	}
	accepted, err := reload(rq.Id)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if accepted.Status != FollowStatusAccepted {
		t.Errorf("status after accept = %q, want %q", accepted.Status, FollowStatusAccepted)
	}
}

func TestFollowRequestDeny(t *testing.T) {
	alice, bob, _ := setupFollowModelTestDB(t)

	rq, err := requestFollow(bob, alice)
	if err != nil {
		t.Fatalf("request follow: %v", err)
	}
	if rq.Status != FollowStatusPending {
		t.Fatalf("status = %q, want %q", rq.Status, FollowStatusPending)
	}

	var denier FollowRequestType
	denier.InvitationRowType.Id = rq.Id
	if err := denier.Deny(alice); err != nil {
		t.Fatalf("Deny() error: %v", err)
	}
	denied, err := reload(rq.Id)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if denied.Status != FollowStatusDenied {
		t.Errorf("status after deny = %q, want %q", denied.Status, FollowStatusDenied)
	}
}

func TestFollowRequestOnlyRecipientCanRespond(t *testing.T) {
	alice, bob, carol := setupFollowModelTestDB(t)

	rq, err := requestFollow(bob, carol)
	if err != nil {
		t.Fatalf("request follow: %v", err)
	}
	if rq.Status != FollowStatusPending {
		t.Fatalf("status = %q, want %q", rq.Status, FollowStatusPending)
	}

	// The sender must not accept their own request.
	bobAccepts := FollowRequestType{}
	bobAccepts.InvitationRowType.Id = rq.Id
	if err := bobAccepts.Accept(bob); !errors.Is(err, ferror.ErrorPermissionDenied) {
		t.Errorf("Accept() by sender error = %v, want %v", err, ferror.ErrorPermissionDenied)
	}
	// An unrelated user must not deny it either.
	aliceDenies := FollowRequestType{}
	aliceDenies.InvitationRowType.Id = rq.Id
	if err := aliceDenies.Deny(alice); !errors.Is(err, ferror.ErrorPermissionDenied) {
		t.Errorf("Deny() by unrelated user error = %v, want %v", err, ferror.ErrorPermissionDenied)
	}
}

func TestFollowRequestUnfollow(t *testing.T) {
	alice, bob, _ := setupFollowModelTestDB(t)

	rq, err := requestFollow(bob, alice)
	if err != nil {
		t.Fatalf("request follow: %v", err)
	}
	if err := acceptByRecipient(rq, alice); err != nil {
		t.Fatalf("accept: %v", err)
	}

	unfollower := FollowRequestType{}
	unfollower.ToUserId = alice
	if err := unfollower.Unfollow(); err != nil {
		t.Fatalf("Unfollow() error: %v", err)
	}
	unfollowed, err := reload(rq.Id)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if unfollowed.Status != FollowStatusUnfollowed {
		t.Errorf("status after unfollow = %q, want %q", unfollowed.Status, FollowStatusUnfollowed)
	}

	if err := unfollower.Unfollow(); !errors.Is(err, ferror.ErrorCannotUnfollow) {
		t.Errorf("second Unfollow() error = %v, want %v", err, ferror.ErrorCannotUnfollow)
	}
}

func TestFollowsTypeGetForUser(t *testing.T) {
	alice, bob, carol := setupFollowModelTestDB(t)

	// bob -> alice (pending, incoming for alice / outgoing for bob)
	pending, err := requestFollow(bob, alice)
	if err != nil {
		t.Fatalf("bob asks alice: %v", err)
	}
	// bob -> carol (pending)
	if _, err := requestFollow(bob, carol); err != nil {
		t.Fatalf("bob asks carol: %v", err)
	}
	// alice accepts bob -> alice, making bob a follower of alice.
	if err := acceptByRecipient(pending, alice); err != nil {
		t.Fatalf("accept: %v", err)
	}

	var follows FollowsType
	if err := follows.GetForUser(bob); err != nil {
		t.Fatalf("GetForUser() error: %v", err)
	}
	if len(follows.Following) != 1 || follows.Following[0].Username != "alice" || follows.Following[0].ToUserId != alice {
		t.Errorf("Following = %+v, want 1 (alice)", follows.Following)
	}
	if len(follows.PendingIncoming) != 0 {
		t.Errorf("PendingIncoming = %+v, want 0", follows.PendingIncoming)
	}
	if len(follows.PendingOutgoing) != 1 || follows.PendingOutgoing[0].Username != "carol" || follows.PendingOutgoing[0].ToUserId != carol {
		t.Errorf("PendingOutgoing = %+v, want 1 (carol)", follows.PendingOutgoing)
	}
	if len(follows.Followers) != 0 {
		t.Errorf("Followers = %+v, want 0", follows.Followers)
	}

	var aliceFollows FollowsType
	if err := aliceFollows.GetForUser(alice); err != nil {
		t.Fatalf("GetForUser() error: %v", err)
	}
	if len(aliceFollows.Followers) != 1 || aliceFollows.Followers[0].Username != "bob" || aliceFollows.Followers[0].FromUserId != bob {
		t.Errorf("Followers = %+v, want 1 (bob)", aliceFollows.Followers)
	}
}
