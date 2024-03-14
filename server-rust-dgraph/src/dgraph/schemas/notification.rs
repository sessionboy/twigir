
pub const NOTIFICATION_SCHEMA: &str = "
type Notification {
  notification_type
  message
  sender
  receiver
  target
  group
  status
  reply
  user
  created_at
}

notification_type: string .
message: string .
sender: uid .
receiver: [uid] .
target: uid .
group: uid .
status: uid .
reply: uid .
user: uid .
";
