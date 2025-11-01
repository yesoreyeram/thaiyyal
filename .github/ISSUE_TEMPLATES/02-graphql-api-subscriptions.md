---
name: GraphQL API with Real-time Subscriptions
about: Implement GraphQL API layer with real-time subscriptions for flexible querying and collaboration
title: '[EPIC] GraphQL API Layer with Real-time Subscriptions'
labels: ['epic', 'enhancement', 'priority:high', 'complexity:high', 'area:api', 'area:frontend']
assignees: ''
---

## Overview

Implement a comprehensive GraphQL API layer with real-time subscriptions to enable flexible data querying, reduce over-fetching, and support real-time collaboration features.

## Problem Statement

**Current State**: No API layer exists. Planned REST API requires multiple round-trips and lacks real-time capabilities.

**Desired State**: GraphQL API that provides:
- Single endpoint for all data operations
- Flexible querying (fetch exactly what you need)
- Real-time subscriptions for workflow updates
- Efficient data loading (no N+1 queries)
- Strong typing and introspection

## Business Value

- **Developer Experience**: Easier frontend development with type safety
- **Performance**: Reduced network traffic and round-trips
- **Real-time Collaboration**: Multiple users editing workflows simultaneously
- **Flexibility**: Frontend can evolve without backend changes
- **API Documentation**: Self-documenting via introspection

## Technical Requirements

### Schema Design

```graphql
# Core Types
type User {
  id: ID!
  email: String!
  name: String!
  createdAt: DateTime!
  workflows(first: Int, after: String): WorkflowConnection!
}

type Workflow {
  id: ID!
  name: String!
  description: String
  nodes: [Node!]!
  edges: [Edge!]!
  version: Int!
  createdAt: DateTime!
  updatedAt: DateTime!
  owner: User!
  executions(first: Int, after: String, status: ExecutionStatus): ExecutionConnection!
  collaborators: [User!]!
  permissions: WorkflowPermissions!
}

type Node {
  id: ID!
  type: NodeType!
  data: JSON!
  position: Position!
}

type Edge {
  id: ID!
  source: String!
  target: String!
  sourceHandle: String
  targetHandle: String
}

type Execution {
  id: ID!
  workflowId: ID!
  workflow: Workflow!
  status: ExecutionStatus!
  startTime: DateTime!
  endTime: DateTime
  results: JSON
  logs: [ExecutionLog!]!
  metrics: ExecutionMetrics!
}

type ExecutionLog {
  id: ID!
  timestamp: DateTime!
  nodeId: String
  level: LogLevel!
  message: String!
  metadata: JSON
}

# Enums
enum NodeType {
  NUMBER
  TEXT_INPUT
  HTTP_REQUEST
  VISUALIZATION
  MATH_OPERATION
  TEXT_OPERATION
  TRANSFORM
  EXTRACT
  CONDITION
  FOR_EACH
  WHILE_LOOP
  SWITCH
  PARALLEL
  JOIN
  SPLIT
  VARIABLE
  CACHE
  ACCUMULATOR
  COUNTER
  RETRY
  TRY_CATCH
  TIMEOUT
  DELAY
}

enum ExecutionStatus {
  PENDING
  RUNNING
  COMPLETED
  FAILED
  CANCELLED
}

enum LogLevel {
  DEBUG
  INFO
  WARN
  ERROR
}

# Mutations
type Mutation {
  # Auth
  register(input: RegisterInput!): AuthPayload!
  login(input: LoginInput!): AuthPayload!
  logout: Boolean!
  
  # Workflows
  createWorkflow(input: CreateWorkflowInput!): Workflow!
  updateWorkflow(id: ID!, input: UpdateWorkflowInput!): Workflow!
  deleteWorkflow(id: ID!): Boolean!
  
  # Execution
  executeWorkflow(workflowId: ID!, input: JSON): Execution!
  cancelExecution(id: ID!): Execution!
  
  # Collaboration
  shareWorkflow(workflowId: ID!, userId: ID!, permission: Permission!): Workflow!
  unshareWorkflow(workflowId: ID!, userId: ID!): Workflow!
}

# Queries
type Query {
  # User
  me: User
  user(id: ID!): User
  users(first: Int, after: String, search: String): UserConnection!
  
  # Workflows
  workflow(id: ID!): Workflow
  workflows(first: Int, after: String, filter: WorkflowFilter): WorkflowConnection!
  
  # Executions
  execution(id: ID!): Execution
  executions(first: Int, after: String, filter: ExecutionFilter): ExecutionConnection!
  
  # Analytics
  workflowStats(workflowId: ID!, period: TimePeriod!): WorkflowStats!
}

# Subscriptions
type Subscription {
  # Workflow changes
  workflowUpdated(workflowId: ID!): Workflow!
  
  # Execution updates
  executionUpdated(executionId: ID!): Execution!
  executionLogAdded(executionId: ID!): ExecutionLog!
  
  # Collaboration
  userPresence(workflowId: ID!): UserPresence!
  nodeSelection(workflowId: ID!): NodeSelection!
}

# Inputs
input CreateWorkflowInput {
  name: String!
  description: String
  nodes: [NodeInput!]!
  edges: [EdgeInput!]!
}

input UpdateWorkflowInput {
  name: String
  description: String
  nodes: [NodeInput!]
  edges: [EdgeInput!]
}

input NodeInput {
  id: String!
  type: NodeType!
  data: JSON!
  position: PositionInput!
}

input EdgeInput {
  id: String!
  source: String!
  target: String!
  sourceHandle: String
  targetHandle: String
}

# Pagination
type WorkflowConnection {
  edges: [WorkflowEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type WorkflowEdge {
  node: Workflow!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Technology Stack

**GraphQL Server**:
- Option A: gqlgen (Go, type-safe, code generation)
- Option B: graphql-go (Go, flexible)

**Subscriptions**:
- WebSocket transport (graphql-ws protocol)
- Redis Pub/Sub for distributed subscriptions

**Data Loading**:
- DataLoader pattern for batch loading
- Caching layer for frequently accessed data

## Acceptance Criteria

### Phase 1: Schema & Infrastructure (Sprint 1-2)
- [ ] Design complete GraphQL schema
- [ ] Set up gqlgen (or chosen framework)
- [ ] Configure code generation
- [ ] Implement schema directives
- [ ] Set up GraphQL Playground/Apollo Studio
- [ ] Design resolver structure
- [ ] Create schema documentation

### Phase 2: Core Resolvers (Sprint 3-4)
- [ ] Implement User resolvers
  - [ ] Query: me, user, users
  - [ ] Field resolvers with DataLoader
- [ ] Implement Workflow resolvers
  - [ ] Query: workflow, workflows
  - [ ] Mutation: create, update, delete
  - [ ] Field resolvers
- [ ] Implement Execution resolvers
  - [ ] Query: execution, executions
  - [ ] Mutation: execute, cancel
  - [ ] Field resolvers
- [ ] Implement pagination
  - [ ] Cursor-based pagination
  - [ ] Connection pattern
  - [ ] totalCount optimization

### Phase 3: DataLoader & Optimization (Sprint 5)
- [ ] Implement DataLoader for Users
- [ ] Implement DataLoader for Workflows
- [ ] Implement DataLoader for Executions
- [ ] Add batch loading
- [ ] Add caching layer
- [ ] Optimize N+1 queries
- [ ] Add query complexity analysis
- [ ] Implement depth limiting

### Phase 4: Real-time Subscriptions (Sprint 6-7)
- [ ] Set up WebSocket server
- [ ] Implement subscription resolvers
  - [ ] workflowUpdated
  - [ ] executionUpdated
  - [ ] executionLogAdded
  - [ ] userPresence
  - [ ] nodeSelection
- [ ] Implement Redis Pub/Sub
- [ ] Add subscription authentication
- [ ] Add subscription authorization
- [ ] Implement connection management
- [ ] Add heartbeat mechanism
- [ ] Handle disconnections gracefully

### Phase 5: Security & Auth (Sprint 8)
- [ ] Implement authentication middleware
- [ ] Add JWT token validation
- [ ] Implement authorization directives
  - [ ] @auth
  - [ ] @hasRole
  - [ ] @isOwner
- [ ] Add field-level permissions
- [ ] Implement rate limiting per user
- [ ] Add query cost calculation
- [ ] Implement query timeout
- [ ] Add introspection security

### Phase 6: Frontend Integration (Sprint 9)
- [ ] Generate TypeScript types
- [ ] Create Apollo Client setup
- [ ] Implement authentication flow
- [ ] Create GraphQL hooks
- [ ] Implement optimistic updates
- [ ] Add cache normalization
- [ ] Implement error handling
- [ ] Add loading states
- [ ] Create subscription components

### Phase 7: Testing & Documentation (Sprint 10)
- [ ] Write resolver unit tests (>80% coverage)
- [ ] Write integration tests
- [ ] Write subscription tests
- [ ] Load test subscriptions
- [ ] Create API documentation
- [ ] Write GraphQL usage guide
- [ ] Create example queries
- [ ] Document best practices

## Technical Design

### Resolver Implementation

```go
package graph

type Resolver struct {
    db          *sql.DB
    userLoader  *dataloader.Loader
    workflowLoader *dataloader.Loader
    pubsub      *redis.Client
}

// User Resolver
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    // Use DataLoader to batch
    return r.userLoader.Load(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context, first *int, after *string, search *string) (*model.UserConnection, error) {
    // Implement cursor-based pagination
    query := r.db.Query().
        Where(user.NameContains(*search)).
        Order(ent.Asc(user.FieldCreatedAt))
    
    // Apply cursor
    if after != nil {
        cursor, _ := decodeCursor(*after)
        query = query.Where(user.IDGT(cursor))
    }
    
    // Apply limit
    limit := getLimit(first)
    users, _ := query.Limit(limit + 1).All(ctx)
    
    // Build connection
    return buildConnection(users, limit)
}

// Subscription Resolver
func (r *subscriptionResolver) WorkflowUpdated(ctx context.Context, workflowId string) (<-chan *model.Workflow, error) {
    // Check authorization
    if !canAccessWorkflow(ctx, workflowId) {
        return nil, ErrUnauthorized
    }
    
    // Create channel
    ch := make(chan *model.Workflow, 1)
    
    // Subscribe to Redis
    topic := fmt.Sprintf("workflow:%s:updates", workflowId)
    subscription := r.pubsub.Subscribe(ctx, topic)
    
    // Handle messages
    go func() {
        defer close(ch)
        for {
            select {
            case <-ctx.Done():
                return
            case msg := <-subscription.Channel():
                var workflow model.Workflow
                json.Unmarshal([]byte(msg.Payload), &workflow)
                ch <- &workflow
            }
        }
    }()
    
    return ch, nil
}
```

### DataLoader Implementation

```go
package dataloader

func NewUserLoader(db *sql.DB) *dataloader.Loader {
    return dataloader.NewBatchedLoader(func(ctx context.Context, keys []string) []*dataloader.Result {
        // Batch load users
        users, err := db.User.Query().
            Where(user.IDIn(keys...)).
            All(ctx)
        
        if err != nil {
            return dataloader.NewResultsWithError(len(keys), err)
        }
        
        // Map results to keys
        userMap := make(map[string]*ent.User)
        for _, u := range users {
            userMap[u.ID] = u
        }
        
        // Return in correct order
        results := make([]*dataloader.Result, len(keys))
        for i, key := range keys {
            if user, ok := userMap[key]; ok {
                results[i] = &dataloader.Result{Data: user}
            } else {
                results[i] = &dataloader.Result{Error: ErrNotFound}
            }
        }
        
        return results
    })
}
```

### Frontend Apollo Client Setup

```typescript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from 'graphql-ws';

const httpLink = new HttpLink({
  uri: 'http://localhost:8080/graphql',
  headers: {
    authorization: `Bearer ${getToken()}`,
  },
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:8080/graphql',
    connectionParams: {
      authorization: `Bearer ${getToken()}`,
    },
  })
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink
);

export const client = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache({
    typePolicies: {
      Workflow: {
        fields: {
          executions: {
            keyArgs: ['filter'],
            merge(existing, incoming, { args }) {
              // Cursor-based pagination merge
              if (!existing) return incoming;
              return {
                ...incoming,
                edges: [...existing.edges, ...incoming.edges],
              };
            },
          },
        },
      },
    },
  }),
});
```

## Non-Functional Requirements

- **Performance**:
  - P95 query latency < 50ms
  - P99 query latency < 200ms
  - Support 1000+ concurrent subscriptions
  - Query complexity limit: 1000 points

- **Security**:
  - Authentication required for all operations
  - Field-level authorization
  - Rate limiting: 1000 requests/minute per user
  - Introspection disabled in production

- **Scalability**:
  - Horizontal scaling for query servers
  - Redis Pub/Sub for distributed subscriptions
  - Connection pooling

## Dependencies

- [ ] PostgreSQL database with workflows schema
- [ ] Redis for Pub/Sub and caching
- [ ] JWT authentication system
- [ ] Monitoring (Prometheus + Grafana)

## Risks & Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| N+1 query problems | High | High | Use DataLoader, optimize queries |
| Complex query DOS | High | Medium | Query complexity analysis, depth limiting |
| Subscription memory leaks | High | Medium | Connection cleanup, heartbeat monitoring |
| Breaking schema changes | Medium | Medium | Schema versioning, deprecation workflow |

## Success Metrics

- [ ] <50ms P95 latency for queries
- [ ] Zero N+1 queries
- [ ] 1000+ concurrent subscriptions
- [ ] <1% error rate
- [ ] 100% type safety on frontend

## Timeline

**Estimated Effort**: 15-20 person-days  
**Recommended Team**: 1 backend engineer + 1 frontend engineer  
**Duration**: 10 weeks

## Related Issues

- #TBD: Type generation for frontend
- #TBD: Real-time collaboration features
- #TBD: GraphQL federation (future)

## References

- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [gqlgen Documentation](https://gqlgen.com/)
- [Apollo Client Documentation](https://www.apollographql.com/docs/react/)
- [DataLoader Pattern](https://github.com/graphql/dataloader)
