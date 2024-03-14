use actix_web::{ web, HttpRequest };
use std::collections::HashMap;
use json_gettext::JSONGetText;
use dgraph_tonic::{ Mutation, Mutate };
use crate::models::groups as group_model;
use crate::api::group::dgraph::mutation as dgraph;
use crate::lib::dgraph::{ DgraphClient };
use crate::api::group::services::member as service;
use crate::lib::auth::{ GuardAuthUser };
use crate::lib::res::{ Res, ServiceResult };
use crate::utils::{ date, parse };
use serde_json::json;

/*
  1，功能：加入小组
  2，path: /group/{id}/join
  3，body：小组信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn group_join(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    group_id: web::Path<String>,
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    // 检查是否已加入
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
    if member_me.is_some() {
        // 已加入小组，不能再加入
        let msg = get_text!(lang_ctx, lang, "group.err.joined").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    let updated_at = date::get_utc_now();
    let p = json!({
        "uid": group_id.clone(),
        "members":{
            "uid": "_:member",
            "dgraph.type": "Member",
            "member_user":{
                "uid": logged_user_id.clone()
            }, 
            "is_owner": false,
            "is_admin": false,
            "level": 1,
            "is_anonymously": false,
            "created_at": updated_at.clone()
        },
        "members_count": group.members_count + 1,
        "updated_at": updated_at.clone()
    });
    client.mutate(&p).await?;
    let msg = get_text!(lang_ctx, lang, "group.success.joined").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能：退出小组
  2，path: /group/{id}/leave
  3，body：小组信息
  4，logged_user：当前登录用户信息
*/ 
pub async fn group_leave(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    group_id: web::Path<String>,
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);

    // 检查是否已加入
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

    // 不是该小组成员
    let member_me = service::get_member(group.members_me.clone());
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.leave.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }

    // 组长不能退出小组
    let me = member_me.unwrap();
    if me.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.leave.owner").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let mut txn = client.mutate_txn();
    let updated_at = date::get_utc_now();

    // 1，删除成员
    let p1 = json!({
        "uid": group_id.clone(),
        "members":{
            "uid": me.uid
        }
    });
    let mut mu_del = Mutation::new();
    mu_del.set_delete_json(&p1)?;
    txn.mutate(mu_del).await?;

    // 2，更新小组信息
    let p2 = json!({
        "uid": group_id.clone(),
        "members_count": group.members_count - 1,
        "updated_at": updated_at
    });
    let mut mu_mut = Mutation::new();        
    mu_mut.set_set_json(&p2)?;        
    txn.mutate(mu_mut).await?;

    let res = txn.commit().await;
    if res.is_err() {
        let msg = get_text!(lang_ctx, lang, "common.action_failure").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let msg = get_text!(lang_ctx, lang, "group.success.leaved").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能：添加小组管理员
  2，path: /group/{id}/add_admin/{member_id}
  3，body：小组信息
  4，logged_user：当前登录用户信息,
  5，ids：参数，第0个是小组id，第1个是成员id
*/ 
pub async fn group_add_admin(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    ids: web::Path<(String,String)>,
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);
    let group_id = &ids.0;
    let member_id = &ids.1;

    // 检查是否是组长，唯有组长才可以添加管理员
    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$member_id", member_id.as_str());
    let _group: Option<group_model::item::GroupWithMeAndMember> = client
        .find_one(dgraph::GROUP_WITH_ME_AND_MEMBER, vars).await?;
        
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let group = _group.unwrap();
    let member_me = service::get_member(group.members_me.clone());
    // "我"不是该小组成员
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.update.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    // "我"不是组长，不能添加管理员
    let me = member_me.unwrap();
    if !me.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.unauthorized").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let _member = service::get_member(group.members.clone());
    // {member_id}不是该小组成员，不能添加为管理员
    if _member.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let member = _member.unwrap();
    // 已经是组长
    if member.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.owner").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    // 已经是管理员
    if member.is_admin {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.admin").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let updated_at = date::get_utc_now();
    let p = json!({
        "uid": group_id.clone(),
        "members":{
            "uid": member_id,
            "is_admin": true
        },
        "updated_at": updated_at
    });
    client.mutate(&p).await?;
    let msg = get_text!(lang_ctx, lang, "group.success.add_admin").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能：开除小组管理员
  2，path: /group/{id}/fire_admin/{member_id}
  3，body：小组信息
  4，logged_user：当前登录用户信息,
  5，ids：参数，第0个是小组id，第1个是成员id
*/ 
pub async fn group_fire_admin(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    ids: web::Path<(String,String)>,
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);
    let group_id = &ids.0;
    let member_id = &ids.1;

    // 检查是否是组长，唯有组长才可以开除管理员
    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$member_id", member_id.as_str());
    let _group: Option<group_model::item::GroupWithMeAndMember> = client
        .find_one(dgraph::GROUP_WITH_ME_AND_MEMBER, vars).await?;
    
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let group = _group.unwrap();
    let member_me = service::get_member(group.members_me.clone());
    // "我"不是该小组成员
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.update.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    // "我"不是组长，不能开除管理员
    let me = member_me.unwrap();
    if !me.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.unauthorized").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let _member = service::get_member(group.members.clone());
    // {member_id}不是该小组成员，不能开除管理员
    if _member.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.member.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let member = _member.unwrap();
    // {member_id}是组长，不能执行此操作
    if member.is_owner {
        let msg = get_text!(lang_ctx, lang, "group.err.fire_admin.is_owner").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    // {member_id}不是管理员
    if !member.is_admin {
        let msg = get_text!(lang_ctx, lang, "group.err.fire_admin").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let updated_at = date::get_utc_now();
    // 更新is_admin
    let p = json!({
        "uid": group_id.clone(),
        "members":{
            "uid": member_id,
            "is_admin": false
        },
        "updated_at": updated_at
    });
    client.mutate(&p).await?;
    let msg = get_text!(lang_ctx, lang, "group.success.fire_admin").unwrap();
    Res::ok(json!({ "message": msg }))
}

/*
  1，功能：禁言
  2，path: /group/{id}/forbidden/{member_id}
  3，logged_user：当前登录用户信息,
  4，ids：参数，第0个是小组id，第1个是成员id
*/ 
pub async fn group_forbidden(
    data: web::Data<DgraphClient>,
    lang_ctx: web::Data<JSONGetText<'static>>,
    req: HttpRequest,
    logged_user: GuardAuthUser,
    ids: web::Path<(String,String)>,
) -> ServiceResult {
    let client = &data.into_inner();
    let logged_user_id = logged_user.user_id();
    let lang = parse::get_accept_language(&req);
    let group_id = &ids.0;
    let member_id = &ids.1;

    let mut vars = HashMap::new(); 
    vars.insert("$uid", group_id.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());
    vars.insert("$member_id", member_id.as_str());
    let _group: Option<group_model::item::GroupWithMeAndMember> = client
        .find_one(dgraph::GROUP_WITH_ME_AND_MEMBER, vars).await?;
    
    if _group.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }
    let group = _group.unwrap();
    let member_me = service::get_member(group.members_me.clone());
    // "我"不是该小组成员
    if member_me.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.update.not_found").unwrap();
        return Res::err(json!({ "message": msg.as_str().unwrap() }));
    }
    // "我"不是组长/管理员，不能禁言
    let me = member_me.unwrap();
    if !me.is_owner && !me.is_admin {
        let msg = get_text!(lang_ctx, lang, "group.err.add_admin.unauthorized").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let _member = service::get_member(group.members.clone());
    // {member_id}不是该小组成员，不能禁言
    if _member.is_none() {
        let msg = get_text!(lang_ctx, lang, "group.err.member.not_found").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let member = _member.unwrap();
    // {member_id}是组长/管理员，不能禁言
    if member.is_owner || member.is_admin {
        let msg = get_text!(lang_ctx, lang, "group.err.forbidden.owner_or_admin").unwrap();
        return Res::err(json!({ "message": msg }));
    }

    let updated_at = date::get_utc_now();
    // 三天后解禁
    let forbidden_date = date::utc_now_distance_day(3);
    let p = json!({
        "uid": group_id.clone(),
        "members":{
            "uid": member_id,
            "forbidden_date": forbidden_date,
            "updated_at": updated_at
        }
    });
    client.mutate(&p).await?;
    let msg = get_text!(lang_ctx, lang, "group.success.fire_admin").unwrap();
    Res::ok(json!({ "message": msg }))
}
