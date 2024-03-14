use actix_web::{ web };
use std::collections::HashMap;
use crate::api::user::dgraph::relation as dgraph;
use crate::models::users as user_model;
use crate::lib::dgraph::{ 
    DgraphClient,
    pagination::{ pagination, PaginationArgs }
};
use crate::lib::res::{ Res, ServiceResult };
use crate::lib::auth::{ AuthUser, GuardAuthUser };
use serde_json::{ json };

/*
  1，功能：获取{id}的关注列表
  2，path: /users/{id}/followings
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/ 
pub async fn get_followings(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    query: web::Query<PaginationArgs>,
    logged_user: AuthUser
) -> ServiceResult {

    let client = &data.into_inner();
    let qr = query.into_inner();

    let logged_user_id = logged_user.user_id_some();
    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::USERS_FOLLOWINGS, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);

    Res::ok(json!({ "data": result }))
}


/*
  1，功能：获取{id}的粉丝列表
  2，path: /users/{id}/followers
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/ 
pub async fn get_followers(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    query: web::Query<PaginationArgs>,
    logged_user: AuthUser
) -> ServiceResult {

    let client = &data.into_inner();
    let qr = query.into_inner();

    let logged_user_id = logged_user.user_id_some();
    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::USERS_FOLLOWERS, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": result }))
}


/*
  1，功能：获取{id}的好友列表
  2，path: /users/{id}/friends
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/ 
pub async fn get_friends(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    query: web::Query<PaginationArgs>,
    logged_user: AuthUser
) -> ServiceResult {

    let client = &data.into_inner();
    let qr = query.into_inner();

    let logged_user_id = logged_user.user_id_some();
    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::USERS_FRIENDS, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": result }))
}


/*
  1，功能：    {id_1}和{id_2} 的共同关注
  2，path:    /users/{id}/same_followings
  3，user_id：目标用户的id
  4，query：  分页参数first、after等
  5，说明：    暂时限制为已登录状态用户{id_1}
*/ 
pub async fn get_same_followings(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    logged_user: GuardAuthUser,
    query: web::Query<PaginationArgs>
) -> ServiceResult {

    let login_user = logged_user.to_slim();
    let client = &data.into_inner();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid_one", login_user.uid.as_str()); 
    vars.insert("$uid_two", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::USERS_SAME_FOLLOWINGS, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": result }))
}


/*
  1，功能：    我关注的{id_1}、{id_2}等人也关注了{user_id}
  2，path:    /users/{id}/relation_followings
  3，user_id：目标用户的id
  4，query：  分页参数first、after等
  5，说明：    暂时限制为已登录状态用户，查询我的关注列表里有哪些人关注了{user_id}
*/ 
pub async fn get_relation_followings(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    logged_user: GuardAuthUser,
    query: web::Query<PaginationArgs>
) -> ServiceResult {
  
    let login_user = logged_user.to_slim();
    let client = &data.into_inner();
    let qr = query.into_inner();

    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid_one", login_user.uid.as_str()); 
    vars.insert("$uid_two", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::ONE_FOLLOWINGS_TOFOLLOW_TWO, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": result }))
}


/*
  1，功能：获取{id}的屏蔽列表
  2，path: /users/{id}/shields
  3，user_id：目标用户的id
  4，query：分页参数first、after等
*/ 
pub async fn get_shields(
    data: web::Data<DgraphClient>,
    user_id: web::Path<String>,
    query: web::Query<PaginationArgs>,
    logged_user: AuthUser
) -> ServiceResult {

    let client = &data.into_inner();
    let qr = query.into_inner();

    let logged_user_id = logged_user.user_id_some();
    let _first = qr.get_add_first_str();
    let _after = qr.get_after();
    let mut vars = HashMap::new();
    vars.insert("$uid", user_id.as_str()); 
    vars.insert("$first", _first.as_str());
    vars.insert("$after", _after.as_str());
    vars.insert("$logged_user_id", logged_user_id.as_str());

    let list: Vec<user_model::item::Follow> = client
        .find_sub_list(dgraph::USERS_SHIELDS, vars)
        .await?;
    let result = pagination(&list, qr.get_first(), _after);
    Res::ok(json!({ "data": result }))
}
