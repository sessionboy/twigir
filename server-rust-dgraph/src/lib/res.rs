use actix_web::{ HttpResponse };
use crate::lib::errors::{ ServiceError };
use serde_json::{ json, Value as JsonValue };

// 响应体
pub type ServiceResult = anyhow::Result<HttpResponse, ServiceError>;

pub struct Res{}

impl Res{
  // 成功
  pub fn ok(value: JsonValue) -> Result<HttpResponse, ServiceError> {
    let ok_res = json!({
      "ok": true, 
      "status": value.get("status").unwrap_or(&json!(200)),
      "data": value.get("data"), 
      "message": value.get("message")
    });
    Ok(HttpResponse::Ok().json(ok_res))
  }

  // 失败
  pub fn err(value: JsonValue) -> Result<HttpResponse, ServiceError> {
    let err_res = json!({
      "ok": false, 
      "status": value.get("status").unwrap_or(&json!(200)),
      "message": value.get("message"), 
      "errors": value.get("errors")
    });
    Ok(HttpResponse::BadRequest().json(err_res))
  }
}
