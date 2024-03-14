
use dgraph_tonic::{ 
  Client, Response, Mutate, Mutation, 
  Query, TxnMutated, TxnReadOnly,TxnMutatedType 
};
use serde::{ Serialize, Deserialize };
use std::collections::HashMap;
use std::hash::Hash;
use std::sync::Arc;
use crate::lib::errors::{ ServiceError };

pub mod pagination;

// 响应体
pub type ServiceResult = anyhow::Result<Response, ServiceError>;
pub type ExistResult = anyhow::Result<bool, ServiceError>;
pub type ItemResult<T> = anyhow::Result<Option<T>, ServiceError>;
pub type ListResult<T> = anyhow::Result<Vec<T>, ServiceError>;

#[derive(Debug, Clone)]
pub struct DgraphClient {
    pub client: Arc<Client>
}


impl DgraphClient { 

  // 注意：需要引入 use dgraph_tonic::{ Mutate };
  pub fn mutate_txn(&self) -> TxnMutated {
    self.client.new_mutated_txn()
  }

  // 注意：需要引入 use dgraph_tonic::{ Query };
  pub async fn query_txn(&self) -> TxnReadOnly {
    self.client.new_read_only_txn()
  }

  // 突变: 创建、更新
  pub async fn mutate<T: ?Sized>(&self, p: &T) -> ServiceResult 
  where T: Serialize
  {
    let mut mu = Mutation::new();
    mu.set_set_json(&p)?;
    let mut txn = self.client.new_mutated_txn();
    let response = txn.mutate(mu).await?;
    txn.commit().await?;

    Ok(response)
  }

  // 突变: 删除
  pub async fn delete<T: ?Sized>(&self, p: &T) -> ServiceResult 
  where T: Serialize
  {
    let mut mu = Mutation::new();
    mu.set_delete_json(&p)?;
    let mut txn = self.client.new_mutated_txn();
    let response = txn.mutate(mu).await?;
    txn.commit().await?;

    Ok(response)
  }

  // 批量操作
  pub async fn upsert_with_vars<Q, K, V>(
    &self, 
    query: Q, 
    vars: HashMap<K, V>,
    mu: Mutation
  ) -> ServiceResult 
  where
    Q: Into<String> + Send + Sync,
    K: Into<String> + Send + Sync + Eq + Hash,
    V: Into<String> + Send + Sync,
    String: From<Q>
  {
    let mut txn = self.client.new_mutated_txn();
    let response = txn.upsert_with_vars(query,vars, vec![mu]).await?;
    txn.commit().await?;
    Ok(response)
  }

  // 查询
  pub async fn query<Q, K, V>(
    &self,
    query: Q,
    vars: HashMap<K, V>,
  ) -> ServiceResult 
  where 
  Q: Into<String> + Send + Sync,
  K: Into<String> + Send + Sync + Eq + Hash,
  V: Into<String> + Send + Sync,
  String: From<Q>
  {
    let response = self.client
        .new_read_only_txn()
        .query_with_vars(query, vars)   
        .await?;    
    Ok(response)
  }

  // 查询是否存在
  pub async fn find_exist<Q, K, V>(
    &self,
    query: Q,
    vars: HashMap<K, V>,
  ) -> ExistResult
  where 
  Q: Into<String> + Send + Sync,
  K: Into<String> + Send + Sync + Eq + Hash,
  V: Into<String> + Send + Sync,
  String: From<Q>
  {
    let res = self.client
        .new_read_only_txn()
        .query_with_vars(query, vars)
        .await?;
    let _json = !res.json.is_empty();
    Ok(_json)
  }

  // 查询单条记录，返回 Option<T>
  pub async fn find_one<Q, K, V, T>(
    &self,
    query: Q,
    vars: HashMap<K, V>,
  ) -> ItemResult<T>
  where 
  Q: Into<String> + Send + Sync,
  K: Into<String> + Send + Sync + Eq + Hash,
  V: Into<String> + Send + Sync,
  T: for<'de> Deserialize<'de>,
  String: From<Q>
  {
    let res = self.client.new_read_only_txn().query_with_vars(query, vars).await?;
    let _result: Item<T> = serde_json::from_slice(&res.json)?;
    let result = _result.to_result();
    Ok(result)
  }

  // 查询列表，返回 [..]
  pub async fn find_list<Q, K, V, T>(
    &self,
    query: Q,
    vars: HashMap<K, V>,
  ) -> ListResult<T>
  where 
  Q: Into<String> + Send + Sync,
  K: Into<String> + Send + Sync + Eq + Hash,
  V: Into<String> + Send + Sync,
  T: for<'de> Deserialize<'de>,
  String: From<Q>
  {
    let res = self.client.new_read_only_txn().query_with_vars(query, vars).await?;
    let _result: List<T> = serde_json::from_slice(&res.json)?;
    let result = _result.to_result();
    Ok(result)
  }

  // 查询二级列表，返回 [..]
  pub async fn find_sub_list<Q, K, V, T>(
    &self,
    query: Q,
    vars: HashMap<K, V>,
  ) -> ListResult<T>
  where 
  Q: Into<String> + Send + Sync,
  K: Into<String> + Send + Sync + Eq + Hash,
  V: Into<String> + Send + Sync,
  T: for<'de> Deserialize<'de>,
  String: From<Q>
  {
    let res = self.client.new_read_only_txn().query_with_vars(query, vars).await?;
    let _result: SubList<T> = serde_json::from_slice(&res.json)?;
    let result = _result.to_result();
    Ok(result)
  }

}

/*
  1，用于查询单条记录 
  2，从response的data列表中提取第一个结果 T
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct Item<T> {
  pub data: Option<Vec<T>>
}
impl<T> Item<T> {
  pub fn to_result(self) -> Option<T>{
    self.data?.into_iter().next()
  }
}

/*
  1，用于列表查询 
  2，返回response的 data列表 Vec<T>
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct List<T> {
  pub data: Option<Vec<T>>
}
impl<T> List<T> {
  pub fn to_result(self) -> Vec<T>{
    if self.data.is_none() {
      let result: Vec<T> = Vec::new();
      return result;
    }
    self.data.unwrap()
  }
}


/*
  1，查询二级列表
  2，返回response的 data列表 Vec<T>
  3，注意：查询包装名必须是 data ，例如: data(func: uid($uid)) {...
*/ 
#[derive(Debug,Clone,Serialize,Deserialize)]
pub struct SubList<T> {
  pub data: Option<Vec<SubEdges<T>>>
}
impl<T> SubList<T> {
  pub fn to_result(self) -> Vec<T> {
    if self.data.is_none() {
      let result: Vec<T> = Vec::new();
      return result;
    }
    let list = self.data.unwrap();
    let result = list.into_iter().next();
    let sub_edges: SubEdges<T> = result.unwrap();
    sub_edges.get_edges()
  }
}


/*
  1，说明：二级列表查询 ，比如查询用户下的statuses列表
  2，注意：二级列表需要以"edges"作为别名，比如 
   {
     uid
     edges: statuses {
       uid
       ...
     }
   }
*/ 
#[derive(Debug,Clone,Default,Serialize,Deserialize)]
pub struct SubEdges<T> {
  pub uid: String,
  pub edges: Option<Vec<T>>
}
impl<T> SubEdges<T> {
  pub fn get_edges(self) -> Vec<T> {
    match self.edges {
      Some(edges)=> edges,
      None=> Vec::new()
    }
  }
}