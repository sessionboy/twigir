use actix_web::{ HttpResponse,ResponseError };
use serde::{Serialize, Deserialize};
use thiserror::Error;
use serde_json::{ Value as JsonValue };
use dgraph_tonic::DgraphError;
use crate::lib::response::{ Reserr };
use crate::lib::res::Res;

#[derive(Debug, Error, Serialize, Deserialize)]
pub enum ServiceError {

    // 参数错误
    #[error("ValidationError")]
    ValidationError(Vec<JsonValue>),

    // 访问错误
    #[error("AuthError")]
    AuthError(JsonValue),

    // 400 参数错误
    #[error("BadRequest")]
    BadRequest(String),

    // 401 未授权
    #[error("未授权")]
    Unauthorized(String),

    // 403 禁止访问
    #[error("禁止访问")]
    Forbidden(String),

    // 404 找不到
    #[error("找不到")]
    NotFound(&'static str),

    // 422 无法处理
    #[error("Unprocessable Entity")]
    UnprocessableEntity(String),

    // 500 服务器错误
    #[error("服务器异常")]
    InternalServerError,
}

impl ServiceError {
    fn to_error(&self) -> i32 {
        let error = &self.to_string();
        error.parse().unwrap_or(-1)
    }
}

// 转换数据库错误
impl From<DgraphError> for ServiceError {
  fn from(error: DgraphError) -> ServiceError {
      println!("{:#?}",&error);
      ServiceError::InternalServerError
  }
}

// 转换数据库错误
impl From<anyhow::Error> for ServiceError {
  fn from(error: anyhow::Error) -> ServiceError {
      println!("DgraphError: {:#?}",&error);
      ServiceError::InternalServerError
  }
}

// 转换json错误
impl From<serde_json::Error> for ServiceError {
  fn from(error: serde_json::Error) -> ServiceError {
      println!("{:#?}",&error);
      ServiceError::InternalServerError
  }
}

impl ResponseError for ServiceError {
    fn error_response(&self) -> HttpResponse { 
      match *self {
        // 400
        ServiceError::ValidationError(ref message) => {
          let resp = Reserr::new(message, "参数错误");
          return HttpResponse::BadRequest().json(resp);
        },
        // 1002
        ServiceError::AuthError(ref message) => {
          Res::err(message.clone()).unwrap()
        },
        // 400
        ServiceError::BadRequest(ref message) => {
          let resp = Reserr::new(self.to_error(), &message.as_str());
          return HttpResponse::BadRequest().json(resp);
        },
        // 401
        ServiceError::Unauthorized(ref message) => {
          return HttpResponse::Unauthorized().json(
            Reserr::new(self.to_error(), &message.as_str())
          );
        },
        // 403
        ServiceError::Forbidden(ref message) => {
          return HttpResponse::Forbidden().json(
            Reserr::new(self.to_error(), &message.as_str())
          );
        },
        // 404
        ServiceError::NotFound(ref message) => {
          return HttpResponse::NotFound().json(
            Reserr::new(self.to_error(), &message)
          );
        },
        // 422
        ServiceError::UnprocessableEntity(ref message) => {          
          return HttpResponse::BadRequest().json(
            Reserr::new(self.to_error(), &message.as_str())
          );
        },
        // 500
        ServiceError::InternalServerError => {
          return HttpResponse::InternalServerError().json(
            Reserr::new(self.to_error(), "服务器异常")
          )
        },
        _=> {
          return HttpResponse::BadRequest().json(
            Reserr::new(self.to_error(), "未知错误")
          );
        }
      }
    }
}
