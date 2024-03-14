use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use dgraph_tonic::{ Mutation };
use json_gettext::JSONGetText;
use crate::models::groups as group_model;
use crate::api::group::dgraph::mutation as dgraph;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::group::services::member as service;
use crate::lib::auth::{ GuardAuthUser };
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ date, parse };
use serde_json::json;

/*
  1，功能：创建小组
  2，path: /group
  3，body：小组信息
  4，logged_user：当前登录用户信息
  5，其他：前端弹窗三步：(1) 名字、简介 -> (2) 隐私设置 -> (3) 头像、封面图
*/ 
pub async fn group_create(
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

    // 检查是否已存在
    let mut vars = HashMap::new(); 
    vars.insert("$group_name", inputs.group_name.as_str());
    let group: Option<group_model::item::GroupWithUid> = client
        .find_one(dgraph::GROUP_NAME_EXIST, vars).await?;
    
    if group.is_some() {
        let msg = get_text!(lang_ctx, lang, "group.err.exist_name").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    let created_at = date::get_utc_now();
    let new_group = json!({
        "uid": "_:group",
        "dgraph.type": "Group",
        "group_name": inputs.group_name.clone(),
        "group_description": inputs.group_description.clone(),
        "access": inputs.access.clone(),
        "visible": inputs.visible.clone(),
        "default_cover": true,
        "is_verified": false,
        "group_creater": {
            "uid": login_user_id.clone()
        },
        "members":{
            "uid": "_:member",
            "dgraph.type": "Member",
            "member_user":{
                "uid": login_user_id.clone()
            }, 
            "is_owner": true,
            "is_admin": true,
            "level": 1,
            "is_anonymously": false,
            "created_at": created_at.clone()
        },
        "members_count": 1,
        "statuses_count": 0,
        "created_at": created_at,
    });
    
    let group_resp = client.mutate(&new_group).await?;
    let group_id = &group_resp.uids["group"];
    let member_id = &group_resp.uids["member"];

    let p = json!({
        "uid": group_id.clone(),
        "owner": {
            "uid": member_id
        }
    });
    client.mutate(&p).await?;
    Res::ok(json!({ 
        "data": { "uid": group_id }
    }))
}

/*
  1，功能：更新小组
  2，path: /group/{id}
  3，body：小组信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn group_update(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    group_id: web::Path<String>,
    body: web::Json<group_model::input::UpdateGroupInput>
) -> ServiceResult {
    let client = &data.into_inner();
    let inputs = &body.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    let _group: Option<group_model::item::GroupWithMe> = client
        .find_one(dgraph::GROUP_WITH_ME, vars)
        .await?;

    // 该小组不存在
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    let group = _group.unwrap();
    let member_me = service::get_member(group.members_me.clone());
    // 不是该小组的成员
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.update.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    // 不是该小组的管理员/组长
    let me = member_me.unwrap();
    if !me.is_admin && !me.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.update.unauthorized").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    // 如果是更新小组名称，检查是否已存在
    if inputs.group_name.is_some() {
        let group_name = inputs.group_name.clone().unwrap();
        let mut vars = HashMap::new(); 
        vars.insert("$group_name", group_name.as_str());
        let group: Option<group_model::item::GroupWithUid> = client
            .find_one(dgraph::GROUP_NAME_EXIST, vars).await?;
        if group.is_some() {
            let msg = get_text!(lang_ctx, lang, "group.err.exist_name").unwrap();
            return Res::err(json!({ "message": msg.as_str().unwrap() }));
        }
    }

    let updated_at = date::get_utc_now();
    let mut update_group = serde_json::to_value(&inputs).unwrap();
    update_group["uid"] = json!(group_id.clone());
    update_group["updated_at"] = json!(updated_at);

    client.mutate(&update_group).await?;
    Res::ok(json!({ 
        "data": { "uid": group_id.as_str() }
    }))
}

/*
  1，功能：删除/解散小组
  2，path: /group/{id}
  3，body：小组信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn group_delete(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    group_id: web::Path<String>
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    // 检查是否是组长，唯有组长才可以解散小组
    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    let _group: Option<group_model::item::GroupWithMe> = client
        .find_one(dgraph::GROUP_WITH_ME, vars).await?;
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    let group = _group.unwrap();
    let member_me = service::get_member(group.members_me.clone());
    // 非小组成员
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.update.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    // 非组长不能解散小组
    let me = member_me.unwrap();   
    if !me.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.delete.unauthorized").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    let q = r#"
        query all($uid: string){
            V as var(func: uid($uid)) {
                M as members
            }
        }
    "#;
    let d = r#"       
        uid(V) * * .
        uid(M) * * .     
    "#;
    let mut _vars = HashMap::new(); 
    _vars.insert("$uid", group_id.as_str());
    let mut mu = Mutation::new();
    mu.set_delete_nquads(d);
    client.upsert_with_vars(q, _vars, mu).await?;
    Res::ok(json!({ 
        "message": "success"
    }))
}
