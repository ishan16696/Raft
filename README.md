## Raft: Distributed Consensus Algorithm.


#### Distributed consensus: 
```
A distributed consensus ensures a consensus of data among nodes in a distributed system or reaches an agreement on a proposal.
```


#### Raft algorithm: 
```
Raft algorithm can be used in distributed systems to achieve the consensus.
In this algorithm distributed consensus is achieved with help of Leader Election and Log Replication of data.
```

Summary of RAFT algorithm:
====

**RAFT setup**:
1. All nodes in a system can only be in 3 states: 
    - **Follower**
    - **Candidate**
    - **Leader**
2. Raft also introduces the "**Term**" concept to identify expired information in a timely manner.
3. Each nodes have timer running(initialised randomly). 
4. Raft algorithm is a leader-based asymmetric model:
    * **Follower**: Initially all the nodes are in Follower state.
    * **Candidate**: Follower Node which triggers the leader-election becomes Candidate.
    * **Leader**: Only one Leader per Term is possible.

**Leader Election:**
1. As the timer is randomly initialised, if time out occurs for one of the node this leads to node change the state from **Follower** to **Candidate** and triggers an **election** to become **Leader**.
2. Candidate node votes for himself and Send Request message to all other nodes asking them to vote for him(Greedy Algorithm).
3. If it get the majority of vote then it will declare himself as a **Leader** and send Heartbeat message to all followers nodes to tell them that it becomes the **Leader** of this Term and reset their Timer.
4. To maintain its authority, an elected **Leader** must continuously send a Heartbeat message to the other nodes in the cluster.
5. If a **follower** does not receive the heartbeat packet during a given time, the leader is considered to have crashed,  **follower** changes its status to **candidate** and starts a leader election.

**Log Replication:**
1. Leader is responsible for the log replication.
2. Client send the request to create/update the data, leader receives the request.
3. Leader will forward the client requests to the followers as AppendEntries message.
4. All follower should acknowledge, if leader receives the majority of follower's acknowledgment then only it will commit this change.

* 1.**Leader Election:**
    * Select one of the node as the cluster Leader. It is achieved by Leader Election.
    * Only Elected Leader node can respond to the client. All other node have to sync-up with Leader.
* 2.**Log Replication:**
    * All the changes have to gone through the Leader.
    * Leader replicates its log to all other followers nodes.

   ![RAFT](raft.png)



### other Useful Resources:
- [RAFT Paper](https://raft.github.io/raft.pdf)
- http://thesecretlivesofdata.com/raft/
- https://medium.com/@kasunindrasiri/understanding-raft-distributed-consensus-242ec1d2f521
