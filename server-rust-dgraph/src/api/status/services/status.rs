use crate::models::statuses::input::CreateStatusInput;
use crate::models::statuses::entity_input::EntityInput;
use crate::models::replies::input::CreateReplyInput;
use serde_json::{ Value, json };
use crate::utils::{ date };

pub fn parse_status(
  body: &CreateStatusInput,
  login_user_id: String 
) -> Value {
  
  let created_at = date::get_utc_now();
  let mut status_value = json!({
     "uid": "_:status",
     "dgraph.type": "Status",
     "user": {
       "uid": login_user_id
     },
     "text": body.text,
     "is_forward": body.forward_to_status.is_some(),
     "replies_count": 0,
     "forwards_count": 0,
     "favorites_count": 0,
     "created_at": created_at
  });

  if body.forward_to_status.is_some() {
    status_value["forward_to_status"] = json!({
      "uid": body.forward_to_status
    })
  }

  if body.entities.is_some() {
    let entities = body.entities.clone().unwrap();
    status_value["entities"] = parse_entities(entities);
  }

  status_value
}

pub fn parse_reply(
  body: &CreateReplyInput,
  login_user_id: String 
) -> Value {

  let created_at = date::get_utc_now();
  let mut status_value = json!({
     "uid": "_:reply",
     "dgraph.type": "Reply",
     "user": {
       "uid": login_user_id
     },
     "text": body.text,
     "is_to_reply": body.to_reply.is_some(),
     "replies_count": 0,
     "forwards_count": 0,
     "favorites_count": 0,
     "created_at": created_at
  });

  if body.to_reply.is_some() {
    status_value["to_reply"] = json!({
      "uid": body.to_reply
    })
  }

  if body.entities.is_some() {
    let entities = body.entities.clone().unwrap();
    status_value["entities"] = parse_entities(entities);
  }

  status_value
}

pub fn parse_entities(
  entities: EntityInput
) -> Value {

  let mut entities_value = json!({
    "uid": "_:entity",
    "dgraph.type": "Entity",
    "media_type": entities.media_type
  });

  if entities.mentions.is_some() {
    let _mentions = entities.mentions.unwrap();
    let mentions = serde_json::to_value(&_mentions).unwrap();
    entities_value["mentions"] = mentions;
  }

  if entities.hashtags.is_some() {
    let _hashtags = entities.hashtags.unwrap();
    let hashtags = serde_json::to_value(&_hashtags).unwrap();
    entities_value["hashtags"] = hashtags;
  }

  if entities.urls.is_some() {
    let mut urls_list: Vec<serde_json::Value>  = Vec::new();
    let _urls = entities.urls.unwrap();
    let mut urls = serde_json::to_value(_urls).unwrap();
    let urls_arr = urls.as_array_mut().unwrap().to_vec();
    for mut url in urls_arr {
      url["dgraph.type"] = json!("Url");
      urls_list.push(url);
    }
    entities_value["urls"] = json!(urls_list);
  }

  if entities.medias.is_some() {
    let mut medias_list: Vec<serde_json::Value>  = Vec::new();
    let _medias = entities.medias.unwrap();
    let mut medias = serde_json::to_value(_medias).unwrap();
    let medias_arr = medias.as_array_mut().unwrap().to_vec();
    for mut media in medias_arr {
      media["dgraph.type"] = json!("Media");
      medias_list.push(media);
    }
    entities_value["medias"] = json!(medias_list);
  }

  entities_value
  
}
