use actix_web::{ HttpResponse };
use serde::{Serialize, Deserialize};
use crate::lib::errors::{ ServiceError };
use serde_json::{ json, Value as JsonValue };

// 响应体
pub type ServiceResult = anyhow::Result<HttpResponse, ServiceError>;

// 成功响应
#[derive(Deserialize, Serialize)]
pub struct Resok<T> where T: Serialize {
    ok: bool,
    message: String,
    data: Option<T>
}
impl<T: Serialize> Resok<T> {
    pub fn new(data: T) -> Self {
      Resok { ok: true, message: "ok".to_owned(), data: Some(data) }
    }

    pub fn to_json_result(&self) -> Result<HttpResponse, ServiceError> {
        Ok(HttpResponse::Ok().json(self))
    }
}

// 错误响应
#[derive(Default, Deserialize, Serialize)]
pub struct Reserr<T> where T: Serialize {
    ok: bool,
    message: Option<String>,
    errors: Option<T>,
    code: Option<String>
}

impl<T: Serialize> Reserr<T> {

  pub fn new(
    errors: T, 
    message: &str
  ) -> Self {
    Reserr { 
      ok: false, 
      errors: Some(errors), 
      message: Some(message.to_owned()),
      code: None 
    }
  }

  pub fn to_json_result(&self) -> Result<HttpResponse, ServiceError> {
    Ok(HttpResponse::BadRequest().json(self))
  }

  pub fn to_json(
    // errors: T, 
    value: JsonValue
  ) -> Result<HttpResponse, ServiceError> {
    // let err_res = Reserr { 
    //   ok: false, 
    //   errors: Some(errors), 
    //   message: Some(value.get("message").unwrap_or(&json!("")).to_string()),
    //   code: Some(value.get("code").unwrap_or(&json!("200")).to_string()),
    // };
    let err_res = json!({
      "ok": false, 
      "status": value.get("status").unwrap_or(&json!(200)),
      "message": value.get("message"), 
      "errors": value.get("errors")
    });
    Ok(HttpResponse::BadRequest().json(err_res))
  }

}
