use crate::models::users::auth::{ LoggedUser,LoggedUserWithPassword };
use crate::models::users::exist::UserExist;
use crate::models::users::input::RegisterInput;
use crate::models::users::args::LoginArg;
use std::collections::HashMap;
use crate::api::user::dgraph::auth as gql;
use serde_json::{ Value as JsonValue };
use serde_json::json;

// 检查是否已注册
pub fn register_exist_errors(
  user_exist: &UserExist, 
  user: &RegisterInput
) -> Vec<JsonValue> {
  let mut errors: Vec<JsonValue> = Vec::new();
  if !user_exist.name.is_empty() {
    errors.push(json!({
        "path": "name",
        "value": &user.name
    }));    
  }
  if !user_exist.username.is_empty() {
    errors.push(json!({
        "path": "username",
        "value": &user.username
    }));    
  }
  if !user_exist.phone.is_empty() {
    errors.push(json!({
        "path": "phone_number",
        "value": &user.phone_number
    }));    
  }
  errors
}

// 查询登录用户 
pub fn login_query(
  body: &LoginArg
) -> (&'static str, HashMap<&str, String>) {
  let mut q: &'static str = "";
  let mut vars: HashMap<&str, String> = HashMap::new();
  if body.username.is_some() {
    vars.insert("$username", body.username.clone().unwrap());
    q = gql::LOGIN_BY_USERNAME;
  } 
  if body.phone_number.is_some() && body.phone_code.is_some() {
    vars.insert("$phone_number", body.phone_number.clone().unwrap());
    q = gql::LOGIN_BY_PHONE;
  }
  if body.email.is_some() {
    vars.insert("$email", body.email.clone().unwrap());
    q = gql::LOGIN_BY_EMAIL;
  } 

  (q, vars)
}

// 登录用户信息
pub fn login_user(user: LoggedUserWithPassword) -> LoggedUser {
  LoggedUser{
    uid: user.clone().uid,
    name: user.clone().name,
    username: user.clone().username,
    role: user.clone().role,
    avatar_url: user.clone().avatar_url,
    lang: user.clone().lang,
    description: user.clone().description,
    is_verified: user.clone().is_verified,
    verified: user.clone().verified
  }
}
