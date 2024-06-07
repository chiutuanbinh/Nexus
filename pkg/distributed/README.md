# Distributed feature of a database service

## I. Problem statement

I want to allow Nexus to work in a distributed manner. Doing so will provide the following benefits:

1. High availability: If one instance is down, another one can take its place and serve incoming requests
2. Scale: One instance is limited to whatever resource can be allocated to a single machine (Memory, Processing power, Storage).
3. Prevent dataloss: If data is only stored in 1 place, it can be easily loss or damaged. By allowing multiple instance of the same data exist, we

## II. Glossary

- Instance: 1 service process instance
- Node: unit of identification, one instance may represent multiple nodes
- Cluster: several nodes working together
- Partition: a unit of storage, contain some keys. In the simplest case, we have only 1 partition with all keys.
- Replication: a partition is saved to several e, each copy is call a replication.
- Partition's Leader: a node which receive e to partition. It is the source of truth
- Partition's Followers: nodes which hold partition replication, exclude the Partition's Leader

## III. Scope

1. Implement a simple partition system where a key is assigned to a partition. Assignment is decided through consistent hashing function.
   a. Hashing a key resulted in its correspondent nodeID. The instance which handles that nodeId will be the one to receive the entry
2. Implement a replication system with partition. Each partition will have k replicas. Replication leader and follower is set from the configuration.

3. P1 Concensus:

- Raft: inspired by https://github.com/hashicorp/raft

## IV. Detail

### Consistent hashing

- Consistent hashing only has to remap `k/n` keys when the hash table size changes, with `k` is number of keys and `n` is the number of slots.
- A common implementation is dividing the hashed value space into `n` ranges (usually equal). Each range has a corresponded slot. If a value is hashed into a particular range, then it is placed into that correspondent slot.

  - If a slot is removed, all value mapped to its range will be transfer to the next ranges's slot. Example: A[1:10], B[11:20], C[21:30], remove B, we have A[1:10], C[11:30]
  - If a new slot is created, another slot must share its range to the newly created one. Example: slots A[1:10], B[11:20], we add a new slot C, then new slots can be A[1:5], C[6:10], B[11:20]

- The mentioned above implementation has some drawbacks

  - Upon hash table size change, slots size are no longer guaranteed to be equal.
  - A lot of keys have to be re-hashed upon table size change

- An improved implementation is divide the hashed space into `n` ranges, but each slot is assigned to some non-continuous ranges. If a slot change, some other slots can take other its ranges. In other word, we create a virtual layer of slot, and our slot only map to these virtual slots instead of hard ranges. This implementation is used by Amazon's Dynamo. In their implementation, virtual slot is called a node, and each slot is assign several nodes.
  - Example: A[1:10], B[11:20], C[21:30], X[A,B], Y[C]. A,B,C never change, only X or Y, if they change, their corresponded virtual slots are transfered to another slot

### Partition

- Using a hash function (usually consistent hashing), we assign a partition to an entry. Each server in our cluster only have to handle some partition, therefore, if the database size increase or we need to serve more request, we can just increase the number of server and distribute the paritition among servers.
- Basically we split the data into multiple pieces so each piece can be served independently

### Replication

- If we keep data in 1 server, if the server go down, we lose access to the data. If we split it into multiple partition in multiple server, if one of them go down, we lose access to some of the data. It can be a disaster, if the downed servers do not go back up, or the data is corrupted and cannot be covered. It mean we lose the data forever. If the paritition is tightly coupled, then we lose the entire database

- To solve this problem, we can backup the data, or replicate it. Each partition will have a backup and if a server go down, all its partition can be restore to a new server to serve the request.

  - Offline backups still require downtime, we cannot serve the request when we restore the data.
  - We can use a live backup, which mean the backup server is always running and can stand in for the backed up server immediately if anything bad happen.
  - Because the backup server contain the same data with the main server, we can let it serve read request instead of letting it sit idle. This is basically how replication work. The live backup is called a replica.

- Instead of one replica, we can have multiple ones. Having multiple replicas increase the cost but also increase reliability, and increase read throught put if we decide to serve them to read request.
- This come at some costs

  - Replicas can go out of sync with the primary data, and keeping them in sync is expensive and complicated.
  - Replicas use storage.

- If we also use partition, replication will be applied at partition scope. Each partition will have a set of replicas. We also need to make sure each partition is kept in sync with its replicas.

### Advance topic

- Highly available write: We only allow write to primary server, what happen if it failed? Read from replica is still possible, but write is no longer possible.
  - Dynamo have a solution: hinted handoff and slopy quorum, along with anti-entropy to repair inconsistent data. It significant increase availability, but also introduce inconsistency.
  - Failover, allow a replica to write, as long as the chosen replica is in sync. Doing this automatically is not easy unless you have a dedicated replica in line to be standin.
- Concensus:
  - Membership change
  - Leader election
  - Shared lock

#### Quorum

- Quorum is concept of majority in a group vote. In database scope it is the majority of replica.
- If we want replica to server read requests, we need to know the consistency requirement. It is assume that we trigger replica synchronization on primary write

  - The highest level is all replica must acknowledge write before we return success
  - The lowest level is none replica acknowledge write, only primary.
  - The middle ground is only majority of replica ack write. It is call a write quorum.

- If we read from a replica, we might also need to read from other replica to know if we have the correct data.
  - The highest level is all other replica (and primary) also return the same data
  - The lowest level is we don't read any other replica
  - The middle ground is we need majority of replica to aggree with value (a read quorum)

### Design 

```mermaid

```