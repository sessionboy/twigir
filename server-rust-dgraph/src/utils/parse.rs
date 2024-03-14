use actix_web::{ HttpRequest };
use std::net::{ SocketAddr };
use accept_language::{ parse };
use woothee::parser::{ Parser,WootheeResult };

// 获取语言列表
pub fn get_accept_language(req: &HttpRequest)-> String {
  let accept_langs = &req.headers().get("accept-language");
  let langs = match accept_langs {
    None=> vec![],
    Some(lang)=> parse(lang.to_str().unwrap())
  };
  if langs.is_empty() {
    return String::from("zh-CN");
  }
  langs.first().unwrap().clone()
}

// 获取cookie
pub fn get_cookie(req: &HttpRequest)-> Vec<String> {
  let cookie = &req.headers().get("cookie").expect("cookie");
  parse(cookie.to_str().unwrap())
}

// 获取user-agent
pub fn get_agent(req: &HttpRequest) -> WootheeResult {
  let parser = Parser::new();
  let agent = &req.headers().get("user-agent").expect("user-agent");
  parser.parse(agent.to_str().unwrap()).unwrap()
}

// 获取ip地址
pub fn get_ip(req: &HttpRequest) -> Option<SocketAddr> {
  let agent = req.peer_addr();
  agent
}