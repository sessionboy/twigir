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
use crate::lib::auth::{ AuthUser,GuardAuthUser };
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ parse };
use serde_json::json;

/*
  1，功能：获取通知列表
  2，path: /notifications
*/ 
pub async fn notifications(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    Res::ok(json!({ "data": "notifications" }))
}
