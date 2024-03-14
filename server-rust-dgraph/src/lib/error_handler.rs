use actix_web::{ error,HttpRequest,HttpResponse };
use serde_json::{ json };

// web::Json 错误处理
pub fn json_error_handler(err: error::JsonPayloadError, _req: &HttpRequest) -> error::Error {
  println!("{:?}",err);
  let error = json!({ 
    "ok": false, 
    "errors": null, 
    "status": 1001,
    "message": "参数错误"
  });
  let response = match &err {
      error::JsonPayloadError::ContentType => {
          HttpResponse::UnsupportedMediaType().json(error)
      }
      _ => HttpResponse::BadRequest().json(error)
  };
  error::InternalError::from_response(err, response).into()
}

// web::Path 错误处理
pub fn path_error_handler(err: error::PathError, _req: &HttpRequest) -> error::Error {
  println!("{:?}",err);
  let error = json!({ 
    "ok": false, 
    "errors": null, 
    "status": 1001,
    "message": "参数错误"
  });
  let response = HttpResponse::BadRequest().json(error);
  error::InternalError::from_response(err, response).into()
}

pub fn query_error_handler(err: error::QueryPayloadError, _req: &HttpRequest) -> actix_web::Error {
  // let error_message: String = match &err {
  //   error::QueryPayloadError::Deserialize(deserialize_error) => format!("{}", deserialize_error),
  // };
  println!("{:?}",err);
  let error = json!({ 
    "ok": false, 
    "errors": null, 
    "status": 1001,
    "message": "参数错误"
  });
  let response = HttpResponse::BadRequest().json(error);
  error::InternalError::from_response(err, response).into()
}
