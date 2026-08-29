# Kin Product Scope and Domain Map

## Purpose

This document defines Kin at a complete but intentionally low-resolution product and domain level. It exists to give humans and coding agents a stable conceptual map before implementation details are chosen.

It describes what the system is responsible for, which business concepts matter, and how major domains interact. It does **not** define database schemas, HTTP endpoints, provider SDK usage, deployment topology, or other infrastructure choices.

Some domain boundaries in this document are hypotheses. They should evolve when real user behavior, domain discovery, or implementation evidence shows that a different boundary is more coherent.

---

## Product mission

Kin helps real friends maintain meaningful relationships by preserving lightweight, permissioned context about each other's lives and interests.

The product is based on a simple observation: adults often still care about their friends, but they lose day-to-day context. When enough context disappears, restarting a conversation takes more effort than it should.

Kin aims to preserve that context continuity without asking users to constantly post, publish, or manage a traditional social feed.

Conceptually:

`Digital Footprints -> AI-Curated Identity Signals -> Permissioned Friendship Context -> Better Human Conversation`

Kin should reduce the effort required to remember what matters to a friend and increase the chance that a useful, natural conversation starts.

---

## What Kin is

Kin is an ambient friendship-context system.

It may ingest explicitly authorized activity signals, derive higher-level social context, apply relationship-aware privacy rules, and surface only the context that is likely to help a real friendship.

Examples of useful context:

- A friend has recently become interested in a new topic.
- A friend is repeatedly researching a specific problem.
- Two friends are independently exploring the same subject.
- A friend's recent behavior suggests a meaningful change in interests.
- A previously active friendship has gone quiet and there is a relevant reason to reconnect.

The system should prefer derived meaning over raw activity logs.

For example:

- Raw signal: liked five videos about AI agents.
- Derived context: recently exploring AI agents for product workflows.

---

## What Kin is not

Kin is not intended to be:

- a high-volume public social network;
- an engagement-maximization feed;
- a system that republishes every digital action;
- a surveillance product;
- a replacement for direct human communication;
- a generic life-logging database;
- a product where users must constantly remember to post in order for it to be useful.

The guiding principle is:

> We do not want the user's attention. We want to give it back to their friends.

---

## Core product principles

### 1. Context continuity over content production

Kin should help friends stay aware of meaningful changes without requiring continuous manual posting.

### 2. Derived context over raw activity

Raw behavior is private input. Shared output should usually be higher-level meaning.

The preferred transformation is:

`Raw Behavior -> Interest Detection -> Significant Change -> Friendship Context`

not:

`Raw Behavior -> Social Feed`

### 3. Permissioned by default

A user should control whether a signal can contribute to shared friendship context and at what relationship level that context may be visible.

### 4. Sensitive by default

When uncertainty exists, the system should prefer withholding or abstracting information rather than oversharing it.

### 5. Relationship-aware disclosure

The same underlying signal may produce different projections for different relationships.

Example:

- Friend: "Recently researching AI products."
- Close Friend: "Researching AI products for non-technical founders."
- Inner Circle: "Thinking about an AI product specifically for friendship maintenance."

### 6. Explainability

When context is surfaced, Kin should be able to explain at a high level why it was considered relevant and why it was allowed to be shared.

### 7. Human conversation is the outcome

The product should optimize for meaningful conversation and relationship continuity rather than impressions, scrolling time, or posting volume.

---

## Product success

Primary success signals should reflect relationship value rather than attention capture.

Candidate outcome metrics include:

- conversations started from Kin context;
- dormant friendships reactivated;
- friendships maintained over time;
- users reporting better awareness of close friends' lives;
- context items that lead to meaningful follow-up;
- useful shared-interest discoveries.

DAU, session time, and feed impressions may be operational metrics, but they are not the primary product purpose.

---

## Primary actors

### User

A person who owns a Kin identity, contributes or authorizes activity signals, manages relationship/privacy settings, receives context about friends, and may initiate conversation from that context.

### Friend

Another Kin user connected through an accepted relationship.

### Relationship participant

A user viewed specifically in the context of one friendship. The same user can have different permissions and projections across different relationships.

### External activity provider

A system that can contribute authorized signals, such as a music, video, reading, browser, AI assistant, or other digital service.

Providers are external actors. They are not part of Kin's domain model and must not define the product architecture.

### AI context generator

An external capability used to classify, summarize, normalize, or derive context. It is an adapter capability, not a source of domain truth.

---

## High-level capability map

Kin's long-term product capabilities can be grouped into the following areas.

### Identity and relationship setup

- establish a user identity;
- discover or invite another person;
- connect two people;
- represent relationship closeness;
- manage relationship lifecycle.

### Activity contribution

- receive explicitly authorized digital activity;
- accept manual contribution when integrations are unavailable;
- normalize heterogeneous activity into Kin concepts;
- distinguish raw signals from socially meaningful context.

### Social-context derivation

- identify recurring interests;
- detect meaningful changes;
- infer candidate context from multiple signals;
- suppress low-signal or repetitive information;
- attach provenance/explanation metadata.

### Privacy and sharing

- determine whether a candidate context may be shared;
- project different detail levels by relationship;
- honor user preferences and sensitive-data rules;
- support opt-out and revocation.

### Relevance and prioritization

- determine whether a valid context item is worth surfacing now;
- rank context for a specific relationship;
- suppress noise and stale information;
- identify serendipitous overlap.

### Friend pulse

- summarize what is worth knowing about a friend;
- surface a small number of useful context items;
- avoid a high-volume chronological feed.

### Conversation support

- generate a context-aware conversation starter;
- help a user ask a useful question about a friend;
- explain why a topic may be worth discussing;
- track whether context led to interaction when appropriate.

### Integration

- connect external providers;
- maintain provider authorization state;
- ingest provider events or pull activity;
- isolate provider-specific payloads and errors.

### Notification

- notify users when meaningful context is available;
- avoid excessive interruption;
- support digest and low-frequency delivery patterns.

### Social memory

Longer-term, Kin may help a user remember meaningful shared history, recurring interests, or prior context without turning private conversation into indiscriminate searchable surveillance.

---

## Candidate domains / bounded contexts

The following are candidate domain boundaries, not final service boundaries.

### Identity

**Responsibility**

Represents a Kin user as a participant in the product.

**Owns**

- stable user identity;
- basic profile-level information required by Kin;
- account-level lifecycle state;
- user-level preferences that are not specific to one friendship.

**Does not own**

- friendship state;
- provider-specific OAuth internals;
- relationship-specific privacy;
- derived social context.

**Classification candidate**

Supporting domain.

---

### Friendship

**Responsibility**

Represents the explicit relationship between two Kin users and the relationship state required for friendship-aware product behavior.

**Owns**

- relationship creation/acceptance;
- active/inactive relationship lifecycle;
- relationship participants;
- relationship-level closeness or trust classification;
- relationship-specific state that is not itself a privacy policy.

**Does not own**

- user identity;
- raw activity;
- social-context generation;
- notification delivery.

**Important invariants**

- a friendship refers to valid participants;
- relationship transitions are explicit;
- relationship-specific behavior must not silently apply globally to a user.

**Classification candidate**

Core domain.

---

### Activity

**Responsibility**

Represents authorized signals about what a user has done, consumed, saved, liked, watched, listened to, researched, or otherwise chosen to contribute to Kin.

**Owns**

- normalized activity concepts;
- activity provenance;
- authorization/contribution intent relevant to the activity;
- timestamps and semantic activity metadata;
- distinction between raw/normalized input and derived context.

**Does not own**

- provider connection lifecycle;
- whether a friend may see the activity;
- final social-context wording;
- friend-specific relevance ranking.

**Important principle**

An activity is not automatically a shareable social item.

**Classification candidate**

Core or supporting domain; this boundary should be revisited after real ingestion use cases exist.

---

### Social Context

**Responsibility**

Turns one or more activities or other authorized signals into a socially meaningful, human-readable context candidate.

**Owns**

- context candidates;
- semantic topic/meaning;
- significance/change interpretation;
- context lifecycle such as candidate, approved, suppressed, expired;
- provenance linking context back to source signals at an abstract level.

**Does not own**

- provider API payloads;
- friendship relationship state;
- final friend-specific privacy projection;
- delivery channel.

**Important invariants**

- raw provider output is not a domain context object;
- context should be validated/normalized before entering the domain;
- one raw action should not automatically become one context item.

**Classification candidate**

Core domain.

---

### Privacy & Sharing

**Responsibility**

Determines whether and at what level a context item can be disclosed to a specific relationship.

**Owns**

- disclosure policy;
- relationship-aware visibility rules;
- user sharing preferences;
- sensitive-content policy;
- projection/detail level;
- revocation and suppression decisions.

**Does not own**

- friendship lifecycle itself;
- generation of the semantic context;
- ranking of otherwise valid context;
- external provider authorization.

**Important invariants**

- default to the least revealing valid representation;
- absence of permission is not permission;
- raw activity should not leak merely because derived context is shareable;
- a disclosure decision is relationship-specific.

**Classification candidate**

Core domain.

---

### Relevance

**Responsibility**

Determines which already-valid context is useful enough to surface for a specific user/relationship at a specific time.

**Owns**

- ranking signals;
- novelty/staleness decisions;
- repetition suppression;
- shared-interest opportunity detection;
- prioritization for a pulse or digest.

**Does not own**

- whether context is legally/privacy-allowed to be shared;
- semantic generation of the context;
- notification transport.

**Classification candidate**

Core or supporting domain depending on how differentiated Kin's ranking behavior becomes.

---

### Conversation / Interaction

**Responsibility**

Represents product behavior that helps friendship context turn into a real interaction.

**Owns**

- conversation-starting intents;
- context-based prompts/questions;
- interaction outcomes that Kin explicitly tracks;
- product-level state for "start a conversation" flows.

**Does not own**

- private third-party chat histories;
- general messaging infrastructure unless Kin later becomes a messaging product;
- friendship permissions;
- context derivation.

**Classification candidate**

Core product experience, but likely a supporting domain in early MVP implementation.

---

### Integration

**Responsibility**

Manages connections to external systems that can contribute or receive data.

**Owns**

- provider connection state;
- external authorization lifecycle;
- sync/checkpoint state;
- provider-specific normalization boundary;
- error translation and retries as integration concerns.

**Does not own**

- Kin domain activities after normalization;
- social-context meaning;
- friendship state;
- privacy disclosure policy.

**Classification candidate**

Generic/supporting domain.

---

### Notification

**Responsibility**

Delivers already-authorized, already-relevant product messages through appropriate channels and timing.

**Owns**

- notification intent;
- delivery scheduling/preferences;
- channel selection;
- delivery status.

**Does not own**

- whether context may be disclosed;
- whether context is semantically meaningful;
- friendship state.

**Classification candidate**

Generic/supporting domain.

---

## Domain interactions

These interactions describe conceptual dependencies, not synchronous service calls.

### Activity contribution to social context

`Integration -> Activity -> Social Context`

An external provider may contribute a signal through the Integration boundary. The signal is normalized into an Activity before Social Context reasons about its meaning.

Provider-specific DTOs must not cross into the Activity or Social Context domain model.

### Context disclosure

`Social Context + Friendship + Privacy & Sharing -> Relationship-specific Context Projection`

A context item is not friend-visible merely because it exists. The relevant relationship and privacy policy must determine whether it can be disclosed and at what detail level.

### Pulse generation

`Visible Context + Relevance -> Friend Pulse`

Relevance operates only on context that is already eligible for the relationship. Ranking must not bypass privacy.

### Conversation support

`Friend Pulse / Visible Context -> Conversation / Interaction`

Conversation support consumes context that the user is already allowed to see. It must not increase disclosure beyond the underlying context projection.

### Notification

`Relevant Product Intent -> Notification`

Notification delivers an already-approved intent. It must not independently reinterpret privacy or generate new social meaning.

---

## Initial ubiquitous language

### Activity

An authorized signal about something a user did, consumed, saved, liked, watched, listened to, researched, or contributed.

### Raw Activity

Provider-shaped or source-shaped activity before Kin normalization and social interpretation.

### Normalized Activity

A provider-independent Kin representation of an activity signal.

### Context Candidate

A possible piece of socially meaningful information derived from one or more authorized signals. It is not necessarily shareable or worth surfacing.

### Social Context

A validated, meaningful description of something worth potentially knowing about a user.

### Context Projection

The relationship-specific representation of a context item after privacy/sharing rules have been applied.

### Friendship

An explicit relationship between two Kin users.

### Relationship Level

A classification of relational closeness/trust that may affect sharing policy. Candidate long-term levels include Acquaintance, Friend, Close Friend, and Inner Circle.

The MVP may intentionally begin with only a single close-friend or inner-circle relationship model rather than exposing the full hierarchy.

### Friend Pulse

A small, prioritized set or summary of meaningful, permissioned context about a friend. It is not a chronological activity feed.

### Significant Change

A meaningful shift or recurring pattern in a user's interests or behavior that is more informative than a one-off raw action.

### Shared Rabbit Hole

A situation where two friends are independently exploring the same or closely related topic and that overlap may be useful for conversation.

### Friendship Drift

A long-term condition where a relationship has become less active or context continuity has weakened. This is a future product concept, not necessarily MVP scope.

### Social Memory

Long-lived relationship context that helps preserve shared history or continuity without exposing unnecessary raw records.

---

## Core, supporting, and generic domain hypotheses

These classifications are provisional.

### Likely core domains

- Friendship
- Social Context
- Privacy & Sharing
- potentially Relevance
- potentially Activity, depending on how much domain-specific behavior emerges around signal meaning

These areas most directly express Kin's differentiation: preserving context for real relationships while controlling what each friend should know.

### Likely supporting domains

- Identity
- Conversation / Interaction in early phases
- Activity if it remains mostly normalization and lifecycle behavior

### Likely generic domains

- Integration plumbing
- Notification delivery
- infrastructure-level provider connectivity

A domain should not be promoted to "core" merely because it is technically difficult.

---

## Privacy model direction

The long-term privacy model should be relationship-aware rather than a single global public/private switch.

Candidate relationship levels:

- Acquaintance
- Friend
- Close Friend
- Inner Circle

The exact taxonomy is not yet fixed.

The domain concept that matters is that one source context can safely project different levels of detail to different relationships.

The MVP should avoid implementing the complete relationship-level system before validating that close-friend context itself is valuable.

---

## AI boundary

AI is an external capability, not a domain authority.

A conceptual application flow may look like:

`Authorized Signals -> Application Use Case -> ContextGenerator Port -> AI Adapter -> Normalization/Validation -> Context Candidate -> Domain Rules`

Rules:

- model names do not belong in domain code;
- prompts do not define domain invariants;
- raw LLM JSON does not become a domain object directly;
- provider-specific failures are translated at the adapter boundary;
- domain rules determine whether a normalized candidate is valid, shareable, or useful.

---

## Long-term capability map vs MVP scope

This document intentionally includes future capabilities so the system can evolve coherently.

Their presence here does **not** authorize implementation.

Possible long-term capabilities include:

- relationship levels beyond close friends;
- Friends Pulse;
- AI Q&A about permitted friend context;
- conversation starters;
- weekly friendship digest;
- serendipity matching;
- friendship drift detection;
- shared rabbit holes;
- social memory;
- AI friendship concierge;
- friend-aware relevance;
- ChatGPT/MCP-style access to permissioned Kin context;
- broader integrations across media, reading, browsing, AI research, and other digital activity.

The MVP roadmap is the source of truth for sequencing and current validation scope.

---

## MVP direction

The initial product should validate the smallest useful form of context continuity among a very small number of close relationships.

Expected direction:

- start with roughly 3-5 close friends / inner-circle relationships;
- use a small number of input sources;
- derive only a small number of meaningful context items;
- prefer 1-3 useful contexts over a dense feed;
- use "Start a conversation" as the important action rather than Like;
- introduce automated provider integrations only after the domain flow can work with a simpler contribution path.

The detailed sequence belongs in `docs/product/mvp-roadmap.md`.

---

## Architectural implications

The product model implies the following implementation direction without defining concrete infrastructure:

- domain responsibilities should be modeled before persistence;
- write behavior and read behavior should remain conceptually separate;
- commands express state-changing intent;
- queries optimize for interactions without distorting write-side models;
- domain events may represent meaningful completed business facts;
- external integrations and AI systems remain adapters;
- infrastructure should plug into domain/application ports;
- bounded contexts are initially modules inside a modular monolith, not independent services;
- Event Sourcing is not required by this product model.

---

## Boundary questions to revisit

The following questions are deliberately unresolved and should be answered through product discovery and implementation evidence:

1. Does Activity become a rich core domain or remain a supporting normalization layer?
2. Does Relevance deserve its own bounded context, or is it initially an application policy around Social Context?
3. Does Conversation / Interaction become a domain with its own lifecycle, or remain a thin application capability?
4. Should relationship closeness live entirely in Friendship, while Privacy & Sharing only references it, or should some trust policy be represented independently?
5. At what point does Social Memory become distinct from Social Context?
6. Which concepts need strong consistency within one aggregate, and which can remain eventually consistent read projections?

These questions should not be resolved by choosing database tables first.

---

## Source-of-truth relationship

This document defines long-term product/domain context.

Implementation precedence remains:

1. current GitHub issue and explicit acceptance criteria/non-goals;
2. current MVP roadmap slice;
3. this product scope and domain map;
4. root and local `AGENTS.md` contracts;
5. architecture/testing/workflow skills;
6. existing implementation conventions.

Future capabilities described here are context, not implementation permission.
