use serde::{Deserialize, Serialize};

#[derive(Debug,SmartDefault,Clone,Serialize,Deserialize)]
pub enum NotificationType {
  #[default]
  User_follow           // {user_id}关注了你
  User_unfollow         // {user_id}取关了你
  User_friend           // {user_id}关注了你，与你成为好友

  Status_aite,          // {user_id}在帖子中@了你
  Status_favorite,      // {user_id}点赞了你的帖子
  Status_forward        // {user_id}转发了你的帖子
  Status_reply          // {user_id}回复了你的帖子

  Reply                 // {user_id}回复了你的回复
  Reply_aite,           // {user_id}在回复中@了你
  Reply_favorite,       // {user_id}点赞了你的回复
  Reply_forward         // {user_id}转发了你的回复

  Group_join            // {user_id}加入了{group_id}小组
  Group_leave           // {user_id}退出了{group_id}小组
  Group_admin           // {user_id}成为了{group_id}小组的管理员
  Group_quit_admin      // {user_id}不再是{group_id}小组的管理员
  Group_forbidden       // 您在{group_id}小组中被禁言n天

  System_admin          // 恭喜您成为社区管理员
  System_quit_admin     // 您不再是社区管理员
}
