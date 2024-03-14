use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use crate::api::group::dgraph::query as dgraph;
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use json_gettext::JSONGetText;
use crate::models::groups as group_model;
use crate::models::statuses as status_model;
use crate::lib::auth::{ AuthUser };
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ parse };
use serde_json::json;

/*
  1，功能：小组详情
  2，path: /group/{id}
  3，status_id：帖子的id
*/ 
pub async fn group(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    group_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let result: Option<group_model::item::Group> = client
        .find_one(dgraph::GROUP_DETAIL, vars)
        .await?;

    // 小组不存在
    if result.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let group = result.unwrap();
    Res::ok(json!({ "data": group }))
}


/*
  1，功能：小组{id}的帖子列表
  2，path: /group/{id}/status
  3，user_id：用户id
*/ 
pub async fn group_status(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    group_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let result: Vec<status_model::item::Status> = client
        .find_sub_list(dgraph::GROUP_STATUS, vars)
        .await?;
    let list = pagination(&result, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}