use crate::dgraph_orm::model::DgraphModel;
use crate::lib::dgraph::DgraphClient;

#[derive(Debug, Clone)]
pub struct Dgraph {
  pub client: DgraphClient,
  pub models: DgraphModel,
}

impl Dgraph {
  fn init(&self){
    // let models = self.models;
    // for model in models.iter() {
    //   println!("{:?}",model);
    // }
  }
}