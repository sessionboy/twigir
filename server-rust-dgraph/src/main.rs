#[macro_use]
extern crate validator_derive;
extern crate validator;
#[macro_use]
extern crate serde_derive;
extern crate serde_json;
#[macro_use]
extern crate dotenv_codegen;
extern crate anyhow;
extern crate regex;
extern crate chrono;
#[macro_use]
extern crate smart_default;
extern crate accept_language;
extern crate woothee;
#[macro_use] extern crate json_gettext;

pub mod config;
pub mod lib;
pub mod api;
pub mod dgraph;
pub mod models;
mod utils;
mod locales;

mod dgraph_orm;

use actix_web::{web, App, HttpServer };
use actix_web::http::header;
use actix_cors::Cors;
use chrono::Duration;
use csrf_token::CsrfTokenGenerator;
use dgraph::{ get_client, drop_all, set_schema };
use std::sync::Arc;
use lib::dgraph::DgraphClient;

async fn index() -> &'static str {
    "Hello world!\r\n"
}


#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct User {
  pub name: String,
  pub username: String,
  pub description: Option<String>,
  pub age: u32,
  pub verified: bool
}

#[actix_rt::main]
async fn main() -> std::io::Result<()> {

    // 清空数据库
    // drop_all().await;
    // 设置schema
    // set_schema().await;

    let client = Arc::new(get_client());

    let csrf_token_header = header::HeaderName::from_lowercase(b"x-csrf-token").unwrap();
    HttpServer::new(move || {
        App::new()
            .wrap(
                Cors::new()
                    .send_wildcard()
                    .allowed_methods(vec!["GET", "POST", "PUT", "DELETE"])
                    .allowed_headers(
                        vec![header::AUTHORIZATION,
                            header::CONTENT_TYPE,
                            header::ACCEPT,
                            csrf_token_header.clone()
                        ]
                    )
                    .expose_headers(vec![csrf_token_header.clone()])
                    .max_age(3600)
                    .finish()
            )
            .data(
                CsrfTokenGenerator::new(
                    dotenv!("CSRF_TOKEN_KEY").as_bytes().to_vec(),
                    Duration::hours(1)
                )
            )
            .wrap(utils::identity::get_identity_service())
            // 本地化数据
            .data(locales::langs_ctx())
            // dgraph
            .data(DgraphClient {
                client: client.clone()
            })
            // web::Json 错误处理
            .app_data(
                web::JsonConfig::default()
                .error_handler(lib::error_handler::json_error_handler)
            )
            .app_data(
                web::QueryConfig::default()
                .error_handler(lib::error_handler::query_error_handler)
            )
            .app_data(
                web::PathConfig::default()
                .error_handler(lib::error_handler::path_error_handler)
            )
            .service(
                web::scope("/v1")
                .configure(api::user::route)
                .configure(api::status::route)
                .configure(api::group::route)
                .configure(api::notification::route)
                .configure(api::message::route)
            )           
            .route("/", web::get().to(index))
    })
    .bind("127.0.0.1:8088")?
    .run()
    .await
}