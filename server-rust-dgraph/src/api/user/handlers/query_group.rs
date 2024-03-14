use actix_web::{ web };
use std::collections::HashMap;
use crate::api::user::dgraph::group as dgraph;
use crate::models::groups::item as group_model;
use crate::models::commons::{ List };
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use crate::lib::res::{ Res, ServiceResult };
use crate::lib::auth::{ AuthUser };
use serde_json::json;

/*
  1，功能：用户{id}的照片/视频等媒体帖子列表
  2，path: /users/{id}/status
  3，user_id：用户id
  4，query: media_type(媒体类型)、first、after
*/ 
pub async fn user_groups(
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
    
    let users: Vec<group_model::GroupItem> = client
        .find_list(dgraph::USER_GROUPS, vars)
        .await?;
    let list = pagination(&users, qr.get_first(), _after);
    Res::ok(json!({ "data": list }))
}