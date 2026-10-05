CREATE TRIGGER group_member_accepted AFTER UPDATE OF status ON group_invitations
WHEN NEW.status = 'accepted' AND OLD.status != 'accepted'
BEGIN
  INSERT INTO group_message_reads (group_id, user_id, last_read_message_id)
  VALUES (NEW.group_id, NEW.to_user_id,
          COALESCE((SELECT MAX(id) FROM group_messages WHERE group_id = NEW.group_id), 0))
  ON CONFLICT(group_id, user_id)
  DO UPDATE SET last_read_message_id = excluded.last_read_message_id;
END;

CREATE TRIGGER group_owner_created AFTER INSERT ON groups
BEGIN
  INSERT INTO group_message_reads (group_id, user_id, last_read_message_id)
  VALUES (NEW.id, NEW.owner_user_id, 0);
END;

-- Existing owners
INSERT OR IGNORE INTO group_message_reads (group_id, user_id, last_read_message_id)
SELECT g.id, g.owner_user_id,
       COALESCE((SELECT MAX(id) FROM group_messages WHERE group_id = g.id), 0)
FROM groups g;

-- Existing accepted members
INSERT OR IGNORE INTO group_message_reads (group_id, user_id, last_read_message_id)
SELECT gi.group_id, gi.to_user_id,
       COALESCE((SELECT MAX(id) FROM group_messages WHERE group_id = gi.group_id), 0)
FROM group_invitations gi
WHERE gi.status = 'accepted';