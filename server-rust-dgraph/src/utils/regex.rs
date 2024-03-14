use lazy_static::lazy_static;
use regex::Regex;

lazy_static! {
  // 名字, 支持中文、英文、日文、韩文
  pub static ref RE_NAME: Regex = Regex::new(r"^[\u4e00-\u9fa5a|\u0800-\u4e00|\uac00-\ud7ff|a-zA-Z0-9·]{1,20}$").unwrap();
  // 主名 username
  pub static ref RE_USERNAME: Regex = Regex::new(r"^[a-zA-Z]([a-zA-Z0-9_]{2,20})$").unwrap();
  // 密码 password
  pub static ref RE_PASSWORD: Regex = Regex::new(r"^[a-zA-Z0-9]([a-zA-Z0-9_]{7,15})$").unwrap();
  // 中国大陆号码 phone_number
  pub static ref RE_PHONE: Regex = Regex::new(r"^1[3,5,7,8]\d{9}$").unwrap();
}
