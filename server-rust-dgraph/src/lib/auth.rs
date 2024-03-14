use actix_web::{ dev::Payload, FromRequest, HttpRequest };
use crate::lib::errors::ServiceError;
use actix_identity::RequestIdentity;
use crate::utils::jwt::{ decode_token };
use crate::models::users::auth::SlimUser;
use serde_json::{ json };

// 1，获取登录用户信息，不做拦截
#[derive(Debug, Clone, Default)]
pub struct AuthUser(pub Option<SlimUser>);
impl From<SlimUser> for AuthUser {
    fn from(slim_user: SlimUser) -> Self {
      AuthUser(Some(slim_user))
    }
}


impl AuthUser {
  pub fn to_slim (self) -> SlimUser {
    let user = self.0.unwrap();
    SlimUser {
      uid: user.uid,
      username: user.username,
      role: user.role
    }
  }
  pub fn user_id (self) -> String {
    match &self.0 {
      Some(user)=> user.uid.clone(),
      None=> String::new()
    }
  }
  pub fn user_id_some (self) -> String {
    match &self.0 {
      Some(user)=> user.uid.clone(),
      None=> String::from("0x0")
    }
  }
}


/*
  1，功能：仅获取登录用户信息，不做拦截
*/ 
impl FromRequest for AuthUser {
  type Error = ServiceError;
  type Future = futures::future::Ready<Result<Self, Self::Error>>;
  type Config = ();

  fn from_request(req: &HttpRequest, _: &mut Payload) -> Self::Future {
      let identity = RequestIdentity::get_identity(req);     
      if let Some(identity) = identity {
        let slim_user = decode_token(&identity);
        match slim_user {
          Ok(claims) => {
            return futures::future::ready(
              Ok(AuthUser(Some(
                SlimUser{
                  uid: claims.uid,
                  role: claims.role,
                  username: claims.username
                }
              ))));
          },
          Err(_) => return futures::future::ready(Ok(AuthUser(None)))
        }
      } else {
        return futures::future::ready(Ok(AuthUser(None)));
      }
  }
}


/*
  1，功能：未登录拦截处理，并获取登录用户信息
*/ 
#[derive(Debug, Clone, Default)]
pub struct GuardAuthUser(pub Option<SlimUser>);
impl From<SlimUser> for GuardAuthUser {
    fn from(slim_user: SlimUser) -> Self {
      GuardAuthUser(Some(slim_user))
    }
}
impl GuardAuthUser {
  pub fn to_slim (self) -> SlimUser {
    let user = self.0.unwrap();
    SlimUser {
      uid: user.uid,
      username: user.username,
      role: user.role
    }
  }
  pub fn user_id (self) -> String {
    let user_id = self.0.unwrap().uid;
    user_id
  }
}

impl FromRequest for GuardAuthUser {

  type Error = ServiceError;
  type Future = futures::future::Ready<Result<Self, Self::Error>>;
  type Config = ();

  fn from_request(req: &HttpRequest, _: &mut Payload) -> Self::Future {
      let identity = RequestIdentity::get_identity(req);     
      if let Some(identity) = identity {
        // 有token
        let slim_user = decode_token(&identity);
        match slim_user {
          Ok(claims) => {
            return futures::future::ready(
              Ok(GuardAuthUser(Some(
                SlimUser{
                  uid: claims.uid,
                  role: claims.role,
                  username: claims.username
                }
              ))));
          },
          Err(_) => return futures::future::ready(Err(
            ServiceError::AuthError(json!({ 
              "message":"token已过期", 
              "status": 1001
            }))
          ))
        }
      } else {

        // 没有token
        return futures::future::ready(Err(
          ServiceError::AuthError(json!({ 
            "message":"未登录", 
            "status": 1002
          }))
        ));
      }
  }
}