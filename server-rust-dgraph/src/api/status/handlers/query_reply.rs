use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use crate::models::replies as reply_model;
use crate::api::status::dgraph::reply as dgraph_reply;
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use crate::lib::auth::{ AuthUser };
use json_gettext::JSONGetText;
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ parse };
use serde_json::json;

/*
  1，功能：回复详情
  2，path: /reply/{id}
  3，reply_id：帖子的id
*/ 
pub async fn reply(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    reply_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", reply_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    let _reply: Option<reply_model::query::Reply> = client
        .find_one(dgraph_reply::REPLY_DETAIL, vars)
        .await?;

    // 回复不存在
    if _reply.is_none() {
        let msg = get_text!(lang_ctx, lang, "status.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let reply = _reply.unwrap();
    
    Res::ok(json!({ "data": reply }))
}

/*
  1，功能：帖子的回复列表
  2，path: /status/{id}/replies
  3，status_id：帖子的id
*/ 
pub async fn status_replies(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    status_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$uid", status_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let list: Vec<reply_model::query::StatusReply> = client
        .find_sub_list(dgraph_reply::STATUS_REPLIES, vars)
        .await?;
    let replies = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": replies }))
}

/*
  1，功能：回复的回复列表
  2，path: /reply/{id}/replies
  3，reply_id：帖子的id
*/ 
pub async fn reply_replies(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    reply_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$uid", reply_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let list: Vec<reply_model::query::StatusReply> = client
        .find_sub_list(dgraph_reply::STATUS_REPLIES, vars)
        .await?;
    let replies = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": replies }))
}