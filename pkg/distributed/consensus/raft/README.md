# RAFT consensus implementation

## State machine

```mermaid
stateDiagram
    [*] --> Follower : start up
    Follower --> Candidate : timeout, start election
    Candidate --> Candidate : timeout, new election
    Candidate --> Leader : receives votes from  majority of servers
    Candidate --> Follower : discover current leader or new term
    Leader --> Follower : discover server with higher term
```

## Log entry

Each log entry contain the current term number, a log index and a command

### Command

Each command is abitrary 

## Architecture

```mermaid
---
title: Components
---
block-beta

  block:Components
    columns 3
    fsm["Finite state machine"]
    RPC["Remote procedure call"]
    Memstg["Memory log storage"]
    PersistStg["Persistence log storage"]
    Snapshot["Log compression and snapshot"]
    block:Telemetry
        Metrics
        Tracing
    end
  end
```

```mermaid
flowchart LR
    subgraph Cluster
        direction TB
        subgraph LeaderSG["Leader"]
            direction TB
            subgraph Write
                direction LR
                RequestConfigChange
                Apply
                Barrier
                LeaderShipChange
                Restore
                VerifyLeader
            end

        end
        subgraph Common
            subgraph Read
                AppliedIndex
                GetConfiguration
                LastContact
                LastIndex
                Leader
                Snapshot
                State
                Stats
            end
        end

        subgraph Follower
            subgraph FWrite["Write"]
                BootstrapCluster
            end
        end
    end

```
