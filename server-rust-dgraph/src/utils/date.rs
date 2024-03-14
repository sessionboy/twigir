use chrono::prelude::*;
use chrono::Duration;

/*
   1，功能: 和当前时间相比相差多少秒
   2，utc_date: 以 rfc3339 格式存储在数据库的utc时间
*/ 
pub fn compare_from_now_secs(utc_date: &String)-> i64 {
  let _date = DateTime::parse_from_rfc3339(utc_date.as_str()).unwrap();
  let _utc_now = Utc::now();
  let res = _utc_now.signed_duration_since(_date);
  res.num_seconds()
}

/*
   1，功能: 是否在当前时间之前
*/ 
pub fn before_from_now(utc_date: &String)-> bool {
  compare_from_now_secs(&utc_date) < 0
}

/*
   1，功能: 获取当前utc时间
   2，utc_date: rfc3339 格式
*/ 
pub fn get_utc_now()-> String {
  Utc::now().to_rfc3339().to_string()
}

/*
   1，功能: 获取n天后或-n天前的的utc时间
   2，utc_date: rfc3339 格式
*/ 
pub fn utc_now_distance_day(day: i64)-> String {
  (Utc::now() + Duration::days(day)).to_rfc3339().to_string()
}
