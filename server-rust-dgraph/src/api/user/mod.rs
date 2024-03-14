use actix_web::{ web };

pub mod dgraph;
pub mod handlers;
pub mod services;

use handlers::{
  mutate_auth,
  mutate_account,
  mutate_relation,
  query_group,
  query_relation,
  query_status,
  query_user,
};

pub fn route(cfg: &mut web::ServiceConfig) {
    cfg
    .service( 
      // 账号相关
      web::scope("/auth")
        .service(web::resource("/register").route(web::post().to(mutate_auth::register)))
        .service(web::resource("/login").route(web::post().to(mutate_auth::login)))    
        .service(web::resource("/logout").route(web::get().to(mutate_auth::logout)))  
        .service(web::resource("/send_phonecode").route(web::post().to(mutate_auth::send_phonecode)))      
        .service(web::resource("/verify_phonecode").route(web::post().to(mutate_auth::verify_phonecode)))      
        .service(web::resource("/send_emailcode").route(web::post().to(mutate_auth::send_emailcode)))      
        .service(web::resource("/verify_emailcode").route(web::post().to(mutate_auth::verify_emailcode)))             
    )
    .service( 
      // 更新用户资料
      web::scope("/account")
        .service(web::resource("/profile").route(web::put().to(mutate_account::profile)))
        .service(web::resource("/name").route(web::put().to(mutate_account::account_name)))
        .service(web::resource("/username").route(web::put().to(mutate_account::account_username)))
        .service(web::resource("/phone").route(web::put().to(mutate_account::account_phone)))
        .service(web::resource("/email").route(web::put().to(mutate_account::account_email)))
        .service(web::resource("/password").route(web::put().to(mutate_account::account_password)))
    )
    .service( 
      web::scope("/users")
        // 查询用户信息
        .service(web::resource("/me").route(web::get().to(query_user::me)))
        .service(web::resource("/{id}").route(web::get().to(query_user::user_detail)))
        // 查询用户关系列表
        .service(web::resource("/{id}/followings").route(web::get().to(query_relation::get_followings)))
        .service(web::resource("/{id}/followers").route(web::get().to(query_relation::get_followers)))   
        .service(web::resource("/{id}/friends").route(web::get().to(query_relation::get_friends)))     
        .service(web::resource("/{id}/same_followings").route(web::get().to(query_relation::get_same_followings)))
        .service(web::resource("/{id}/relation_followings").route(web::get().to(query_relation::get_relation_followings)))       
        .service(web::resource("/{id}/shields").route(web::get().to(query_relation::get_shields)))
        // 更新用户关系
        .service(web::resource("/{id}/follow").route(web::post().to(mutate_relation::follow)))
        .service(web::resource("/{id}/unfollow").route(web::post().to(mutate_relation::unfollow)))      
        .service(web::resource("/{id}/shield").route(web::post().to(mutate_relation::shield)))   
        .service(web::resource("/{id}/unshield").route(web::post().to(mutate_relation::unshield)))        
        // 查询用户帖子
        .service(web::resource("/{id}/status").route(web::get().to(query_status::user_status)))
        .service(web::resource("/{id}/media_status").route(web::get().to(query_status::user_status_media)))
        .service(web::resource("/{id}/favorites").route(web::get().to(query_status::user_status_favorite)))
        .service(web::resource("/{id}/groups_status").route(web::get().to(query_status::user_groups_status)))
        // 查询用户小组
        .service(web::resource("/{id}/groups").route(web::get().to(query_group::user_groups)))
    );
}
