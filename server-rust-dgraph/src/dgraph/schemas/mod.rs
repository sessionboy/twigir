
mod user;
mod status;
mod group;
mod notification;

use user::USER_SCHEMA;
use status::STATUS_SCHEMA;
use group::GROUP_SCHEMA;
use notification::NOTIFICATION_SCHEMA;

pub fn get_schema() -> String {
  // STATUS_SCHEMA.to_string() 
  USER_SCHEMA.to_string() 
  + &STATUS_SCHEMA.to_string() 
  + &GROUP_SCHEMA.to_string() 
  + &NOTIFICATION_SCHEMA.to_string() 
}
