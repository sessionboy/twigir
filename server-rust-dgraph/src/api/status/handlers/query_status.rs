use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use crate::models::statuses as status_model;
use crate::api::status::dgraph::query as dgraph;
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use crate::lib::auth::{ AuthUser, GuardAuthUser };
use json_gettext::JSONGetText;
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ parse };
use serde_json::json;

/*
  1，功能：帖子详情
  2，path: /status/{id}
  3，status_id：帖子的id
*/ 
pub async fn status(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    status_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", status_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let _status: Option<status_model::item::Status> = client
        .find_one(dgraph::STATUS_DETAIL, vars.clone()).await?;

    // 帖子不存在
    if _status.is_none() {
        let msg = get_text!(lang_ctx, lang, "status.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let status = _status.unwrap();
    Res::ok(json!({ "data": status }))
}

/*
  1，功能：用户主页帖子列表
  2，path: /status/{id}
  3，status_id：帖子的id
*/ 
pub async fn status_user_homepage(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    let users: Vec<status_model::item::Status> = client
        .find_list(dgraph::USER_HOME_STATUS, vars).await?;
    let list = pagination(&users, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}
