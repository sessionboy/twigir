use actix_web::{ web };
use std::collections::HashMap;
use crate::api::user::dgraph::status as dgraph;
use crate::models::statuses::item::{ Status, StatusMedia };
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use crate::lib::res::{ Res, ServiceResult };
use crate::lib::auth::{ AuthUser };
use serde_json::json;

/*
  1，功能：用户主页帖子列表
  2，path: /users/{id}/status
  3，user_id：用户id
*/ 
pub async fn user_status(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    user_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$user_id", user_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let statuses: Vec<Status> = client
        .find_sub_list(dgraph::USER_STATUS, vars)
        .await?;
    let list = pagination(&statuses, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}

/*
  1，功能：用户{id}的照片/视频等媒体帖子列表
  2，path: /users/{id}/status
  3，user_id：用户id
  4，query: media_type(媒体类型)、first、after
*/ 
pub async fn user_status_media(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    media_query: web::Query<StatusMedia>,
    user_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();
    let media = media_query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$user_id", user_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    vars.insert("$media_type", media.media_type.as_str());

    let statuses: Vec<Status> = client
        .find_sub_list(dgraph::USER_STATUS_MEDIA, vars)
        .await?;
    let list = pagination(&statuses, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}

/*
  1，功能：用户{id}喜欢的帖子列表
  2，path: /users/{id}/favorite
  3，user_id：用户id
  4，query: media_type(媒体类型)、first、after
*/ 
pub async fn user_status_favorite(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    user_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$user_id", user_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let statuses: Vec<Status> = client
        .find_sub_list(dgraph::USER_STATUS_FAVORITE, vars)
        .await?;
    let list = pagination(&statuses, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}

/*
  1，功能：用户{id}加入的小组的帖子列表
  2，path: /users/{id}/groups_status
  3，user_id：用户id
*/ 
pub async fn user_groups_status(
    data: web::Data<DgraphClient>,
    query: web::Query<PaginationArgs>,
    user_id: web::Path<String>,
    logged_user: AuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id_some();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new(); 
    vars.insert("$user_id", user_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let statuses: Vec<Status> = client
        .find_sub_list(dgraph::USER_GROUPS_STATUS, vars)
        .await?;
    let list = pagination(&statuses, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}
