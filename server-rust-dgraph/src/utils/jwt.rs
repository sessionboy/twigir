use chrono::{Duration, Utc};
use jsonwebtoken::{decode, encode, Header,EncodingKey,DecodingKey,Validation};
use crate::models::users::auth::SlimUser;
use crate::models::users::_enum::Role;
use crate::lib::errors::{ ServiceError };

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub iat: i64,
    pub exp: i64,
    pub uid: String,
    pub username: String,
    pub role: Role
}

impl From<Claims> for SlimUser {
  fn from(claims: Claims) -> Self {
      SlimUser {
          uid: claims.uid,
          username: claims.username,
          role: claims.role
      }
  }
}

pub fn generate_token(
  user: SlimUser
) -> String {
  let iat = Utc::now().timestamp();
  let exp = (Utc::now() + Duration::days(7)).timestamp();
  let claims = Claims { 
    iat,
    exp,
    uid: user.uid, 
    username: user.username,
    role: user.role
  };

  let header = Header::default();
  encode(
      &header,
      &claims,
      &EncodingKey::from_secret("secret".as_ref())
  )
  .expect("token")
}

pub fn decode_token(token: &str) -> Result<Claims,ServiceError> {
  decode::<Claims>(
      token,
      &DecodingKey::from_secret("secret".as_ref()),
      &Validation::default(),
  )
  .map(|data| data.claims)
  .map_err(|e| ServiceError::BadRequest(e.to_string()))
}