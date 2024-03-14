
pub const MESSAGE_SCHEMA: &str = "
type Message {
  message_type
  from
  to
  msg
  msg_entities
  created_at
}

message_type: string .
from: uid .
to: uid .
msg: string .
msg_entities: uid .
";
