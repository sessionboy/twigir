use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use dgraph_tonic::{ Mutation };
use json_gettext::JSONGetText;
use crate::models::groups as group_model;
use crate::api::group::dgraph::mutation as dgraph;
use crate::lib::dgraph::{ DgraphClient };
use crate::lib::auth::{ GuardAuthUser };
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ date, parse };
use serde_json::json;

/*
  1，功能：新建消息
  2，path: /message
  3，body：小组信息
  4，logged_user：当前登录用户信息
  5，其他：前端弹窗三步：(1) 名字、简介 -> (2) 隐私设置 -> (3) 头像、封面图
*/ 
pub async fn message_create(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    body: web::Json<group_model::input::CreateGroupInput>
) -> ServiceResult {
    let client = &data.into_inner();
    let inputs = &body.into_inner();
    let login_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    Res::ok(json!({ 
        "data": { "uid": login_user_id }
    }))
}
