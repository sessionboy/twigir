pub mod schemas;
pub use schemas::get_schema;
use dgraph_tonic::{ Client, Operation };

pub fn get_client() ->Client {
  // http://localhost:9080
  let client: Client = Client::new("http://47.99.243.195:9080").expect("connected client");
  client
}

pub async fn drop_all() {
  let client = get_client();
  let op = Operation {
      drop_all: true,
      ..Default::default()
  };
  client.alter(op).await.expect("set schema");
  println!("数据库已清空!");
}

pub async fn set_schema(){
  let client = get_client();
  let schema = get_schema();
  let op = Operation {
      schema,
      ..Default::default()
  };
  client.alter(op).await.expect("set schema");
  println!("schema已设置成功!");
}